package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

func parseAssign(token clexer.Token, tk []clexer.Token, c Communicator) error {
	if len(tk) == 0 {
		return nerr(empty, "Nothing can't be assigned to a Variable")
	}
	if token.Type != clexer.VARIABLE {
		return nerr(empty, fmt.Sprintf("Can't use %s %s as variable to assign too", token.Stringer(), token.Value))
	}

	switch tk[0].Type {
	case clexer.FUNC:
		r, err := parseAssignFunc(tk)
		if err != nil {
			return err
		}
		c.AddVariable(token.Value, Variable{Type: CONSTANT, Content: r})
	default:
		c.AddVariable(token.Value, Variable{Type: CONSTANT, Content: tk})
	}
	c.Write([]byte(fmt.Sprintf("Variable %s assigned succesfully", token.Value)))
	return nil
}
