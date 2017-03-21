fetch('/api/transactions/')
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
