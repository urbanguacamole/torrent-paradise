app = new Vue({
    el: '#app',
    data: { showsearchbox: false, error: "", resultPage: "resultpage/", resultPageHeight: 1, entries: -1}
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