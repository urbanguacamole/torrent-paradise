package main

import (
	"testing"
)

func TestCrawlYts(t *testing.T) {
	torrents := CrawlYts()
	if len(torrents) < 1 {
		t.Error("no torrents crawled from yts")
	}
}
