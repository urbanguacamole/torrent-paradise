package main

import (
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

func main() {
	f, err := os.Open("/home/nextgen/torrent_dump_full.csv.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	db := initDb()

	cr := csv.NewReader(gr)
	cr.LazyQuotes = true
	cr.Comma = rune(';')
	const layout = "2006-Jan-02 15:04:05"
	if err != nil {
		log.Fatal(err)
	}
	_, err = cr.Read() // read first line and throw it away
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, error := cr.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			if perr, ok := err.(*csv.ParseError); ok && perr.Err == csv.ErrFieldCount {
				log.Println(err)
			}
		}

		added, err := time.Parse(layout, line[0])
		if err != nil {
			log.Println(err)
		}
		ihBytes, _ := base64.StdEncoding.DecodeString(line[1])
		ih := hex.EncodeToString(ihBytes)
		_, err = db.Exec("INSERT INTO torrent (infohash, name, length, added) VALUES ($1, $2, $3, $4)", ih, line[2], line[3], added)
		if err, ok := err.(*pq.Error); ok { //dark magic
			if err.Code != "23505" {
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
