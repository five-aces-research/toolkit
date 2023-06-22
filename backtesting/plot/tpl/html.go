package tpl

var HtmlStart = `<!DOCTYPE html>
<html>
<head>
    <title>Candlestick Chart</title>
    <!-- Include Highcharts library -->
    <script src="https://code.highcharts.com/stock/highstock.js"></script>
    <script src="https://code.highcharts.com/highcharts.js"></script>
    <!-- Include Highstock module for candlestick charts -->
    <script src="https://code.highcharts.com/modules/series-label.js"></script>
    <script src="https://code.highcharts.com/stock/modules/accessibility.js"></script>
    <script src="https://code.highcharts.com/stock/modules/data.js"></script>
<script src="https://code.highcharts.com/stock/modules/exporting.js"></script>
    </head>
<body>`

var HtmlEnd = `
</body>
</html>`

//    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
