package main

import (
	"os"
	"text/template"

	"github.com/five-aces-research/toolkit/fas"
)

func DrawFunction(filename string, chart []fas.Candle, results *PnlStruct) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = file.WriteString(htmlstart)
	if err != nil {
		return nil
	}

	if chart != nil {
		tpl, err := template.New("chart").Parse(chartTemplate)
		if err != nil {
			return nil
		}

		tpl.Execute(file, chart)
	}

	if results != nil {
		tpl, err := template.New("pnl").Parse(pnlTemplate)
		if err != nil {
			return nil
		}

		tpl.Execute(file, *results)
	}

	_, err = file.WriteString(htmlend)
	return err
}
