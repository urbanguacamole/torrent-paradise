go build
scp ./spider user@server:/home/nextgen/
ssh user@server systemctl restart spider