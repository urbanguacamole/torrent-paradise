package main

import (
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/lib/pq"
	"github.com/urbanguacamole/dht"
)

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func main() {

	db := initDb()

	w := dht.NewWire(65536, 2048, 512)
	go handleResponses(w, db)
	go w.Run()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)
	d.Run()
}

func handleResponses(w *dht.Wire, db *sql.DB) {
	for resp := range w.Response() {
		metadata, err := dht.Decode(resp.MetadataInfo)
		if err != nil {
			continue
		}
		info := metadata.(map[string]interface{})

		if _, ok := info["name"]; !ok {
			continue
		}

		bt := bitTorrent{
			InfoHash: hex.EncodeToString(resp.InfoHash),
			Name:     info["name"].(string),
		}

		var length int

		if v, ok := info["files"]; ok {
			files := v.([]interface{})
			bt.Files = make([]file, len(files))

			for _, item := range files {
				f := item.(map[string]interface{})
				length += f["length"].(int)
			}
		} else if _, ok := info["length"]; ok {
			length = info["length"].(int)
		}

		_, err = db.Exec("INSERT INTO torrent (infohash, name, length) VALUES ($1, $2, $3)", bt.InfoHash, bt.Name, length)
		if err, ok := err.(*pq.Error); ok { //dark magic
			if err.Code != "23505" {
				log.Println(err)
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
