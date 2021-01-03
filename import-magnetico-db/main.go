package main

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lib/pq"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	db := initDb()

	cr := csv.NewReader(f)
	cr.LazyQuotes = true
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			if perr, ok := err.(*csv.ParseError); ok && perr.Err == csv.ErrFieldCount {
				log.Println(err)
			}
		}

		infohash := line[0]
		if len(infohash) != 40 {
			log.Fatal("bad infohash length " + line[0])
		}

		name := line[1]
		if len(name) < 2 {
			log.Println("bad name length " + line[1])
			continue
		}
		if !utf8.ValidString(name) {
			log.Println("utf8 invalid name")
			log.Println(name)
			continue
		}

		length := line[2]

		addedUnix, err := strconv.ParseInt(line[3], 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		added := time.Unix(addedUnix, 0)

		name = strings.ToLower(name)

		//fmt.Printf("Ih %v name %v len %v added %v", infohash, name, length, added)
		_, err = db.Exec("INSERT INTO torrent (infohash, name, length, added) VALUES ($1, $2, $3, $4)", infohash, name, length, added)
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
