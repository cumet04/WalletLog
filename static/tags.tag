<transactions>
    <table border=1>
        <tr>
            <th>time</th>
            <th>price</th>
            <th>content</th>
            <th>source</th>
        </tr>
        <tr each={ items }>
            <td>{ time }</td>
            <td align='right'>{ price }</td>
            <td>{ content }</td>
            <td>{ source }</td>
        </tr>
    </table>

    <style>
        :scope { border-collapse: collapse }
    </style>

    <script>
    this.items = opts.items.map((x) => {
        x.time = x.time.replace(/T00:00:00\+09:00/, '');
        x.price = String(x.price).replace(/(\d)(?=(\d\d\d)+(?!\d))/g, '$1,');
        return x;
    });
    </script>
</transactions>