package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/etix/goscrape"
	_ "github.com/lib/pq"
)

const retryLimit = 10
const waitTime = 500 // in ms
var trackers = [4]string{"udp://tracker.coppersurfer.tk:6969", "udp://tracker.internetwarriors.net:1337", "udp://tracker.opentrackr.org:1337", "udp://tracker.pirateparty.gr:6969/announce"}

func main() {
	db := initDb()
	trackerResponses := make(chan trackerResponse, 100)
	trackerRequests := make(chan []string, 1000)

	for _, tracker := range trackers {
		go runScraper(trackerRequests, trackerResponses, tracker, waitTime)
	}

	go runPersister(trackerResponses, db)

	rows, err := db.Query("SELECT infohash FROM torrent WHERE NOT EXISTS (SELECT FROM peercount WHERE infohash = torrent.infohash)")
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
			trackerRequests <- infohashes
			infohashes = []string{}
		}
	}
	trackerRequests <- infohashes
	for len(trackerRequests) > 0 {
		time.Sleep(2 * time.Second)
		log.Println("Tracker requests left to send: " + strconv.Itoa(len(trackerRequests)))
	}
	close(trackerRequests)

	for len(trackerResponses) > 0 {
		time.Sleep(2 * time.Second)
		log.Println("Tracker responses left to save: " + strconv.Itoa(len(trackerResponses)))
	}
	time.Sleep(time.Duration(waitTime*retryLimit) * time.Millisecond)
	close(trackerResponses)
}

//Runs a tracker that scrapes the given tracker. Takes requests from trackerRequests and sends responses to trackerResponses
//waittime is in miliseconds
func runScraper(trackerRequests chan []string, trackerResponses chan trackerResponse, tracker string, waittime int) {
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
			trackerResponses <- trackerResponse{tracker, res}
		}

		time.Sleep(time.Duration(waittime) * time.Millisecond)
	}
}

func runPersister(trackerResponses chan trackerResponse, db *sql.DB) {
	for res := range trackerResponses {
		for _, scrapeResult := range res.scrapeResult {
			_, err := db.Exec("INSERT INTO peercount (infohash, tracker, seeders, leechers, completed) VALUES ($1, $2, $3, $4, $5)", scrapeResult.Infohash, res.tracker, scrapeResult.Seeders, scrapeResult.Leechers, scrapeResult.Completed)
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS peercount (
		infohash char(40),
		tracker varchar,
		seeders int NOT NULL,
		leechers int NOT NULL,
		completed int NOT NULL,
		scraped timestamp DEFAULT current_timestamp,
		PRIMARY KEY (infohash, tracker, scraped)
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
