# About
## What is this?

If you don't know what Torrent Paradise is, see the [website](https://torrent-paradise.ml/about.html).

This is a repository of all the tools I use to build and run torrent-paradise.ml. The 'code name' of the project is nextgen (next gen torrent search), so don't be surprised if it comes up somewhere.

## Can you help me?
Maybe, open an issue. Be sure to demonstrate an effort that you tried to solve the problem yourself.

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
  - node v10.15 & npm
  - Python 3 (required only for index-generator/fix-metajson.py)

The programs create their own tables in the DB that they need. Database name is "nextgen".

What I did first after getting the server up and running was importing the TPB dump. Download https://thepiratebay.org/static/dump/csv/torrent_dump_full.csv.gz to the import-tpb-dump directory and run `go run`.

There is a complete database dump available in torrentparadise-staticbackup.torrent, so you don't have to do that. This same database dump is available on https://mega.nz/#!ddESlChb!3YBqfxG-a4fwpXzPG3QsXa-C6FeQ9AbNSGXxY7W7xm4. It contains the same data as the torrent, only .xz compressed.

# Usage

## Generate the index

See `update-index.sh`.

## Spider the DHT

Run `go build` in spider/ to compile and scp the binary it to the server. You can use the systemd service file in `spider/spider.service` to start the spider on server boot.

## Scraping trackers for seed/leech data

Run `go build` in seedleech-daemon/ to compile and scp the binary it to the server. You can use the systemd service file in `seedleech-daemon/seedleech.service`.

## IPFS vs 'static'

The directory website gets deployed to IPFS, static gets deployed to the server. Static calls the API, the IPFS version doesn't.

# Contributing

Before working on something, open an issue to ask if it would be okay. I would love to [KISS](https://en.wikipedia.org/wiki/KISS_principle). 
