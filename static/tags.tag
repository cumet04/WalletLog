<transactions>
    <table>
        <tr>
            <th>ID</th>
            <th>type</th>
            <th>time</th>
            <th>price</th>
            <th>content</th>
        </tr>
        <tr each={ items }>
            <td>{ id }</td>
            <td>{ type }</td>
            <td>{ time }</td>
            <td>{ price }</td>
            <td>{ content }</td>
        </tr>
    </table>

    <script>
    var self = this;
    fetch('/api/transactions/')
        .then((response) => {
            if (!response.ok) {
                console.log(`transactions [GET] API returns error: code=${response.status}, msg=${response.statusText}`)
                return;
            }
            response.json().then((list) => {
                self.items = list;
                self.update();
            });
        })
        .catch((error) => {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
    </script>
</transactions>