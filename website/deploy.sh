scp -r . user@server:/www/torrent-paradise.ml
ssh user@server "cd /www/torrent-paradise.ml; cat adsnippet >> index.html; cat adsnippet >> about.html; cat adsnippet >> ipfs.html"