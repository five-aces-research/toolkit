package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

type VariableType int

const (
	FUNCTION VariableType = iota
	CONSTANT
	DELETE
)

type Variable struct {
	Type    VariableType
	Content interface{}
	//Raw     string
}

/*
	Parser geht durch tk und ersetzt jede Variable durch seinen gegenwert.

*/

func parseVariables(tk []clexer.Token, c Communicator) ([]clexer.Token, error) {
	var ntk []clexer.Token
	for _, v := range tk {
		if v.Type == clexer.VARIABLE {

		}
	}
	ntk = tk
	return ntk, nil
}

func parseVariable(v Variable, tk []clexer.Token) (new []clexer.Token, rest []clexer.Token, err error) {
	switch v.Type {
	case FUNCTION:
		nk, err := parseFunc(v.Content, tk)
		if err != nil {
			return tk, tk, err
		}
		return nk, tk, nil
	case CONSTANT:
		nk, ok := v.Content.([]clexer.Token)
		if !ok {
			return tk, tk, nerr(empty, "Error Parse Variable, Variable not existing")
		}
		return append(nk, tk...), tk, nil
	}

	return tk, nil, nerr(empty, "Error while Parsing a Variable Something went wrong")
}

func parseFunc(v interface{}, tk []clexer.Token) ([]clexer.Token, error) {
	fun, ok := v.(function)
	if !ok {
		return tk, nerr(empty, fmt.Sprintf("Unexpected ERROR with %v ", fun))
	}
	e, _, err := fun.Parse(tk)
	return e, err
}
