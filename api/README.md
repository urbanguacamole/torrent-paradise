Expecting a future where clickstream data are more and more important, because ranking will be moving from strict keyword based seeders sorted to fuzzy matching with machine learned sorting, I'll be adding limited anonymized telemetry to Torrent Paradise.

## Telemetry collected after searching
- query
- session id
    - required to find out if this is a followup query (after not being able to find results with a previous query) or a first
    - pseudorandom identifier generated on each visit, not persisted across browser tabs
- action #
    - needed to find the order of actions

## Telemetry collected after clicking magnet link
- infohash
- query (once again)
- torrent name
- seeders
- leechers
- torrent size
- session id
- action #
- rank #
- screen resolution
    - expected to be critical for ranking of video
- User Agent
    - expected to be critical for ranking of software
- language
    - could help for ranking video, might predict internet speed and thus required file size


Your IP address is never logged. Telemetry is always sent encrypted. HTTPS is terminated at the caddy web server.