# About
## What is this?

If you don't know what Torrent Paradise is, see the [website](https://torrent-paradise.ml/about.html).

This is a repository of all the tools I use to build and run torrent-paradise.ml. The 'code name' of the project is nextgen (next gen torrent search), so don't be surprised if it comes up somewhere.

## Can you help me?
Maybe, open an issue. Be sure to demonstrate an effort that you tried to solve the problem yourself.

# Setup

Here's what the setup looks like rn:
- VPS, Debian Bullseye, 8 GB RAM
  - user with username nextgen on the server
- my laptop w/ Linux
  - Go toolchain installed
  - node & npm
  - Python 3 (required only for index-generator/fix-metajson.py)

Read the server-setup.md file for more precise info.

The programs create their own tables in the DB that they need. Database name is "nextgen". You need to create the materialized views (fresh and search). You can find some useful SQL code in snippets.sql.

Each of the daemons (api, crawl-rss, seedleech-daemon) is its own standalone Go package and resulting binary. You have to compile the binaries yourself. There are systemd .service files available for each of the daemons.

The torrent collection is a mashup of the (now no longer provided) TPB dumps, my own DHT spidering efforts, and [magnetico community database dumps](https://github.com/boramalper/magnetico/issues/218).

The easiest way to get your own site up and running is to start with my .csv dump. It should be easy to import into any kind of system. It contains seed/leech counts too (!). If I were to import it, I'd modify import-magnetico-db.

__Torrent Paradise dump__ download: [MEGA](https://mega.nz/file/IFcTBCKZ#v3OCPNeja4lRC5baccVDeTaQUE150wqqGyS6A1mxglc) [BayFiles](https://bayfiles.com/nc43i3H5oe/torrentparadise-dump-200720.csv_xz)

# Usage

## Generate the index

See `update-index.sh`.

Generation of the IPFS index will prob take a long time, a machine with high single-core perf recommended (ipfsearch runs on node.js)

## Spider the DHT

Run `go build` in spider/ to compile and scp the binary it to the server. You can use the systemd service file in `spider/spider.service` to start the spider on server boot.

## Scraping trackers for seed/leech data

Run `go build` in seedleech-daemon/ to compile and scp the binary it to the server. You can use the systemd service file in `seedleech-daemon/seedleech.service`.

## Import a recent magnetico community dump

Use sqlite3 on a the decompressed dump to generate a .csv file. Format: infohash,name,length(bytes). Optionally quoted.

Then use the go binary in import-magnetico-db to do the import.

## IPFS vs 'static'

The directory website gets deployed to IPFS, static gets deployed to the server. Static calls the API, the IPFS version doesn't.

# Contributing

Before working on something, open an issue to ask if it would be okay. I would love to [KISS](https://en.wikipedia.org/wiki/KISS_principle). 
