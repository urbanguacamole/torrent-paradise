# This script updates the index. Testing and uploading to server/IPFS is done manually.

echo "Scraping trackers for seed/leech data"
mosh nextgen@server "~/tracker-scraper" # you can use ssh instead of mosh aswell

ssh nextgen@server "psql -c 'REFRESH MATERIALIZED VIEW fresh'"

echo "Generating index dump"
rm index-generator/dump.csv
ssh nextgen@server "psql -c '\copy (select * from fresh) to stdout with (format csv)'" > index-generator/dump.csv

(cd index-generator; node --max-old-space-size=10000 main.js)
python3 index-generator/fix-metajson.py website/generated/inx

generate-top-torrents/generate-top-torrents > website/generated/top.json

echo "Uploading website"
cd website
rsync -ar ./ root@server:/www/torrent-paradise.ml # consider using --progress