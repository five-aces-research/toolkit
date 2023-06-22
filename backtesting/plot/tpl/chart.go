package tpl

import "github.com/five-aces-research/toolkit/fas"

type Chart struct {
	Id     string
	Ticker string
	Data   []fas.Candle
}

var ChartHtml = `<div id="chart{{.Id}}" style="height: 600px; width: 100%;"></div>
<script>
    // Define the candlestick data
    var candlestickData = [{{range .Data}}
        [{{.StartTime.UnixMilli}},{{.Open}},{{.High}},{{.Low}},{{.Close}}],{{end}}];

    // Create the chart
    Highcharts.stockChart('chart{{.Id}}', {
        rangeSelector: {
            selected: 5
        },
        title: {
            text: '{{.Ticker}}'
        },
        series: [{
            type: 'candlestick',
            name: 'Candlestick Data',
            data: candlestickData,
            tooltip: {
            }
        }
    ]
    });
</script>`
