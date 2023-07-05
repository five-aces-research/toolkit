package plot

import (
	"errors"
	"log"
	"os"
	"text/template"

	"github.com/five-aces-research/toolkit/backtesting/plot/tpl"
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

/*
Plot has Multiple Usecase. Its main usecase is to provide data visualisation for Trading Results
But it also supports Charting OHCLV Data, together with On Chart Indicator Like SMA or Bolling Bands
PseudoCode
plot.NewChart(ch)

*/

func SimplePnl(Filename string, chart ta.Chart, balance float64, results []*strategy.BackTestStrategy) error {

	ff, err := os.Create(Filename)
	if err != nil {
		return err
	}

	_, err = ff.WriteString(tpl.HtmlStart)
	if err != nil {
		return err
	}

	if chart != nil {
		templ, err := template.New("chart").Parse(tpl.ChartHtml)
		if err != nil {
			return err
		}
		templ.Execute(ff, tpl.Chart{Id: "chart", Ticker: chart.Name(), Data: chart.Data()})
	}

	if results != nil {
		var pnl tpl.Pnl
		for _, v := range results {
			vv, err := getPnlInfo(balance, v)
			if err != nil {
				log.Println(err)
				continue
			}

			pnl.Data = append(pnl.Data, *vv)
			pnl.Id = "pnl"
		}
		templ, err := template.New("chart").Parse(tpl.PnlHTML)
		if err != nil {
			return err
		}
		err = templ.Execute(ff, pnl)
		if err != nil {
			return err
		}
	}

	_, err = ff.WriteString(tpl.HtmlEnd)
	return err
}

func getPnlInfo(balance float64, b *strategy.BackTestStrategy) (*tpl.PnlResult, error) {
	if len(b.Trade()) == 0 {
		return nil, errors.New("no trades")
	}
	var pr tpl.PnlResult
	pp := make([]tpl.PnlPoint, 0, len(b.Trade())+1)
	pp = append(pp, tpl.PnlPoint{Balance: balance, Time: b.Trade()[0].EntrySignalTime})

	for _, v := range b.Trade() {
		balance += v.RealisedPNL()
		pp = append(pp, tpl.PnlPoint{Time: v.CloseSignalTime, Balance: balance})
	}

	pr.Pnl = pp
	pr.Name = b.Name
	return &pr, nil
}
