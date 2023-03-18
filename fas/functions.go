package fas

import "fmt"

/*
GenerateResolutionFunc returns a function
Usually exchanges only support specific resolutions. like 24h,4h,1h,30min
If you want to have a different resolution this function together with ConvertChartResolution
converts you the resolution you want.
To get a function that converts the resolution for you exchange, add the
supported resolution in desc order e.g. GenerateResolutionFunc(86400,14400,3600,900,300,60,15)
*/
func GenerateResolutionFunc(SupportedResInMinutes ...int64) func(int64) int64 {
	return func(r int64) int64 {
		var newRes int64
		for _, v := range SupportedResInMinutes {
			if r == v {
				newRes = r
				return newRes
			}
		}
		for _, v := range SupportedResInMinutes {
			if r >= v && r%v == 0 {
				return v
			}
		}
		return 60 // default hour
	}
}

// ConvertChartResolution Converts a Chart to a different Resolution
func ConvertChartResolution(fromResolution, toResolution int64, ch []Candle) ([]Candle, error) {
	if toResolution == fromResolution {
		return ch, nil
	}
	if fromResolution > toResolution || toResolution%fromResolution != 0 {
		return ch, fmt.Errorf("New Res %v and old %v do not fit", toResolution, fromResolution)
	}

	quotient := int(toResolution / fromResolution)
	var newChart []Candle = make([]Candle, 0, len(ch)/quotient)

	for _, c := range ch {
		if c.StartTime.Unix()%toResolution != 0 {
			ch = ch[1:]
		} else {
			break
		}
	}
	for {
		if len(ch) < quotient {
			break
		}
		newChart = append(newChart, ConvertCandleResolution(ch[:quotient]))
		ch = ch[quotient:]
	}
	if len(ch) != 0 {
		newChart = append(newChart, ConvertCandleResolution(ch))
	}
	return newChart, nil
}

// ConvertCandleResolution converts Candles from  lower resolution into a higher resolution
func ConvertCandleResolution(c []Candle) Candle {
	var out Candle = Candle{c[0].Open, c[0].High, c[0].Close, c[0].Low, c[0].Volume, c[0].StartTime}
	if len(c) == 1 {
		return c[0]
	}
	for _, i := range c[1:] {
		out.Close = i.Close
		out.Volume += i.Volume
		if i.High > out.High {
			out.High = i.High
		}
		if i.Low < out.Low {
			out.Low = i.Low
		}
	}
	return out
}

func CheckForHoles(ch []Candle, res int64) {
	c := ch[0]
	resS := res * 60
	fmt.Println("Start looking for holes...")
	for i, v := range ch[1:] {
		if v.StartTime.Unix()-c.StartTime.Unix() != resS {
			fmt.Print(i, " ")
		}
		c = v
	}
	fmt.Println("\nFinished")
}
