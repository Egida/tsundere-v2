package ansi

import "io"

type Ansi struct {
	Writer io.Writer
}

func NewAnsi(writer io.Writer) *Ansi {
	return &Ansi{Writer: writer}
}
