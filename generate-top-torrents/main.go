package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type torrent struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Len  int    `json:"len"`
	S    int    `json:"s"`
	L    int    `json:"l"`
	C    int    `json:"c"`
}

func main() {
	// https://www.dotnetperls.com/csv-go
	f, err := os.Open("index-generator/dump.csv") // expects that you will run it from the parent dir
	if err != nil {
		log.Fatal(err)
	}
	records, err := csv.NewReader(bufio.NewReader(f)).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(records, func(i, j int) bool {
		return toInt(records[i][3]) > toInt(records[j][3])
	})
	var topResults [500]torrent
	for i, record := range records {
		if i < 500 {
			topResults[i] = torrent{string(record[0]), string(record[1]), toInt(record[2]), toInt(record[3]), toInt(record[4]), toInt(record[5])}
		}
	}
	out, err := json.Marshal(topResults)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
