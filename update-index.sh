# This script updates the index. Testing and uploading to server/IPFS is done manually.

echo "Refreshing database"
ssh nextgen@server "psql -c 'REFRESH MATERIALIZED VIEW fresh'"
echo "Downloading dump"
rm index-generator/dump.csv
ssh nextgen@server "psql -c '\copy (select fresh.infohash, torrent.name, torrent.length, fresh.s, fresh.l, fresh.c from fresh inner join torrent on torrent.infohash = fresh.infohash) to stdout with (format csv)'" > index-generator/dump.csv
echo "Generating index"
(cd index-generator; node --max-old-space-size=10000 main.js)
python3 index-generator/fix-metajson.py website/generated/inx
echo "Generating top torrents list"
generate-top-torrents/generate-top-torrents > website/generated/top.json
echo "Uploading website"
cd website
rsync -ar ./ root@server:/www/torrent-paradise.ml # consider using --progress