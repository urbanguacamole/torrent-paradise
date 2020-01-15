app = new Vue({
    el: '#app',
    data: {
        results: undefined,
        resultsFound: false,
        showProgress: false,
        progress: 0
    }
})

let actionid = 0
let sessionid
let query

window.onmessage = function(e){
    if (e.data.type == "results") {
        let results = JSON.parse(e.data.results)
        app.results = results.map((result) => {
            result.length = formatBytes(result.len)
            return result
        })
        app.resultsFound = true
        setTimeout(updateSize,1)
    } else if (e.data.type == "query") {
        query = e.data.query
        console.log("Query sent, sending limited anonymized telemetry.")
        sendTelemetry({"query":query})
    }
};

function updateSize(){
    window.parent.postMessage(parseInt(document.body.scrollHeight),"*")
}

function formatBytes(a,b){if(0==a)return"0 Bytes";var c=1024,d=b||2,e=["B","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+" "+e[f]}

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

    fetch('/api/telemetry', {
        method: 'POST',
        body: JSON.stringify(payload)
    })
}

/**
 * Reports anonymized data about which result you picked. Smart result sorting is planned and I expect that screen resolution, language and OS will be critical in determining which result you prefer.
 */
function resultClicked(ih,name,s,l,len){
    console.log("Result clicked, sending limited anonymized telemetry")
    payload = {ih: ih, n: name, s: s, l:l, len: len}
    payload.os = platform.os
    payload.w = window.screen.width*window.devicePixelRatio
    payload.h = window.screen.height*window.devicePixelRatio
    payload.lang = navigator.language || navigator.userLanguage
    payload.query = query
    sendTelemetry(payload)
}
