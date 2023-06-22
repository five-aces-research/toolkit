package main

import "time"

type PnlStruct struct {
	Id      string
	Results []PnlResult
}

type Pnl struct {
	Time    time.Time
	Balance float64
}

type PnlResult struct {
	Name string
	Pnl  []Pnl
}

//htmlstart, pnls, chart, htmlend

var htmlstart = `<!DOCTYPE html>
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

var htmlend = `
</body>
</html>`

var pnlTemplate = `
    <div id="pnlContainer" style="height: 600px; width: 100%;"></div>
    <script>
    Highcharts.chart('pnlContainer', {
        chart: {
            type: 'spline'
        },
        title: {
            text: '{{.Name}}'
        },
        xAxis: {
            type: 'datetime',
            title: {
                text: 'Date'
            }
        },
        yAxis: {
            title: {
                text: 'Account Balance'
            },
            min: 0
        },
        tooltip: {
            headerFormat: '<b>{series.name}</b><br>',
            pointFormat: '{point.x:%e. %b}: {point.y:.2f} $'
        },
    
        plotOptions: {
            series: {
                marker: {
                    enabled: true,
                    radius: 2.5
                }
            }
        },
    
        colors: ['#FF0000', '#00FF00', '#0000FF', '#000000', '#FF00FF', '#00FFFF', '#FFA500', '#800080', '#008080', '#FFC0CB'],
    
        // Define the data points. All series have a year of 1970/71 in order
        // to be compared on the same x axis. Note that in JavaScript, months start
        // at 0 for January, 1 for February etc.
        series: [ {{range .Results}}
            {
                name: '{{.Name}}',
                data: [{{range .Pnl}}
                [{{.Time.UnixMilli}},{{.Balance}}],{{end}}    
            ]   
            },
            {{end}}
        ]
    });
    </script>
    `

var chartTemplate = `<div id="chartContainer" style="height: 600px; width: 100%;"></div>
<script>
    // Define the candlestick data
    var candlestickData = [{{range .}}
        [{{.StartTime.UnixMilli}},{{.Open}},{{.High}},{{.Low}},{{.Close}}],{{end}}];

    // Create the chart
    Highcharts.stockChart('chartContainer', {
        rangeSelector: {
            selected: 5
        },
        title: {
            text: 'Candlestick Chart'
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

var columnTemplate = `
<div id="columnContainer" style="height: 600px; width: 100%;"></div>
<script>
Highcharts.chart('columnContainer', {
    chart: {
        type: 'column'
    },
    title: {
        text: 'Trades Performance'
    },
    xAxis: {
        categories: [
            '-8<',
            '-8 - -6',
            '-6 - -4',
            '-4 - -2',
            '-2 - -0.5',
            '-0.5 - 0.5',
            '0.5 - 2',
            '2 - 4',
            '4 - 6',
            '6 - 8',
            '<8'
        ],
        crosshair: true
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Rainfall (mm)'
        }
    },
    tooltip: {
        headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
        pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
            '<td style="padding:0"><b>{point.y:.0f} </b></td></tr>',
        footerFormat: '</table>',
        shared: true,
        useHTML: true
    },
    plotOptions: {
        column: {
            pointPadding: 0.2,
            borderWidth: 0
        }
    },
    series: [
        {{range .}}
        {
        name: '{{.Name}}',
        data: [{{range .Distribution}}{{.}},{{end}}]

    },{{end}}]
});
</script>`
