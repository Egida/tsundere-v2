package goterm2lite

// Text is an exported piece of information being wrote to the terminal
type Text struct {
	content string
	X       int
	Y       int

	inherit *GoTerm2
}

// NewText will create the text on the terminal
func (g *GoTerm2) NewText(x, y int, content string) {
	g.elements = append(g.elements, &Text{
		content: content,
		X:       x, Y: y,
		inherit: g,
	})
}

// draw implements the functionality of writing
func (t *Text) draw(matrix [][]string) {
	if t.Y > t.inherit.Y {
		return
	}

	feed := Split(t.content)
	for data := range feed {
		if t.X+data >= len(matrix[t.Y]) {
			continue
		}

		matrix[t.Y][t.X:][data] = feed[data]
	}
}

// coordinates returns the position of the data on the screen
func (t *Text) coordinates() map[int][]int {
	coordinates := make(map[int][]int)
	coordinates[t.Y] = make([]int, 0)
	for pos := range Split(t.content) {
		coordinates[t.Y] = append(coordinates[t.Y], t.X+pos)
	}

	return coordinates
}

// clicked represents a click event on the text.
func (t *Text) clicked() bool {
	return false
}
