app = new Vue({
    el: '#app',
    data: { showsearchbox: false, error: "", resultPage: "", resultPageHeight: 1, entries: -1 }
})
window.addEventListener("message", receiveMessage, false);

function receiveMessage(event) {
    app.resultPageHeight = event.data
}

const searchbox = document.getElementById('searchbox')
if (searchbox != null) {
    searchbox.onkeydown = function (event) {
        if (event.keyCode == 13) {
            searchTriggered()
        }
    }
}

const urlParams = new URLSearchParams(window.location.search);
const query = urlParams.get('q');
if (query != null) {
    searchbox.value = query
    searchTriggered()
}
