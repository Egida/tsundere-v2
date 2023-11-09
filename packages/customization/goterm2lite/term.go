package goterm2lite

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// Writeable is an object type which can be written within the Run function
type Writeable interface {
	draw([][]string)
	coordinates() map[int][]int
	clicked() bool
}

// GoTerm2 holds all the data about a terminal
type GoTerm2 struct {
	tabPos   int
	X, Y     int
	term     io.ReadWriter
	elements []Writeable
}

// New creates a new GoTerm2 terminal
func New(term io.ReadWriter, x, y int) *GoTerm2 {
	if _, err := term.Write([]byte(fmt.Sprintf("\033[8;%d;%dt\033c", y, x))); err != nil {
		return nil
	}

	return &GoTerm2{
		X:        x,
		Y:        y,
		tabPos:   -1,
		term:     term,
		elements: make([]Writeable, 0),
	}
}

func (g *GoTerm2) Shake(rotations int, timer time.Duration, channel ssh.Channel) error {
	current, err := GetWindowPosition(channel)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(current)

	for times := 0; times < rotations; times++ {
		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal+3, current.Vertical+3)))
		time.Sleep(timer)

		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal, current.Vertical+3)))
		time.Sleep(timer)

		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal+3, current.Vertical)))
		time.Sleep(timer)

		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal-3, current.Vertical-3)))
		time.Sleep(timer)

		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal-3, current.Vertical)))
		time.Sleep(timer)

		channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal, current.Vertical-3)))
		time.Sleep(timer)
	}

	channel.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", current.Horizontal, current.Vertical)))
	return nil
}

// Run will draw and then continue to read from the terminal
func (g *GoTerm2) Run() error {
	matrix := make([][]string, g.Y+1)
	for line := range matrix {
		matrix[line] = strings.Split(strings.Repeat(" ", g.X), "")
	}

	/* builds the elements on the matrix layer */
	for _, element := range g.elements {
		element.draw(matrix)
	}

	/* builds the matrix later now */
	buf := make([]string, 0)
	for _, x := range matrix {
		buf = append(buf, strings.Join(x, ""))
	}

	defer g.term.Write([]byte("\x1b[?1000l"))
	if _, err := g.term.Write([]byte(strings.Join(buf, "\r\n") + "\033[?1000h\033[?25l")); err != nil {
		return err
	}

	for {
		buf := make([]byte, 1024)
		size, err := g.term.Read(buf)
		if err != nil {
			return err
		}

		buf = buf[:size]
		w, ok := g.handleClick(buf)
		if !ok || w == nil {
			if buf[0] == 9 && g.tab() {
				return nil
			}

			continue
		}

		if w.clicked() {
			return nil
		}
	}
}

func (g *GoTerm2) Display() error {
	if _, err := g.term.Write([]byte("\033[2J")); err != nil {
		return err
	}

	if _, err := g.term.Write([]byte("\033c")); err != nil {
		return err
	}

	matrix := make([][]string, g.Y+1)
	for line := range matrix {
		matrix[line] = strings.Split(strings.Repeat(" ", g.X), "")
	}

	/* builds the elements on the matrix layer */
	for _, element := range g.elements {
		element.draw(matrix)
	}

	/* builds the matrix later now */
	buf := make([]string, 0)
	for _, x := range matrix {
		buf = append(buf, strings.Join(x, ""))
	}

	defer g.term.Write([]byte("\x1b[?1000l"))
	if _, err := g.term.Write([]byte(strings.Join(buf, "\r\n") + "\033[?1000h\033[?25l")); err != nil {
		return err
	}

	return nil
}

// handleClick will try find the click item and return it once we've found it
func (g *GoTerm2) handleClick(buf []byte) (Writeable, bool) {
	if len(buf) != 6 || !bytes.HasPrefix(buf, []byte{27, 91, 77, 32}) {
		return nil, false
	}

	y, x := int(buf[4:][1])-33, int(buf[4:][0])-33
	for _, c := range g.elements {
		yx, ok := c.coordinates()[y]
		if !ok || !slices.Contains(yx, x) {
			continue
		}

		//ignores the text clicks
		if _, ok := c.(*Text); ok {
			continue
		}

		return c, true
	}

	return nil, false
}

// tab handles the key press for the tab key
func (g *GoTerm2) tab() bool {
	for g.tabPos <= len(g.elements) {
		g.tabPos++
		if g.tabPos >= len(g.elements) {
			g.tabPos = 0
		}

		if _, ok := g.elements[g.tabPos].(*Input); ok {
			break
		}

	}

	/*
		tab handles the tab key press, this means it scrolls through all the interfaces like intput, and
		can select them, selections mimic clicks.
	*/

	return g.elements[g.tabPos].clicked()
}
