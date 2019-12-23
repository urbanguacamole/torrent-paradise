function searchTriggered() {
    let searchbox = document.getElementById("searchbox");
    let query = searchbox.value
    searchFor(query);
}

async function searchFor(query) {
    var url = new URL("https://torrent-paradise.ml/api/search")
    url.searchParams.append("q",query)
    const res = await fetch(url)
    let results = await res.json();
    if (results == null){
        console.error("No results.")
        results = []
    }
    passResultToResultpage(results)
}

function passResultToResultpage(results) {
    let resultPageIframe = document.getElementById("resultPage");
    resultPageIframe.contentWindow.postMessage({
        type: "results",
        results: JSON.stringify(results)
    }, '*');
}
/**
 * Sends telemetry payload, adds actionid and sessionid to it. IP is never logged.
 */
function sendTelemetry(payload){
    payload.aid = actionid;
    actionid = actionid + 1
    if (sessionid == undefined){
        sessionid = Math.round((Math.random()-0.5)*Math.pow(2,32))
        payload.sid = sessionid;
    }else{
        payload.sid = sessionid;
    }
    

    (async (payload) => {
        await fetch('https://torrent-paradise.ml/api/telemetry', {
            method: 'POST',
            body: JSON.stringify(payload)
        })
    })(payload);
}