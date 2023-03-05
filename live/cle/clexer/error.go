package clexer

import "fmt"

type lexerError struct {
	input      string
	err        error
	errmessage string
}

func nerr(input string, err error, errmessage string) *lexerError {
	return &lexerError{input, err, errmessage}
}

func (l *lexerError) Error() string {
	return fmt.Sprintf("Text: %s, Error: %v + %s", l.input, l.err, l.errmessage)
}
