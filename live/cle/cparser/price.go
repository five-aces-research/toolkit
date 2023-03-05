package cparser

import (
	"errors"
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"strconv"
)

type PriceType int

const (
	PRICE PriceType = iota
	DIFFERENCE
	PERCENTPRICE
	MARKET
)

type Price struct {
	Type        PriceType
	PriceSource string
	Duration    int64      //Optional
	IsLaddered  [2]bool    //Optional 0,0 -> no, 1,0 -> laddered; 1,1 -> exponential laddered
	Values      [3]float64 // [0] Seperation [1]Value1 [2]Value2
}

func ParsePrice(tk []clexer.Token) (p Price, err error) {
	if len(tk) == 0 {
		return p, nerr(empty, "Error Parse Price no Input")
	}

	p.PriceSource = "close"
	if tk[0].Type == clexer.SOURCE {
		switch tk[0].Value {
		case "high", "low", "close", "open":
			p.PriceSource = tk[0].Value
		default:
			return p, nerr(empty, fmt.Sprintf("Parse Price, Invalid Source Value %s", tk[0].Value))
		}

		if len(tk) < 2 || tk[1].Type != clexer.DURATION {
			return p, nerr(empty, fmt.Sprintf("after -%s a duration has to follow", p.PriceSource))
		}

		ss := tk[1].Value
		n, err := strconv.Atoi(ss[:len(ss)-1])
		if err != nil {
			return p, err
		}

		switch ss[len(ss)-1] {
		case 'h':
			n *= 3600
		case 'm':
			n *= 60
		case 'd':
			n *= 3600 * 24
		default:
			return p, nerr(empty, fmt.Sprintf("Error Price Parsing Duration with %s !!", ss))
		}
		p.Duration = int64(n)

		tk = tk[2:]
	}

	switch tk[0].Type {
	case clexer.FLOAT: // 30000 places order at $30000
		p.Type = PRICE
	case clexer.DFLOAT: // -300 places order $300 below the marketprice
		p.Type = DIFFERENCE
	case clexer.PERCENT: // 2% places order 2% below the marketprice
		p.Type = PERCENTPRICE
	case clexer.MARKET: // -market market buys
		p.Type = MARKET
	case clexer.FLAG: // -l -le for laddered Orders
		err = ParsePriceFlag(&p, tk[0].Value, tk[1:])
		return p, err
	default:
		return p, nerr(empty, fmt.Sprintf("Error Price Parsing, %v %s is not a valid price", tk[0].Type, tk[0].Value))
	}

	p.Values[0], err = strconv.ParseFloat(tk[0].Value, 64)
	if err != nil {
		return p, nerr(err, fmt.Sprintf("Error Price Parsing %s is not a Number", tk[0].Value))
	}

	return p, nil
}

// ParsePriceFlag parses laddered Order
func ParsePriceFlag(p *Price, flag string, tl []clexer.Token) (err error) {
	if len(tl) > 3 {
		return nerr(empty, "Parse Price Flag Error, Not Enough Arguments")
	}

	switch flag {
	case "l": //laddered Order
		p.IsLaddered = [2]bool{true, false}
	case "le": //exponential laddered Order
		p.IsLaddered = [2]bool{true, true}
	default:
		return errors.New("This Flag is not supported: " + flag)
	}

	if len(tl) < 3 {
		return errors.New("Not enough Arguments for a laddered order")
	}

	if tl[0].Type == clexer.FLOAT { //First Value sets up how many orders are placed
		num, err := strconv.Atoi(tl[0].Value)
		if err != nil {
			return err
		}

		if num > 25 || num < 2 {
			return nerr(empty, "Error Parse Price Flag, number of seperation to high, max is 25")
		}
		p.Values[0] = float64(num)
	} else {
		return nerr(empty, fmt.Sprintf("Error Parse Price Flage, First Value %s must be a Number ", tl[0].Value))
	}

	if tl[1].Type != tl[2].Type {
		return nerr(empty, "Values 2 and 3 Arguments must be from same type")
	}

	switch tl[1].Type {
	case clexer.FLOAT:
		p.Type = PRICE
	case clexer.DFLOAT:
		p.Type = DIFFERENCE
	case clexer.PERCENT:
		p.Type = PERCENTPRICE
	default:
		return nerr(empty, fmt.Sprintf("Error Parsing Price Flag! %+v is not a legit Pricevalue", tl[1]))
	}
	v1, err := strconv.ParseFloat(tl[1].Value, 64)
	if err != nil {
		return nerr(err, "Error Parsing Price Flag!")
	}

	if err != nil {
		return err
	}
	v2, err := strconv.ParseFloat(tl[2].Value, 64)

	p.Values[1], p.Values[2] = v1, v2

	return nil
}
