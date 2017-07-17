window.addEventListener('load', () => {
    fetch_purchases("");
    document.getElementById('query_form').addEventListener('submit', () => {
        fetch_purchases(document.getElementById('query_box').value);
        return false;
    });
});

function fetch_purchases(query) {
    fetch(`/api/purchase/?query=${encodeURI(query)}`)
        .then((response) => {
            if (!response.ok) {
                console.log(`transactions [GET] API returns error: code=${response.status}, msg=${response.statusText}`)
                return;
            }
            response.json().then((list) => {
                riot.mount('transactions', { items: list})
            });
        })
        .catch((error) => {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
}