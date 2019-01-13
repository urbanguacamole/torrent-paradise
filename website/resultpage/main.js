app = new Vue({
    el: '#app',
    data: {
        results: undefined,
        resultsFound: false
    }
})

window.onmessage = function(e){
    if (e.data.type == "results") {
        let results = JSON.parse(e.data.results)
        results = results.sort((a,b) => {
            if(a.s > b.s){
                return -1;
            }else if(a.s == b.s){
                return 0;
            }else{
                return 1;
            }
        })
        app.results = results.map((result) => {
            result.len = formatBytes(result.len)
            return result
        })
        app.resultsFound = true
        setTimeout(updateSize,1)
    }
};

function updateSize(){
    window.parent.postMessage(parseInt(document.body.scrollHeight),"*")
}

function formatBytes(a,b){if(0==a)return"0 Bytes";var c=1024,d=b||2,e=["B","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+" "+e[f]}