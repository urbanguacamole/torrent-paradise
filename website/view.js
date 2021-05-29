app = new Vue({
    el: '#app',
    data: { showsearchbox: false, error: "", resultPage: "", resultPageHeight: 1, entries: -1 }
})
window.addEventListener("message", receiveMessage, false);

function receiveMessage(event) {
    app.resultPageHeight = event.data
}
searchbox = document.getElementById('searchbox')
if (searchbox != null) {
    searchbox.onkeydown = function (event) {
        if (event.keyCode == 13) {
            searchTriggered()
        }
    }
}

// If the URL has a hash (e.g. #ubuntu iso), search for it.
if (window.location.hash) {
    // `substr(1)` to trim off the leading `#`, `decodeURIComponent` to handle things like `%20` for ` `.
    searchbox.value = decodeURIComponent(window.location.hash.substr(1));
    searchTriggered();
}
