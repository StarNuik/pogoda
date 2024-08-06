package pgwire

import (
	"io"
)

type Wire struct {
	w io.Writer
	r io.Reader
}

func NewWire(w io.Writer, r io.Reader) *Wire {
	return &Wire{
		w: w,
		r: r,
	}
}
