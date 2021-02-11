package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/lib/pq"
)

type Results []Result

type Result struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Len  int    `json:"len"`
	S    int    `json:"s"`
	L    int    `json:"l"`
}

type ingestRequest struct {
	Ih     string
	Name   string
	Length int
}

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

	var alreadyIngested sync.Map

	http.HandleFunc("/api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		_, err = db.Exec("INSERT INTO telemetry (payload) VALUES ($1)", string(body))
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	})

	// INSECURE! BE SURE TO WALL OFF THIS PROPERLY IN PROD
	http.HandleFunc("/api/ingest", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req ingestRequest
		err := decoder.Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		_, cached := alreadyIngested.Load(req.Ih)
		if !cached {
			_, err = db.Exec("INSERT INTO torrent (infohash, name, length) VALUES ($1, $2, $3)", req.Ih, req.Name, req.Length)
			if err, ok := err.(*pq.Error); ok { //dark magic
				if err.Code != "23505" {
					log.Println(err)
				}
			}
			alreadyIngested.Store(req.Ih, true)
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if len(q) == 0 {
			log.Print("/api/search received empty q argument")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pedoWords := []string{"pedo", "preteen", "14yo", "13yo", "12yo", "15yo", "16yo", "17yo", "11y", "underage", "11yo", "10yo", "10y", "9yo", "mafiasex", "pedofilia", "violacion"}
		for _, word := range pedoWords {
			if strings.Contains(q, word) {
				w.WriteHeader(500)
				return
			}
		}
		rows, err := db.Query("select infohash, name, length, s, l from search where vect @@ websearch_to_tsquery($1) and copyrighted = 'f' order by s desc limit 150", q)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var results Results
		for rows.Next() {
			var (
				infohash string
				name     string
				length   int
				s        int
				l        int
			)
			err := rows.Scan(&infohash, &name, &length, &s, &l)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}
			result := Result{infohash, name, length, s, l}
			results = append(results, result)
		}
		w.WriteHeader(http.StatusOK)
		marshaledResults, err := json.Marshal(results)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(marshaledResults)
	})

	http.ListenAndServe(":8000", nil)
}
