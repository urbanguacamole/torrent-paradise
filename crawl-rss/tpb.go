package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func CrawlTPB48hTop() []Torrent {
	return parseApibayJSON("https://apibay.org/precompiled/data_top100_48h.json")
}

func CrawlTPBVideoRecent() []Torrent {
	return parseApibayJSON("https://apibay.org/q.php?q=category%3A200")
}

func parseApibayJSON(url string) []Torrent {
	httpresp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer httpresp.Body.Close()
	body, err := ioutil.ReadAll(httpresp.Body)

	var resp []ApibayTorrent
	err = json.Unmarshal(body, &resp)

	var torrents []Torrent
	for _, apibayTorr := range resp {
		torrents = append(torrents, Torrent{apibayTorr.Info_hash, apibayTorr.Name, apibayTorr.Size})
	}
	return torrents
}

// ApibayTorrent Structure returned from apibay. For unmarshaling from JSON. Not all fields that are returned from Apibay are in this struct; YAGNI
type ApibayTorrent struct {
	ID        int
	Info_hash string
	Name      string
	Size      int
	Added     int
}
