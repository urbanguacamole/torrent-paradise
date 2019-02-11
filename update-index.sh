# This script updates the index. Testing and uploading to server/IPFS is done manually.

echo "Scraping trackers for seed/leech data"
mosh nextgen@server "~/tracker-scraper"
echo "Generating SQL dump"
ssh nextgen@server pg_dump --data-only --inserts nextgen > index-generator/dump.sql

sed -i -e 's/public.peercount/peercount/g' index-generator/dump.sql
sed -i -e 's/public.torrent/torrent/g' index-generator/dump.sql
tail -n +21 index-generator/dump.sql > index-generator/newdump.sql # remove headers
mv index-generator/newdump.sql index-generator/dump.sql
rm index-generator/db.sqlite3
echo "Preparing sqlite DB"
sqlite3 index-generator/db.sqlite3 "CREATE TABLE peercount ( infohash char(40), tracker varchar, seeders int, leechers int, completed int, scraped timestamp, ws boolean);"
sqlite3 index-generator/db.sqlite3 "CREATE TABLE torrent( infohash char(40), name varchar, length bigint, added timestamp);"
echo """Do the following: 
$ sqlite3 index-generator/db.sqlite3

sqlite> BEGIN;
sqlite> .read index-generator/dump.sql
sqlite> END;"""
bash
echo "Generating index now..."
(cd index-generator; node --max-old-space-size=10000 main.js)
echo "Check meta.json, add resultPage:'resultpage', fix invURLBase, inxURLBase"
nano website/generated/inx.meta.json
echo "Uploading website"
cd website
scp -r . user@server:/www/torrent-paradise.ml
echo "Finished uploading website to server."