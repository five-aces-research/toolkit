package tpl

type Stats struct {
	Name    string
	Number  int
	Mean    float64
	Stdev   float64
	Highest float64
	Lowest  float64
}

type Column struct {
	Id           string
	Categories   []string
	Height       string
	Width        string
	Distribution []Distribution
	Stats        []Stats
}

type Distribution struct {
	Name   string
	Values []int
}

const DEFAULTHEIGHT = "600px"
const DEFAULTWIDTH = "100%"

var DEFAULTCATEGORIES = []string{
	"-8 - -6",
	"-6 - -4",
	"-4 - -2",
	"-2 - -0.5",
	"-0.5 - 0.5",
	"0.5 - 2",
	"2 - 4",
	"4 - 6",
	"6 - 8",
	"<8"}

var ColumnHTML = `
<div id="column{{.Id}}" style="height: {{.Height}}; width: {{.Width}};"></div>
<script>
Highcharts.chart('column{{.Id}}', {
    chart: {
        type: 'column'
    },
    title: {
        text: 'Trades Performance'
    },
    xAxis: {
        categories: [
			{{range .Categories}}'{{.}}',
			{{end}}
        ],
        crosshair: true
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Number of Trades'
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
        {{range .Distribution}}
        {
        name: '{{.Name}}',
        data: [{{range .Values}}{{.}},{{end}}]

    },{{end}}]
});

</script>
<table class="table table-dark w-100">
  <thead>
    <tr>
      <th>Algo</th>
      <th>Number of Trades</th>
      <th>Mean</th>
      <th>Stdev</th>
      <th>Max Loss</th>
      <th>Max Win</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>{{.Name}}</td>
      <td>{{.Number}}</td>
      <td>{{.Mean}}</td>
      <td>{{.Stdev}}</td>
      <td>{{.Highest}}</td>
      <td>{{.Lowest}}</td>
    </tr>
  </tbody>
</table>

`
