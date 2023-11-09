package goterm2lite

import (
	"bytes"
	"fmt"
)

// Input is an input which can be clicked on and dismissed.
type Input struct {
	buffer *bytes.Buffer

	mask                 string
	content              []string
	X, Y, RX, RY, maxLen int

	inherit  *GoTerm2
	onSubmit Writeable
}

// NewInput will create a brand new input based upon a button config
func (g *GoTerm2) NewInput(x, rx, y, ry, maxLen int, mask string, content ...string) *Input {
	input := &Input{
		buffer:  bytes.NewBuffer(make([]byte, 0)),
		inherit: g,
		content: content,
		X:       x, RX: rx, Y: y, RY: ry, maxLen: maxLen,
		mask: mask,
	}

	g.elements = append(g.elements, input)
	return input
}

// draw implements the functionality of writing a component
func (i *Input) draw(matrix [][]string) {
	for p := range i.content {
		if i.Y+p-1 > i.inherit.Y {
			break
		}

		feed := Split(i.content[p])
		for data := range feed {
			if i.X+data-1 >= len(matrix[i.Y]) {
				continue
			}

			matrix[i.Y+p-1][i.X:][data] = feed[data]
		}
	}
}

// coordinates will calculate the coordinates of the button
func (i *Input) coordinates() map[int][]int {
	product := make(map[int][]int)
	for pos := range i.content {
		product[i.Y+pos] = make([]int, 0)
		for directional := range Split(i.content[pos]) {
			product[i.Y+pos-1] = append(product[i.Y+pos-1], i.X+directional)
		}
	}

	return product
}

// Value is the current value of the text
func (i *Input) Value() string {
	return i.buffer.String()
}

func (i *Input) SetValue(v string) {
	i.buffer.Reset()
	i.buffer.WriteString(v)
}

// OnSubmit allows you to set custom submit properties
func (i *Input) OnSubmit(w Writeable) {
	i.onSubmit = w
}

// clicked is called whenever the input field is clicked
func (i *Input) clicked() bool {
	if _, err := i.inherit.term.Write([]byte(fmt.Sprintf("\x1b[%d;%dH\033[?25h\033[?0c", i.RY, i.RX+1+i.buffer.Len()))); err != nil {
		return true
	}

	defer func() {
		i.inherit.term.Write([]byte("\033[?25l"))
	}()

	for {
		buf := make([]byte, 128)
		size, err := i.inherit.term.Read(buf)
		if err != nil {
			return true
		}

		buf = buf[:size]
		if bytes.HasPrefix(buf, []byte{27, 91, 77, 32}) && len(buf) == 6 {
			click, ok := i.inherit.handleClick(buf)
			if !ok || click == nil {
				return false
			}

			return click.clicked()
		}

		/* handles the incoming buffer */
		switch buf[0] {

		case 9:
			if i.inherit.tab() {
				return true
			}

			continue

		case 127:
			if i.buffer.Len() <= 0 {
				continue
			}

			i.buffer.Truncate(i.buffer.Len() - 1)
			if _, err := i.inherit.term.Write([]byte{27, 91, 49, 68, 32, 27, 91, 49, 68}); err != nil {
				return true
			}

		case 13:
			if _, err := i.inherit.term.Write([]byte{13, 10}); err != nil {
				return true
			}

			if i.onSubmit == nil {
				return false
			}

			ok := i.onSubmit.clicked()
			if !ok {
				return false
			}

			return ok

		case 32, 96, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 45, 61, 194, 172, 33, 34, 163, 36, 37, 94, 38, 42, 40, 41, 95, 43, 113, 119, 101, 114, 116, 121, 117, 105, 111, 112, 91, 93, 81, 87, 69, 82, 84, 89, 85, 73, 79, 80, 123, 125, 97, 115, 100, 102, 103, 104, 106, 107, 108, 59, 39, 35, 65, 83, 68, 70, 71, 72, 74, 75, 76, 58, 64, 126, 92, 122, 120, 99, 118, 98, 110, 109, 44, 46, 47, 124, 90, 88, 67, 86, 66, 78, 77, 60, 62, 63:
			if i.buffer.Len() > i.maxLen {
				continue
			}

			i.buffer.WriteByte(buf[0])

			if len(i.mask) >= 1 {
				buf = []byte(i.mask)
			}

			if _, err := i.inherit.term.Write([]byte{buf[0]}); err != nil {
				return true
			}
		}
	}
}
