go build
scp ./tracker-scraper user@server:/home/nextgen/
ssh user@server sudo -u nextgen /home/nextgen/tracker-scraper