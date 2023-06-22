package ta

type DIVERGENCE struct {
	ERS[float64]
}

/*
Divergence 1
HiddenDiv 2
BearishCon 3
BullishCon 4
*/

func LongDiv(src Series, Indicator Series, trigger Condition) Series {
	s := new(DIVERGENCE)
	s.st = trigger.StartTime()
	s.res = src.Resolution()

	cl, ind, c1 := Lowest(src, 2).Data(), OffS(Indicator, 1).Data(), trigger.Data()

	shortest := len(c1)
	cl = cl[len(cl)-shortest:]
	ind = ind[len(ind)-shortest:]
	c1 = c1[len(c1)-shortest:]

	out := make([]float64, 0, len(cl))
	var tcl, tind float64 //tempClose, tempIndicator

	for i, v := range c1 {
		if v {
			out = append(out, divHelp(cl[i], ind[i], tcl, tind))
			tcl = cl[i]
			tind = ind[i]
		} else {
			out = append(out, 0.0)
		}
	}
	s.data = out

	return s

}

func ShortDiv(src Series, Indicator Series, trigger Condition) Series {
	s := new(DIVERGENCE)
	s.st = trigger.StartTime()
	s.res = src.Resolution()

	cl, ind, c1 := Highest(src, 2).Data(), OffS(Indicator, 1).Data(), trigger.Data()

	shortest := len(c1)
	cl = cl[len(cl)-shortest:]
	ind = ind[len(ind)-shortest:]
	c1 = c1[len(c1)-shortest:]

	out := make([]float64, 0, len(cl))
	var tcl, tind float64 //tempClose, tempIndicator
	for i, v := range c1 {
		if v {
			out = append(out, divHelp2(cl[i], ind[i], tcl, tind))
			tcl = cl[i]
			tind = ind[i]
		} else {
			out = append(out, 0.0)
		}
	}

	s.data = out

	return s

}

func divHelp(cl, ind, cl1, ind1 float64) float64 {
	if cl >= cl1 {
		if ind >= ind1 {
			return 4.0
		} else {
			return 2.0
		}
	} else {
		if ind >= ind1 {
			return 1.0
		} else {
			return 3.0
		}
	}
}

func divHelp2(cl, ind, cl1, ind1 float64) float64 {
	if cl >= cl1 {
		if ind >= ind1 {
			return 4.0
		} else {
			return 1.0
		}
	} else {
		if ind >= ind1 {
			return 2.0
		} else {
			return 3.0
		}
	}
}
