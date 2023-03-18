package cle

import (
	"github.com/five-aces-research/toolkit/fas"
	"io"
)

type CLEIO interface {
	fas.Public
	fas.Privat
}

func Execute(Commands []string, f CLEIO, writer io.Writer) error {
	return nil
}
