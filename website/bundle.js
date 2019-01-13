class IndexFetcher {
    constructor() { this.combinedIndex = new Map(); this.shardsFetched = new Map(); }
    async fetchShard(shardid) {
        if (this.shardsFetched.has(shardid)) { console.debug("not needing to fetch shard " + shardid); return; }
        console.debug("started fetching inx shard " + shardid); this.shardsFetched.set(shardid, false); let shard = await loadIndexFromURL(meta.inxURLBase + shardid.toString()); for (let i of shard.keys()) {
            if (!inxFetcher.combinedIndex.has(i)) { inxFetcher.combinedIndex.set(i, shard.get(i)); }
            else { if (i != "") { console.warn("srsly weird"); } }
        }
        console.debug("shard " + shardid + " fetched!"); inxFetcher.shardsFetched.set(shardid, true);
    }
    getIndexFor(token) {
        let needle = 0; while (meta.inxsplits[needle] < token) { needle++; }
        if (needle !== 0) { return needle - 1; }
        else
            return needle;
    }
}
class InvertedIndexFetcher extends IndexFetcher {
    constructor() { super(...arguments); this.combinedInvIndex = new Map(); }
    async fetchShard(shardid) {
        if (this.shardsFetched.has(shardid)) { return; }
        console.debug("started fetching invinx shard " + shardid); this.shardsFetched.set(shardid, false); let shard = await loadInvertedIndexFromURL(meta.invURLBase + shardid.toString()); for (let i of shard.keys()) {
            if (!invinxFetcher.combinedInvIndex.has(i)) { invinxFetcher.combinedInvIndex.set(i, shard.get(i)); }
            else { if (i != "") { console.warn("srsly weird"); } }
        }
        console.debug("invinx shard " + shardid + " fetched!"); invinxFetcher.shardsFetched.set(shardid, true);
    }
    getIndexFor(token) {
        let needle = 0; while (meta.invsplits[needle] < token) { needle++; }
        if (needle !== 0) { return needle - 1; }
        else
            return needle;
    }
}
var inxFetcher = new IndexFetcher(); var invinxFetcher = new InvertedIndexFetcher(); var meta; var app; let ipfsGatewayURL; const NUMRESULTS = 60; function onLoad() {
    let params = new URLSearchParams(location.search); if (params.get("index")) { loadMeta(params.get("index")).then(function () { document.getElementById("app").style.visibility = ""; }); }
    else { document.getElementById("app").style.visibility = ""; }
}
async function loadMeta(metaURL) {
    let response; if (metaURL.startsWith("/ipfs/") || metaURL.startsWith("/ipns/")) { response = await fetch((await getIpfsGatewayUrlPrefix()) + metaURL); }
    else { response = await fetch(metaURL); }
    const json = await response.text(); try { meta = JSON.parse(json); }
    catch (e) { app.error = "Unable to find index at " + metaURL; return; }
    if (meta.invURLBase.startsWith("/ipfs/") || meta.invURLBase.startsWith("/ipns/")) { meta.invURLBase = (await getIpfsGatewayUrlPrefix()) + meta.invURLBase; }
    if (meta.inxURLBase.startsWith("/ipfs/") || meta.inxURLBase.startsWith("/ipns/")) { meta.inxURLBase = (await getIpfsGatewayUrlPrefix()) + meta.inxURLBase; }
    console.log("meta fetched"); app.showmeta = false; app.showsearchbox = true; app.indexAuthor = meta.author; app.indexName = meta.name; let ts = new Date(meta.created); app.indexTimestamp = ts.getDate().toString() + "/" + (ts.getMonth()+1).toString() + "/" + ts.getFullYear().toString(); if (meta.resultPage == undefined) { app.resultPage = "/basicresultpage"; }
    else {
        if (meta.resultPage.startsWith("/ipfs/") || meta.resultPage.startsWith("/ipns/")) { app.resultPage = (await getIpfsGatewayUrlPrefix()) + meta.resultPage; }
        else { app.resultPage = meta.resultPage; }
    }
}
async function getIpfsGatewayUrlPrefix() {
    if (ipfsGatewayURL !== undefined) { return ipfsGatewayURL; }
    if (window.location.protocol === "https:") {
        if (await checkIfIpfsGateway("")) { ipfsGatewayURL = window.location.protocol + "//" + window.location.host; }
        else { app.error = "ipfsearch is currently being served from a HTTPS host that is not an IPFS node. This prevents it from using a local IPFS gateway. The node operator should fix this and run an ipfs gateway."; }
    }
    else if (await checkIfIpfsGateway("http://localhost:8080")) { ipfsGatewayURL = "http://localhost:8080"; }
    else if (await checkIfIpfsGateway("http://" + window.location.host)) { ipfsGatewayURL = "http://" + window.location.host; }
    else { app.error = "Loading of the index requires access to the IPFS network. We have found no running IPFS daemon on localhost. Please install IPFS from <a href='http://ipfs.io/docs/install'>ipfs.io</a> and refresh this page."; throw new Error("Couldn't get a IPFS gateway."); }
    return ipfsGatewayURL;
}
async function checkIfIpfsGateway(gatewayURL) {
    let response = await fetch(gatewayURL + "/ipfs/QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"); if ((await response.text()).startsWith("hello world")) { return true; }
    else { return false; }
}
function searchTriggered() { let searchbox = document.getElementById("searchbox"); let querytokens = searchbox.value.split(" "); querytokens = querytokens.map(querytoken => { return stemmer(querytoken); }); console.debug("searching for: " + querytokens.join(" ")); searchFor(querytokens.join(" ")); }
function searchFor(query) {
    let runningFetches = []; let tokenizedquery = tokenize(query); tokenizedquery.forEach((token) => { runningFetches.push(invinxFetcher.fetchShard(invinxFetcher.getIndexFor(token))); }); Promise.all(runningFetches).then(() => {
        let candidates = getAllCandidates(tokenizedquery, invinxFetcher.combinedInvIndex); console.log("candidates prefilter: " + candidates.size); console.debug(candidates); candidates = filterCandidates(candidates, tokenizedquery.length); console.log("candidates postfilter: " + candidates.size); let resultIds; resultIds = []; let foundIdealCandidate; for (let key of candidates.keys()) {
            if (candidates.get(key) == tokenizedquery.length) { foundIdealCandidate = true; }
            resultIds.push(key);
        }
        console.debug(candidates); if (foundIdealCandidate) {
            console.info("Found an ideal candidate in prefetch sorting&filtering. Filtering out all non-ideal candidates..."); resultIds = resultIds.filter((resultId) => {
                if (candidates.get(resultId) != tokenizedquery.length) { return false; }
                else { return true; }
            });
        }
        else {
            console.debug(resultIds); resultIds = resultIds.sort((a, b) => {
                let ascore = candidates.get(a); let bscore = candidates.get(b); if (ascore > bscore) { return -1; }
                else if (ascore > bscore) { return 1; }
                else { return 0; }
            });
        }
        console.debug("resultIds after prefetch sorting & filtering: "); console.debug(resultIds); resultIds = resultIds.slice(0, NUMRESULTS); fetchAllDocumentsById(resultIds).then((results) => { passResultToResultpage(results); });
    });
}
function passResultToResultpage(results) { let resultPageIframe = document.getElementById("resultPage"); resultPageIframe.contentWindow.postMessage({ type: "results", results: JSON.stringify(results) }, '*'); }
async function fetchAllDocumentsById(ids) {
    let runningDocumentFetches; runningDocumentFetches = []; for (let id in ids) { runningDocumentFetches.push(getDocumentForId(ids[id])); }
    return Promise.all(runningDocumentFetches).then((results) => { return results; });
}
function filterCandidates(candidates, tokensInQuery) {
    if (tokensInQuery >= 2) {
        let filteredCandidates; filteredCandidates = new Map(); for (let key of candidates.keys()) { if (candidates.get(key) >= (tokensInQuery / 2)) { filteredCandidates.set(key, candidates.get(key)); } }
        candidates = undefined; return filteredCandidates;
    }
    else { return candidates; }
}
function getAllCandidates(query, index) {
    let candidates; candidates = new Map(); for (let i in query) {
        let result = index.get(query[i]); for (let j in result) {
            if (candidates.has(result[j])) { candidates.set(result[j], candidates.get(result[j]) + 1); }
            else { candidates.set(result[j], 1); }
        }
    }
    return candidates;
}
var stemmer = (function () {
    var step2list = { "ational": "ate", "tional": "tion", "enci": "ence", "anci": "ance", "izer": "ize", "bli": "ble", "alli": "al", "entli": "ent", "eli": "e", "ousli": "ous", "ization": "ize", "ation": "ate", "ator": "ate", "alism": "al", "iveness": "ive", "fulness": "ful", "ousness": "ous", "aliti": "al", "iviti": "ive", "biliti": "ble", "logi": "log" }, step3list = { "icate": "ic", "ative": "", "alize": "al", "iciti": "ic", "ical": "ic", "ful": "", "ness": "" }, c = "[^aeiou]", v = "[aeiouy]", C = c + "[^aeiouy]*", V = v + "[aeiou]*", mgr0 = "^(" + C + ")?" + V + C, meq1 = "^(" + C + ")?" + V + C + "(" + V + ")?$", mgr1 = "^(" + C + ")?" + V + C + V + C, s_v = "^(" + C + ")?" + v; return function (w) {
        var stem, suffix, firstch, re, re2, re3, re4, origword = w; if (w.length < 3) { return w; }
        firstch = w.substr(0, 1); if (firstch == "y") { w = firstch.toUpperCase() + w.substr(1); }
        re = /^(.+?)(ss|i)es$/; re2 = /^(.+?)([^s])s$/; if (re.test(w)) { w = w.replace(re, "$1$2"); }
        else if (re2.test(w)) { w = w.replace(re2, "$1$2"); }
        re = /^(.+?)eed$/; re2 = /^(.+?)(ed|ing)$/; if (re.test(w)) { var fp = re.exec(w); re = new RegExp(mgr0); if (re.test(fp[1])) { re = /.$/; w = w.replace(re, ""); } }
        else if (re2.test(w)) {
            var fp = re2.exec(w); stem = fp[1]; re2 = new RegExp(s_v); if (re2.test(stem)) {
                w = stem; re2 = /(at|bl|iz)$/; re3 = new RegExp("([^aeiouylsz])\\1$"); re4 = new RegExp("^" + C + v + "[^aeiouwxy]$"); if (re2.test(w)) { w = w + "e"; }
                else if (re3.test(w)) { re = /.$/; w = w.replace(re, ""); }
                else if (re4.test(w)) { w = w + "e"; }
            }
        }
        re = /^(.+?)y$/; if (re.test(w)) { var fp = re.exec(w); stem = fp[1]; re = new RegExp(s_v); if (re.test(stem)) { w = stem + "i"; } }
        re = /^(.+?)(ational|tional|enci|anci|izer|bli|alli|entli|eli|ousli|ization|ation|ator|alism|iveness|fulness|ousness|aliti|iviti|biliti|logi)$/; if (re.test(w)) { var fp = re.exec(w); stem = fp[1]; suffix = fp[2]; re = new RegExp(mgr0); if (re.test(stem)) { w = stem + step2list[suffix]; } }
        re = /^(.+?)(icate|ative|alize|iciti|ical|ful|ness)$/; if (re.test(w)) { var fp = re.exec(w); stem = fp[1]; suffix = fp[2]; re = new RegExp(mgr0); if (re.test(stem)) { w = stem + step3list[suffix]; } }
        re = /^(.+?)(al|ance|ence|er|ic|able|ible|ant|ement|ment|ent|ou|ism|ate|iti|ous|ive|ize)$/; re2 = /^(.+?)(s|t)(ion)$/; if (re.test(w)) { var fp = re.exec(w); stem = fp[1]; re = new RegExp(mgr1); if (re.test(stem)) { w = stem; } }
        else if (re2.test(w)) { var fp = re2.exec(w); stem = fp[1] + fp[2]; re2 = new RegExp(mgr1); if (re2.test(stem)) { w = stem; } }
        re = /^(.+?)e$/; if (re.test(w)) { var fp = re.exec(w); stem = fp[1]; re = new RegExp(mgr1); re2 = new RegExp(meq1); re3 = new RegExp("^" + C + v + "[^aeiouwxy]$"); if (re.test(stem) || (re2.test(stem) && !(re3.test(stem)))) { w = stem; } }
        re = /ll$/; re2 = new RegExp(mgr1); if (re.test(w) && re2.test(w)) { re = /.$/; w = w.replace(re, ""); }
        if (firstch == "y") { w = firstch.toLowerCase() + w.substr(1); }
        return w;
    };
})(); function tokenize(name) { let tokens = name.split(' ').join(',').split('.').join(',').split('(').join(',').split(')').join(',').split('-').join(',').split('_').join(',').split(','); tokens.forEach((value, index, array) => { array[index] = value.toLowerCase(); }); tokens.forEach((value, index, array) => { array[index] = stemmer(value); }); return tokens; }
async function loadIndexFromURL(url) {
    let response = await fetch(url); let responsetext; if (response.ok) { responsetext = await response.text(); }
    else { throw new Error(response.statusText); }
    let parsedResponse = JSON.parse(responsetext); let parsedIndex; parsedIndex = new Map(); for (let object of parsedResponse) {
        let id = ""; let document = new Object(); for (let property of Object.keys(object)) {
            if (property === "id") { id = object[property]; }
            else { document[property] = object[property]; }
        }
        parsedIndex.set(id, document);
    }
    return parsedIndex;
}
function loadInvertedIndexFromURL(url) {
    return fetch(url).then((response) => {
        if (response.ok) { return response.text(); }
        throw new Error("Couldn't fetch shard at URL " + url);
    }).then((response) => {
        let loadedIndex = new Map(); let lineNumber = 0; let lines = response.split("\n"); let version; lines.forEach((line) => {
            if (lineNumber === 0) {
                if (parseInt(line) != 1 && parseInt(line) != 2) { throw "Error while parsing invinx: Invalid version, must be 1 or 2!"; }
                else { version = parseInt(line); }
                lineNumber++; return;
            }
            let cols = line.split(","); let tokenname = decodeURIComponent(cols[0]); cols.shift(); if (version === 2) { cols = cols.map(function (value) { return value.replace("%2C", ","); }); }
            loadedIndex.set(tokenname, cols); lineNumber++;
        }); return (loadedIndex);
    });
}
async function getDocumentForId(docid) {
    docid = docid.replace("%2C", ","); await inxFetcher.fetchShard(inxFetcher.getIndexFor(docid)); if (inxFetcher.combinedIndex.get(docid) === undefined) { console.error("No document found for docid " + docid); return { text: "no document found", id: docid }; }
    let doc = inxFetcher.combinedIndex.get(docid); doc["id"] = docid; return inxFetcher.combinedIndex.get(docid);
}