package main

import (
	"database/sql"
	"log"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

func initDb() *sql.DB {
	connStr := "user=nextgen dbname=nextgen host=/var/run/postgresql"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS telemetry (
		payload jsonb NOT NULL,
		time timestamp DEFAULT current_timestamp
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := initDb()

    http.HandleFunc("/api/telemetry", func (w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			log.Print(err)
		}
		_, err = db.Exec("INSERT INTO telemetry (payload) VALUES ($1)", string(body))
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
    })
	http.ListenAndServe(":8000", nil)
}