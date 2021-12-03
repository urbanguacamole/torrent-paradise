

# Run as root

Assumming Debian 11 Bullseye

```apt update
apt upgrade
apt install -y htop kitty-terminfo screenfetch postgresql-13 mosh nload
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | tee /etc/apt/trusted.gpg.d/caddy-stable.asc
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt install caddy

mkdir /www
adduser nextgen
```

### /etc/caddy/Caddyfile

```
http://torrentparadise.ml, http://torrent-paradise.ml {
        reverse_proxy /api/* http://localhost:8000
        root * /www/torrent-paradise.ml/
        file_server
}
```

### /etc/postgresql/13/main/postgresql.conf

Use https://pgtune.leopard.in.ua/

Just get the optimal settings and paste them at the end of the file, they override the defaults.

### Set up nextgen user and database in Postgres

```
postgres $ createuser -d nextgen
nextgen $ createdb nextgen
```

### Ship compiled static executables to server

You can either build it on the server or just ship the binaries to the server via scp. In the end, you need binaries in /home/nextgen and .service files in /etc/systemd/system/.

Might come in handy: a way to build go binaries truly statically (incl glibc) `CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .`

### Ship contents of static/ to /www/torrent-paradise.ml

Use scp.