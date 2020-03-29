package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
)

type Torrent struct {
	Infohash string
	Name     string
	Length   int
}

func main() {
	db := initDb()
	crawled := make(map[string]bool) // set to not needlessly send all torrents to db to check if we found them already
	for {
		torrents := CrawlYts()
		for _, torrent := range torrents {
			addTorrent(db, torrent, &crawled)
		}
		time.Sleep(time.Minute * 45)
		go refresh(db)
	}
}

func refresh(db *sql.DB) {
	db.Exec("REFRESH MATERIALIZED VIEW fresh")
	db.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY search")
}

func addTorrent(db *sql.DB, torr Torrent, crawled *map[string]bool) {
	if !(*crawled)[string(torr.Infohash)] {
		_, err := db.Exec("INSERT INTO torrent (infohash, name, length) VALUES ($1, $2, $3)", torr.Infohash, torr.Name, torr.Length)
		if err, ok := err.(*pq.Error); ok { //dark magic
			if err.Code != "23505" {
				log.Fatal(err)
			}
		}
		(*crawled)[torr.Infohash] = true
	}
}

//todo https://rarbg.to/rssdd.php?category=44
func CrawlYts() []Torrent {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://yts.mx/rss/0/all/all/0")
	if err != nil {
		log.Fatal(err)
	}
	var torrents []Torrent
	for _, item := range feed.Items {
		size, err := parseSizeYts(item.Description)
		if err != nil {
			log.Print(err)
			continue
		}
		ih, err := parseInfohashYts(item.Enclosures[0].URL)
		torrents = append(torrents, Torrent{ih, item.Title, size})
	}
	return torrents
}

// Parses torrent length from YTS description
func parseSizeYts(description string) (int, error) {
	s := strings.Split(description, "<br />Size: ")
	if len(s) == 0 {
		return 0, errors.New("Couldn't find '<br />Size: ' in item description")
	}
	s = strings.Split(s[1], "B<br />Runtime: ")
	if len(s) == 0 {
		return 0, errors.New("Couldn't find 'B<br />Runtime: ' in item description")
	}
	length, err := strconv.ParseFloat(s[0][:len(s[0])-2], 64)
	if err != nil {
		return 0, err
	}
	if s[0][len(s[0])-1:] == "G" {
		return int(length * 1000000000), nil
	} else if s[0][len(s[0])-1:] == "M" {
		return int(length * 1000000), nil
	} else {
		return 0, errors.New("Invalid char in place of length specifier")
	}
}

func parseInfohashYts(url string) (string, error) {
	s := strings.Split(url, "torrent/download/")
	if len(s) == 0 {
		return "", errors.New("invalid URL")
	}
	return strings.ToLower(s[1]), nil
}

func initDb() *sql.DB {
	connStr := "user=nextgen dbname=nextgen host=/var/run/postgresql"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS torrent (
		infohash char(40) PRIMARY KEY NOT NULL,
		name varchar NOT NULL,
		length bigint,
		added timestamp DEFAULT current_timestamp
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
