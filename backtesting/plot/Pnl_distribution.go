package plot

import (
	"errors"
	"log"
	"math"
	"os"
	"text/template"

	"github.com/five-aces-research/toolkit/backtesting/plot/tpl"
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/google/uuid"
)

func DrawPnlDistributionColumn(Filename string, MainStrat []*strategy.Trade, backtests []*strategy.BackTestStrategy) error {
	data := pnlDistribution(MainStrat, backtests)

	ff, err := os.Create(Filename)
	if err != nil {
		return err
	}

	_, err = ff.WriteString(tpl.HtmlStart)
	if err != nil {
		return err
	}

	temp, err := template.New("columen").Parse(tpl.ColumnHTML)
	if err != nil {
		return err
	}
	if err := temp.Execute(ff, data); err != nil {
		return err
	}

	_, err = ff.WriteString(tpl.HtmlEnd)
	return err
}

func pnlDistribution(MainStrat []*strategy.Trade, filtered []*strategy.BackTestStrategy) tpl.Column {
	// Create []float64
	// Standardabweichung
	var data tpl.Column

	if MainStrat != nil {
		data.Stats = append(data.Stats, Stdev("main", MainStrat))
	}

	for _, v := range filtered {
		st, dist, err := CreateStatsAndDistribution(v.Name, v.Trade())
		if err != nil {
			log.Println(v.Name, err)
			continue
		}
		data.Stats = append(data.Stats, st)
		data.Distribution = append(data.Distribution, dist)
	}

	data.Categories = tpl.DEFAULTCATEGORIES
	data.Height = tpl.DEFAULTHEIGHT
	data.Width = tpl.DEFAULTWIDTH
	data.Id = uuid.NewString()

	return data
}

func Stdev(name string, tr []*strategy.Trade) tpl.Stats {
	var sum, max, min float64
	stdev := make([]float64, 0, len(tr))

	if len(tr) == 0 || tr == nil {
		return tpl.Stats{}
	}
	max = tr[0].PnlPercent()
	min = max
	for _, t := range tr {
		value := t.PnlPercent()
		sum += value
		stdev = append(stdev, value)
		if value > max {
			max = value
		} else if value < min {
			min = value
		}
	}
	mean := sum / float64(len(stdev))
	variance := calculateVariance(stdev, mean)
	stdDeviate := math.Sqrt(variance)

	return tpl.Stats{Name: name, Number: len(stdev), Mean: mean, Stdev: stdDeviate, Highest: max, Lowest: min}
}

func CreateStatsAndDistribution(name string, tr []*strategy.Trade) (tpl.Stats, tpl.Distribution, error) {
	var sum, max, min float64
	stdev := make([]float64, 0, len(tr))
	distribution := make([]int, 11, 11)
	if len(tr) == 1 || tr == nil {
		return tpl.Stats{}, tpl.Distribution{}, errors.New("no trades")
	}

	max = tr[0].PnlPercent()
	min = max

	for _, t := range tr {
		value := t.PnlPercent()
		sum += value
		stdev = append(stdev, value)
		if value > max {
			max = value
		} else if value < min {
			min = value
		}
		switch {
		case value < -8:
			distribution[0]++
		case value >= -8 && value < -6:
			distribution[1]++
		case value >= -6 && value < -4:
			distribution[2]++
		case value >= -4 && value < -2:
			distribution[3]++
		case value >= -2 && value < -0.5:
			distribution[4]++
		case value >= -0.5 && value < 0.5:
			distribution[5]++
		case value >= 0.5 && value < 2:
			distribution[6]++
		case value >= 2 && value < 4:
			distribution[7]++
		case value >= 4 && value < 6:
			distribution[8]++
		case value >= 6 && value < 8:
			distribution[9]++
		case value >= 8:
			distribution[10]++
		default:
		}
	}
	mean := sum / float64(len(stdev))
	variance := calculateVariance(stdev, mean)
	stdDeviate := math.Sqrt(variance)

	return tpl.Stats{Name: name, Number: len(stdev), Mean: mean, Stdev: stdDeviate, Highest: max, Lowest: min}, tpl.Distribution{Name: name, Values: distribution}, nil
}

func calculateVariance(data []float64, mean float64) float64 {
	squaredDiffSum := 0.0
	for _, value := range data {
		diff := value - mean
		squaredDiffSum += diff * diff
	}
	return squaredDiffSum / float64(len(data))
}
