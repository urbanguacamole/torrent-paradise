package main

import "testing"

func TestCrawlYts(t *testing.T) {
	torrents := CrawlYts()
	if len(torrents) < 1 {
		t.Error("no torrents crawled from yts")
	}
}

func TestCrawlEztv(t *testing.T) {
	t.Log("t.log")
	torrents := CrawlEztv()
	if len(torrents) < 1 {
		t.Error("no torrents crawled from eztv.io")
	}
	t.Log(torrents[0].Name)
	t.Log(torrents[1].Name)
}
