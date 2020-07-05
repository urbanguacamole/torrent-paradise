package main

import (
	"strconv"
	"testing"
)

func TestCrawlTPB48hTop(t *testing.T) {
	torrents := CrawlTPB48hTop()
	if len(torrents) < 1 {
		t.Error("no torrents crawled from tpb")
	}
	for i, torrent := range torrents {
		if torrent.Length < 10 {
			t.Error("bad length of torrent "+strconv.Itoa(i))
		}
		if len(torrent.Name) < 2 {
			t.Error("weirdly short name of torrent "+strconv.Itoa(i))
		}
	}
}