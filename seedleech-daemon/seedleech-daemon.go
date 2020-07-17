// the seedleech daemon is designed to keep seed/leech counts as fresh as possible automatically and run 24/7
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/etix/goscrape"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Config file definition. Loaded from config.toml in workdir. All fields mandatory.
type config struct {
	trackers    []string
	waitTime    time.Duration         // time to wait between requests to one tracker
	logInterval time.Duration         // interval between stats dumps to console
	categories  map[int]time.Duration // Defines acceptable freshness of seed/leech counts for categories of torrents. Category number is the minimum seed count for torrent to be assigned to a category. Each torrent/tracker pair is fetched independently. All torrent/tracker pairs are in the highest category available to it.
}

type trackerResponse struct {
	tracker      string
	scrapeResult []*goscrape.ScrapeResult
}

var conf config

func main() {
	conf = loadConfig()
	db := initDb()

	trackerResponses := make(chan trackerResponse, 100)
	trackerRequests := make(map[string]chan []string)

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	quit := false
	go func() {
		<-sigterm
		quit = true
	}()

	for _, tracker := range conf.trackers {
		trackerRequests[tracker] = make(chan []string) // non buffered (!)
		go runScraper(trackerRequests[tracker], trackerResponses, tracker)
		for minseed, delay := range conf.categories {
			go runWorkFetcher(trackerRequests[tracker], tracker, minseed, delay, &quit, db)
		}
	}

	go runPersister(trackerResponses, db)

	for {
		time.Sleep(conf.logInterval / 2)
		if quit {
			return
		}
		time.Sleep(conf.logInterval / 2)
		if quit {
			return
		}
		for _, tracker := range conf.trackers {
			for minSeed, maxAge := range conf.categories {
				freshlimit := time.Now().Local().Add(-maxAge)
				if minSeed != 0 {
					var res int
					row := db.QueryRow("SELECT count(1) FROM trackerdata WHERE tracker = $1 AND seeders > $2 AND scraped < $3", tracker, minSeed, freshlimit)
					row.Scan(&res)
					if res > 0 {
						fmt.Println("Tracker " + tracker + ", seeds > " + strconv.Itoa(minSeed) + ": " + strconv.Itoa(res))
					}
				} else {
					var res int
					row := db.QueryRow("SELECT count(1) from torrent")
					row.Scan(&res)
					totalTorrents := res
					row = db.QueryRow("SELECT count(1) from trackerdata where tracker = $1", tracker)
					row.Scan(&res)
					if (totalTorrents - res) > 0 {
						fmt.Println("Tracker " + tracker + ", seeds = 0: " + strconv.Itoa(totalTorrents-res))
					}
				}
			}
		}
	}
}

// a work fetcher for a given tracker and category combination
func runWorkFetcher(trackerRequests chan []string, tracker string, minseed int, maxAge time.Duration, quit *bool, db *sql.DB) {
	for {
		if *quit {
			fmt.Println("Workfetcher for category " + strconv.Itoa(minseed) + ", tracker " + tracker + " stopping.")
			return
		}
		freshlimit := time.Now().Local().Add(-maxAge)
		var rows *sql.Rows
		var err error
		if minseed != 0 {
			rows, err = db.Query("SELECT infohash FROM trackerdata WHERE tracker = $1 AND seeders > $2 AND scraped < $3 LIMIT 740", tracker, minseed, freshlimit)
		} else {
			//time.Sleep(time.Duration(int64(rand.Intn(12000)) * int64(time.Second))) //sleep for random time between 0 mins and 200 mins
			rows, err = db.Query("SELECT infohash FROM torrent WHERE NOT EXISTS (SELECT from trackerdata WHERE infohash = torrent.infohash AND tracker = $1 AND scraped > $2)", tracker, freshlimit)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var infohashes []string
		for rows.Next() {
			var infohash string
			if err := rows.Scan(&infohash); err != nil {
				log.Fatal(err)
			}
			if len(infohashes) < 74 {
				infohashes = append(infohashes, infohash)
			} else {
				if *quit {
					fmt.Println("Workfetcher for category " + strconv.Itoa(minseed) + ", tracker " + tracker + " stopping.")
					return
				}
				trackerRequests <- infohashes
				infohashes = []string{}
			}
		}
		trackerRequests <- infohashes
		time.Sleep(time.Minute)
	}
}

// a scraper for one tracker
func runScraper(trackerRequests chan []string, trackerResponses chan trackerResponse, tracker string) {
	s, err := goscrape.New(tracker)
	s.SetTimeout(conf.waitTime)
	s.SetRetryLimit(1)
	if err != nil {
		log.Fatal("Error:", err)
	}
	success := 0 //how many times request to this tracker has succeeded
	failure := 0 //how many times request to this tracker has failed
	for req := range trackerRequests {
		infohashes := make([][]byte, len(req))
		for i, v := range req {
			if len(v) != 40 { //infohashes are 40 chars long in string representation.
				panic("Infohash in trackerRequest with index " + strconv.Itoa(i) + " isn't 40 chars long, it's " + strconv.Itoa(len(v)) + " long.")
			}
			infohashes[i] = []byte(v)
		}
		res, err := s.Scrape(infohashes...)
		if err != nil {
			failure++
		} else {
			trackerResponses <- trackerResponse{tracker, res}
			success++
		}

		if (failure + success) > 99 {
			if failure > 50 {
				log.Println("unable to communicate with tracker, " + strconv.Itoa(failure) + "reqs of " + strconv.Itoa(failure+success) + " failed")
				log.Println(tracker)
				time.Sleep(time.Hour)
			}
			failure = 0
			success = 0
		}

		time.Sleep(conf.waitTime)
	}
}

func runPersister(trackerResponses chan trackerResponse, db *sql.DB) {
	for res := range trackerResponses {
		for _, scrapeResult := range res.scrapeResult {
			// TODO check if trackerdata for torrent/tracker combo aren't in DB already, if no, insert, if yes, update
			timestamp := time.Now()
			_, err := db.Exec("INSERT INTO trackerdata (infohash, tracker, seeders, leechers, completed, scraped) VALUES ($1, $2, $3, $4, $5, $6)", scrapeResult.Infohash, res.tracker, scrapeResult.Seeders, scrapeResult.Leechers, scrapeResult.Completed, timestamp)
			if pgerr, ok := err.(*pq.Error); ok {
				if pgerr.Code == "23505" {
					_, err := db.Exec("UPDATE trackerdata SET seeders = $3, leechers = $4, completed = $5, scraped = $6 WHERE infohash = $1 AND tracker = $2", scrapeResult.Infohash, res.tracker, scrapeResult.Seeders, scrapeResult.Leechers, scrapeResult.Completed, timestamp)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Fatal(err)
				}
			}
		}
	}
}

func initDb() *sql.DB {
	connStr := "user=nextgen dbname=nextgen host=/var/run/postgresql"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tracker') THEN
        CREATE TYPE tracker AS ENUM ('udp://tracker.coppersurfer.tk:6969', 'udp://tracker.leechers-paradise.org:6969/announce', 'udp://exodus.desync.com:6969', 'udp://tracker.pirateparty.gr:6969','udp://tracker.opentrackr.org:1337/announce','udp://tracker.internetwarriors.net:1337/announce','udp://tracker.cyberia.is:6969/announce','udp://9.rarbg.to:2920/announce');
    END IF;
END$$`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trackerdata (
		infohash char(40),
		tracker tracker,
		seeders int NOT NULL,
		leechers int NOT NULL,
		completed int NOT NULL,
		scraped timestamp,
		PRIMARY KEY (infohash, tracker)
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
