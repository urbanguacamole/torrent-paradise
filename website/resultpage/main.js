app = new Vue({
    el: '#app',
    data: {
        results: undefined,
        resultsFound: false,
        showProgress: false,
        progress: 0
    }
})

window.onmessage = function(e){
    if (e.data.type == "results") {
        let results = JSON.parse(e.data.results)
        app.results = results.map((result) => {
           result.length = formatBytes(result.len)
           return result
        })
        app.resultsFound = true
        setTimeout(updateSize,1)
    } else if (e.data.type == "progress") {
        if(e.data.progress == 1){
            app.showProgress = false
        }else{
            app.showProgress = true
        }
        app.progress = e.data.progress * 100
        setTimeout(updateSize,1)
    }
};

function updateSize(){
    window.parent.postMessage(parseInt(document.body.scrollHeight),"*")
}

function formatBytes(a,b){if(0==a)return"0 Bytes";var c=1024,d=b||2,e=["B","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+" "+e[f]}