package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

	if err != nil {
		log.Println(err)
		return nil
	}

	var resp []ApibayTorrent
	err = json.Unmarshal(body, &resp)

	if err != nil {
		var respFromIdiots []ApibayTorrentTheyAreIdiots
		err = json.Unmarshal(body, &respFromIdiots)
		if err != nil {
			log.Println(err)
			return nil
		}

		for _, torrByIdiot := range respFromIdiots {
			var transl ApibayTorrent
			transl.Info_hash = torrByIdiot.Info_hash
			transl.Name = torrByIdiot.Name
			transl.Size, err = strconv.Atoi(torrByIdiot.Size)
			transl.Added, err = strconv.Atoi(torrByIdiot.Added)
			resp = append(resp, transl)
		}
	}

	var torrents []Torrent
	for _, apibayTorr := range resp {
		torrents = append(torrents, Torrent{apibayTorr.Info_hash, apibayTorr.Name, apibayTorr.Size})
	}
	return torrents
}

// ApibayTorrent Structure returned from apibay. For unmarshaling from JSON. Not all fields that are returned from Apibay are in this struct; YAGNI
type ApibayTorrent struct {
	Info_hash string
	Name      string
	Size      int
	Added     int
}

type ApibayTorrentTheyAreIdiots struct {
	Info_hash string
	Name      string
	Size      string
	Added     string
}
