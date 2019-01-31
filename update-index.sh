# This script updates the index and pushes it to IPFS. Should be run often.

echo "Scraping trackers for seed/leech data"
mosh nextgen@dev.ipfsearch.xyz "~/tracker-scraper"
echo "Generating SQL dump"
ssh nextgen@dev.ipfsearch.xyz pg_dump --data-only --inserts nextgen > index-generator/dump.sql

sed -i -e 's/public.peercount/peercount/g' index-generator/dump.sql
sed -i -e 's/public.torrent/torrent/g' index-generator/dump.sql
tail -n +21 index-generator/dump.sql > index-generator/newdump.sql # remove headers
mv index-generator/newdump.sql index-generator/dump.sql
rm index-generator/db.sqlite3
echo """Do the following: 
$ sqlite3 index-generator/db.sqlite3

sqlite> CREATE TABLE peercount ( infohash char(40), tracker varchar, seeders int, leechers int, completed int, scraped timestamp, ws boolean);
sqlite> CREATE TABLE torrent( infohash char(40), name varchar, length bigint, added timestamp);
sqlite> BEGIN;
sqlite> .read index-generator/dump.sql
sqlite> END;"""
bash
echo "Generating index now..."
cd index-generator
node --max-old-space-size=10000 main.js
cd ..
echo "Check meta.json, add resultPage='resultpage', fix invURLBase, inxURLBase"
nano website/generated/inx.meta.json
echo "Uploading website"
cd website
scp -r . user@server:/www/torrent-paradise.ml
echo "Finished uploading website to server. Adding to IPFS"
ssh user@server sudo -u ipfs ipfs add -r /www/torrent-paradise.ml/
echo "Check if it works, maybe publish to IPNS."