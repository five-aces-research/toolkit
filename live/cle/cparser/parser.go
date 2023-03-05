package cparser

import (
	"errors"
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"io"
)

// Command Line Execution
// This package is used for executing trades with a user terminal

type Communicator interface {
	io.Writer
	io.Reader
	AddVariable(string, Variable)         //Add an Variable
	GetVariable(string) (Variable, error) //Return a Variable
	ErrorMessage(s error)
}

type Parser interface {
	Evaluate(w Communicator, f cle.CLEIO)
}

// Parse returns a Parser which then gets Evaluated and returns
func Parse(tk []clexer.Token, c Communicator) (Parser, error) {
	if len(tk) == 0 {
		return nil, nerr(empty, "Error nothing got lexed")
	}

	//Check if is assign or function
	if len(tk) > 2 {
		if tk[1].Type == clexer.ASSIGN {
			return nil, parseAssign(tk[0], tk[2:], c)
		}
	}
	var err error
	if tk, err = parseVariables(tk, c); err != nil {
		return nil, err
	}

	var o Parser

	switch tk[0].Type {
	case clexer.SIDE: // buy, sell
		o, err = ParseOrder(tk[0].Value, tk[1:])
	case clexer.STOP: //stop
		//o, err = ParseStop(nk[1:])
	case clexer.CANCEL: //cancel
		//o, err = ParseCancel(nk[1:])
	case clexer.CLOSE: //fclose
		//o, err = ParseClose(nk[1:])
	case clexer.VARIABLE:
	case clexer.FUNDINGPAYS:
	case clexer.FUNDINGRATES:
	default:
		return o, nerr(empty, fmt.Sprintf("Invalid Type Error during Parsing %v", nk[0].Type))
	}

	if err != nil {
		return o, err
	}

	return o, nil
}

type parseError struct {
	err error
	msg string
}

var empty = errors.New("")

func nerr(err error, msg string) *parseError {
	return &parseError{err, msg}
}

func (e *parseError) Error() string {
	return fmt.Sprintf("Message:%s : %v", e.msg, e.err)
}
