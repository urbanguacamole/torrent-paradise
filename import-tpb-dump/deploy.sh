go build
scp ./import-tpb-dump user@server:/home/nextgen/
ssh user@server sudo -u nextgen /home/nextgen/import-tpb-dump