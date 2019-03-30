package main

import (
	"log"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

type Tomlconf struct {
	Trackers    []string
	WaitTime    string            // time to wait between requests to one tracker
	LogInterval string            // interval between stats dumps to console
	Categories  map[string]string // Defines acceptable freshness of seed/leech counts for categories of torrents. Category number is the minimum seed count for torrent to be assigned to a category. Each torrent/tracker pair is fetched independently. All torrent/tracker pairs are in the highest category available to it.
}

func loadConfig() config {
	var tomlconf Tomlconf //stores parsed config, ready for translation from strings to durations.
	if _, err := toml.DecodeFile("config.toml", &tomlconf); err != nil {
		log.Fatal(err)
	}

	conf := config{trackers: tomlconf.Trackers}
	wt, err := time.ParseDuration(tomlconf.WaitTime)
	if err != nil {
		log.Println("f")
		log.Println(tomlconf)
		log.Fatal(err)
	}
	conf.waitTime = wt
	li, err := time.ParseDuration(tomlconf.LogInterval)
	if err != nil {
		log.Println(tomlconf.LogInterval)
		log.Fatal(err)
	}
	conf.logInterval = li
	conf.categories = make(map[int]time.Duration)
	for k, v := range tomlconf.Categories {
		ka, err := strconv.Atoi(k)
		conf.categories[ka], err = time.ParseDuration(v)
		if err != nil {
			log.Println(v)
			log.Fatal(err)
		}
	}
	return conf
}
