package tpl

import "time"

type Pnl struct {
	Id   string
	Data []PnlResult
}

type PnlResult struct {
	Name string
	Pnl  []PnlPoint
}

type PnlPoint struct {
	Time    time.Time
	Balance float64
}

var PnlHTML = `
    <div id="pnl{{.Id}}" style="height: 600px; width: 100%;"></div>
    <script>
    Highcharts.chart('pnl{{.Id}}', {
        chart: {
            type: 'spline'
        },
        title: {
            text: 'PNL'
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
        series: [ {{range .Data}}
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
