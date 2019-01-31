# About
## What is this?

This is a repository of all the tools I use to build and run torrent-paradise.ml. Some people asked for a source, so I'm just putting this out here. I did make *some* effort to clean it up. The 'code name' of the project is nextgen (next gen torrent search), so don't be surprised if it comes up somewhere.

## Can you help me?
Maybe, open an issue. Be sure to demonstrate an effort that you tried to solve the problem yourself.

## This is a big mess. Fix it maybe?
WIP ❤️

# Setup

Here's what the setup looks like rn:
- VPS, Debian Stretch, 2 GB RAM
  - PostgreSQL 9.6. pg_hba.conf contains this:

    ```
    local   all             all                                      peer
    # IPv4 local connections:
    host    nextgen         nextgen          localhost               md5
    ```
  - IPFS v0.4.18
  - user with username nextgen on the server
- my laptop w/ Linux
  - Go toolchain installed
  - node v10.9.0 & npm

Schema for the database is sth like this (taken from index-generator/README, runs on sqlite, probably also on pg.)
```sql
CREATE TABLE peercount ( infohash char(40), tracker varchar, seeders int, leechers int, completed int, scraped timestamp);

CREATE TABLE torrent( infohash char(40), name varchar, length bigint, added timestamp);
```



What I did first after getting the server up and running was importing the TPB dump. Download https://thepiratebay.org/static/dump/csv/torrent_dump_full.csv.gz to the import-tpb-dump directory and run `go run`.

I probably forgot sth. Open an issue!

# Usage

## Generate the index

This is a half-broken process that is partially described in update-index.sh. Read the script to understand what it does.

## Spider the DHT

Run `go build` in spider/ to compile and scp the binary it to the server. You can use the systemd service file in `spider/spider.service` to start the spider on boot.


# Contributing

Before working on something, open an issue to ask if it would be okay. I would love to [KISS](https://en.wikipedia.org/wiki/KISS_principle). 