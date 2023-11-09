package termengine

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"tsundere/packages/customization/termengine/mouse"
)

type TermEngine struct {
	Io            io.ReadWriter
	Width, Height int
	Elements      []Element
}

type Dimension struct {
	Width, Height int
}

func New(term io.ReadWriter, width int, height int) *TermEngine {
	return &TermEngine{
		Io:       term,
		Width:    width,
		Height:   height,
		Elements: make([]Element, 0),
	}
}

func (t *TermEngine) Run() error {
	if _, err := t.Io.Write([]byte(fmt.Sprintf("\x1b[8;%d;%dt\x1bc", t.Height, t.Width))); err != nil {
		return err
	}

	if err := t.Clean(); err != nil {
		return err
	}

	var elements = make(map[Element]*Dimension)
	var buf = make([]string, 0)
	var matrix = make([][]string, t.Height)

	for i := range matrix {
		matrix[i] = make([]string, t.Width)
		for j := range matrix[i] {
			matrix[i][j] = " "
		}
	}

	for _, element := range t.Elements {
		elementX := element.X()
		elementY := element.Y()

		for len(matrix) <= elementY {
			matrix = append(matrix, make([]string, t.Width))
		}

		str, err := element.draw()
		if err != nil {
			return err
		}

		lines := strings.Split(str, "\n")
		for i, line := range lines {
			for len(matrix[elementY+i]) <= elementX+len(line)-1 {
				matrix[elementY+i] = append(matrix[elementY+i], " ")
			}

			lineChars := strings.Split(line, "")
			copy(matrix[elementY+i][elementX:], lineChars)
		}

		var width = 0
		for _, _ = range lines[0] {
			width++
		}

		elements[element] = &Dimension{Width: width, Height: len(lines)}
	}

	for _, row := range matrix {
		buf = append(buf, strings.Join(row, ""))
	}

	if _, err := t.Io.Write([]byte(strings.Join(buf, ""))); err != nil {
		return err
	}

	for {
		event, err := mouse.Listen(t.Io)
		if err != nil {
			continue
		}

		for e, d := range elements {
			if (event.X >= e.X() && event.X <= e.X()+d.Width) && (event.Y >= e.Y()+1 && event.Y <= (e.Y()+d.Height)) {
				if e.click != nil && e.click(event) {
					t.Io.Write([]byte("\x1b[?1000l"))
					return nil
				}

				continue
			}
		}
	}
}

func (t *TermEngine) Clean() error {
	var buffer bytes.Buffer
	buffer.WriteString("\x1b[0m")

	for i := 0; i < t.Height; i++ {
		buffer.WriteString(strings.Repeat(" ", t.Width) + "\r\n")
	}

	_, err := t.Io.Write(buffer.Bytes())
	return err
}
