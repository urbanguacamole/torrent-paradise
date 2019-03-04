const ipfsearch = require("ipfsearch-index");
const parse = require("csv-parse");
const fs = require("fs");
let indexer = new ipfsearch.Indexer();
let i = 0;
const parser = parse();
fs.createReadStream("dump.csv").pipe(parser);
parser.on('readable', function () {
    let record;
    while (record = parser.read()) {
        if (parseInt(record[3]) > 0) {
            indexer.addToIndex(new Torrent(record[0], record[1], parseInt(record[2]), parseInt(record[3]), parseInt(record[4]), parseInt(record[5])));
            i++;
        }
    }
});
parser.on('error', function (err) {
    console.error(err.message);
});
parser.on('end', function () {
    console.log("Read all " + i + " records. Persisting.");
    indexer.persist("../website/generated/inv", "../website/generated/inx", "Urban Guacamole", "Torrent Paradise index", "", 1000);
});
class Torrent extends ipfsearch.Document {
    constructor(id, text, size, seeders, leechers, completed) {
        super(id, text);
        this.len = size;
        this.s = seeders;
        this.l = leechers;
        this.c = completed;
    }
}
