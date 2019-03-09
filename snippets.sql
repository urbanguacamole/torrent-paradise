-- SELECT for the index-ready materialized view "fresh"
select distinct on (infohash) trackerdata.infohash, torrent.name, torrent.length, trackerdata.seeders, trackerdata.leechers, trackerdata.completed, tracker from trackerdata inner join torrent on (trackerdata.infohash = torrent.infohash) where torrent.copyrighted != 't' order by infohash, scraped asc, seeders desc;
-- to create this view, run: create materialized view AS (paste from above)

select sum(seeders),tracker from trackerdata group by tracker;

get highest seed/leech count found for every torrent, using data from tracker with most seeds
select distinct on (infohash) infohash, seeders, leechers, completed, tracker from trackerdata where completed != 0 order by infohash, scraped asc, seeders desc;



generate top 100 by seeders:
select * from (select distinct on (trackerdata.infohash) trackerdata.infohash, torrent.name, seeders, leechers from trackerdata inner join torrent on (trackerdata.infohash = torrent.infohash) where completed != 0 order by infohash, seeders desc) as subquery order by seeders desc limit 100;

SELECT added::date, count(infohash)
from torrent where added > '2019-01-15'::date
group by added::date order by count desc;