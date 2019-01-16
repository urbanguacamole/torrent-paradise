go build
ssh user@server systemctl stop spider
scp ./spider user@server:/home/nextgen/
ssh user@server systemctl start spider