function searchTriggered() {
    let searchbox = document.getElementById("searchbox");
    let query = searchbox.value
    searchFor(query);
    passQueryToResultpage(query)
}

async function searchFor(query) {
    var url = new URL(location.origin+"/api/search")
    url.searchParams.append("q",query)
    const res = await fetch(url)
    let results = await res.json();
    if (results == null){
        console.error("No results.")
        alert("No results found.")
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

function passQueryToResultpage(query) {
    let resultPageIframe = document.getElementById("resultPage");
    resultPageIframe.contentWindow.postMessage({
        type: "query",
        query: query
    }, '*');
}