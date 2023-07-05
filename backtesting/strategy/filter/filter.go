package filter

import "github.com/five-aces-research/toolkit/backtesting/strategy"

func Greater(index1 int, index2 int) strategy.Filter {
	return func(sf []strategy.SafeFloat) bool {
		if sf[index1].Safe && sf[index2].Safe {
			return sf[index1].Value > sf[index2].Value
		}
		return false
	}
}

func GreaterAs(index int, val float64) strategy.Filter {
	return func(sf []strategy.SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value > val
		}
		return false
	}
}

func SmallerAs(index int, val float64) strategy.Filter {
	return func(sf []strategy.SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value < val
		}
		return false
	}
}

func Between(index int, lowerBound, upperBound float64) strategy.Filter {
	return func(sf []strategy.SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value > lowerBound && sf[index].Value < upperBound
		}
		return false
	}
}
