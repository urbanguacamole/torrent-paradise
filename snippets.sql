-- see what trackers are the most used
select sum(seeders),tracker from trackerdata group by tracker;

-- generate top 100 by seeders:
SELECT torrent.name, fresh.* from fresh INNER JOIN torrent ON torrent.infohash = fresh.infohash ORDER BY s desc limit 100;

SELECT added::date, count(infohash)
from torrent where added > '2019-01-15'::date
group by added::date order by count desc;

CREATE MATERIALIZED VIEW fresh AS 
 SELECT infohash,
    max(seeders) AS s,
    max(leechers) AS l,
    max(completed) AS c
   FROM trackerdata
  GROUP BY infohash;

--- size of table
SELECT pg_size_pretty(pg_total_relation_size('"<schema>"."<table>"'));

--- count rows
SELECT reltuples::bigint AS estimate FROM pg_class where relname='mytable';

--- create fulltext table
CREATE TABLE search AS select torrent.*, fresh.s as s, fresh.l as l, to_tsvector(torrent.name) as vect from torrent inner join fresh on fresh.infohash = torrent.infohash;