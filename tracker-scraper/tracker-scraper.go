package main

import (
	"database/sql"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/etix/goscrape"
	_ "github.com/lib/pq"
)

const retryLimit = 3
const waitTime = 250 // in ms
var trackers = [3]string{"udp://tracker.coppersurfer.tk:6969", "udp://exodus.desync.com:6969", "udp://tracker.pirateparty.gr:6969"}

func main() {
	db := initDb()
	trackerResponses := make(chan trackerResponse, 100)
	trackerRequests := make(map[string]chan []string)

	var counter uint64 //count of torrents scraped

	quitCounter := make(chan bool)
	go func() {
		for {
			select {
			case <-quitCounter:
				return
			default:
				log.Println("Torrents scraped so far: " + strconv.Itoa(int(atomic.LoadUint64(&counter))))
				time.Sleep(2 * time.Second)
			}
		}
	}()

	for _, tracker := range trackers {
		trackerRequests[tracker] = make(chan []string, 1000)
		go runScraper(trackerRequests[tracker], trackerResponses, tracker, waitTime, &counter)
	}

	datestamp := time.Now().Local().Format("2006-01-02")

	go runPersister(trackerResponses, db, datestamp)

	rows, err := db.Query("SELECT infohash FROM torrent WHERE NOT EXISTS (SELECT FROM trackerdata WHERE infohash = torrent.infohash AND scraped = '" + datestamp + "')")
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
			for _, tracker := range trackers {
				trackerRequests[tracker] <- infohashes
			}
			infohashes = []string{}
		}
	}
	for _, tracker := range trackers {
		trackerRequests[tracker] <- infohashes
	}

	quitCounter <- true

	for len(trackerRequests) > 0 {
		time.Sleep(2 * time.Second)
		var left int
		for i, tracker := range trackers {
			left = len(trackerRequests[tracker])
			if left != 0 {
				log.Println("Tracker " + strconv.Itoa(i) + " requests left to send: " + strconv.Itoa(len(trackerRequests[tracker])))
			}
		}
	}

	for _, tracker := range trackers {
		close(trackerRequests[tracker])
	}

	for len(trackerResponses) > 0 {
		time.Sleep(2 * time.Second)
		log.Println("Tracker responses left to save: " + strconv.Itoa(len(trackerResponses)))
	}
	time.Sleep(time.Duration(waitTime*retryLimit) * time.Millisecond)
	close(trackerResponses)
}

//Runs a tracker that scrapes the given tracker. Takes requests from trackerRequests and sends responses to trackerResponses
//waittime is in miliseconds
func runScraper(trackerRequests chan []string, trackerResponses chan trackerResponse, tracker string, waittime int, counter *uint64) {
	s, err := goscrape.New(tracker)
	s.SetTimeout(time.Duration(waitTime) * time.Millisecond)
	s.SetRetryLimit(retryLimit)
	if err != nil {
		log.Fatal("Error:", err)
	}
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
			log.Println(err)
		} else {
			atomic.AddUint64(counter, uint64(len(infohashes)))
			trackerResponses <- trackerResponse{tracker, res}
		}

		time.Sleep(time.Duration(waittime) * time.Millisecond)
	}
}

func runPersister(trackerResponses chan trackerResponse, db *sql.DB, datestamp string) {
	for res := range trackerResponses {
		for _, scrapeResult := range res.scrapeResult {
			_, err := db.Exec("INSERT INTO trackerdata (infohash, tracker, seeders, leechers, completed, scraped) VALUES ($1, $2, $3, $4, $5, $6)", scrapeResult.Infohash, res.tracker, scrapeResult.Seeders, scrapeResult.Leechers, scrapeResult.Completed, datestamp)
			if err != nil {
				log.Fatal(err)
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

	/*_, err = db.Exec(`CREATE TYPE tracker AS ENUM ('udp://tracker.coppersurfer.tk:6969', 'udp://exodus.desync.com:6969', 'udp://tracker.pirateparty.gr:6969')`)
	if err != nil {
		log.Fatal(err)
	}*/

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trackerdata (
		infohash char(40),
		tracker tracker,
		seeders int NOT NULL,
		leechers int NOT NULL,
		completed int NOT NULL,
		scraped char(10),
		PRIMARY KEY (infohash, scraped, tracker)
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type trackerResponse struct {
	tracker      string
	scrapeResult []*goscrape.ScrapeResult
}
