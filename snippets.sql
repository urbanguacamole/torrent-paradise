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