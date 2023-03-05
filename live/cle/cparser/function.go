package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

type function struct {
	NumbersOfParameters int
	ParameterPosition   []int
	Token               []clexer.Token
}

func (f *function) Parse(tk []clexer.Token) (new []clexer.Token, rest []clexer.Token, err error) {
	if len(tk) < 3 {
		return new, tk, nerr(empty, "Function Syntax Error, Brackets Missing")
	}

	if tk[0].Type != clexer.LBRACKET {
		return new, tk, nerr(empty, "Function Syntax Error, no bracket")
	}
	var count int
	for i, v := range tk[1:] {
		switch v.Type {
		case clexer.RBRACKET:
			count = i
			break
		default:
			new = append(new, v)
		}
	}

	if len(new) != f.NumbersOfParameters {
		return new, tk, nerr(empty, fmt.Sprintf("Function Error, wrong number of parameters. Want %d got %d", f.NumbersOfParameters, len(rest)))
	}

	for i, t := range new {
		n := f.ParameterPosition[i]
		f.Token[n] = t
	}

	return f.Token, tk[count:], nil
}

func parseAssignFunc(tk []clexer.Token) (f function, err error) {
	if len(tk) == 0 {
		return f, nerr(empty, "Empty Func can't be assigned to a variable")
	}
	if tk[0].Type != clexer.LBRACKET {
		return f, nerr(empty, fmt.Sprintf("%s INVALID Syntax, after a func there must be a '(' "))
	}
	//A map that track which variable is on which position of the tokenlist
	m := make(map[string]int)

	nl := tk[1:]
	var count int

L:
	for _, v := range nl {
		switch v.Type {
		case clexer.RBRACKET:
			break L
		case clexer.VARIABLE:
			m[v.Value] = count //
		default:
			return f, nerr(empty, fmt.Sprintf("invalid variable Name %s: ", v.Value))
		}
		count++
	}

	f.ParameterPosition = make([]int, count)
	tk = tk[count+2:]

	f.NumbersOfParameters = count
	for i, v := range tk {
		if v.Type == clexer.VARIABLE {
			n, ok := m[v.Value]
			if ok {
				f.ParameterPosition[n] = i
				delete(m, v.Value)
			}
		}
		f.Token = append(f.Token, v)
	}

	if len(m) != 0 {
		return f, nerr(empty, fmt.Sprintf("Not all Variables got assigned: %+v", m))
	}

	return f, nil
}
