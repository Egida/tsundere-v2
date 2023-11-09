package goterm2lite

// Button represents a button object on the terminal
type Button struct {
	//a buttons display value can span across multiple lines
	content []string

	X int
	Y int

	inherit *GoTerm2
	onClick func() bool
}

// NewButton creates a brand new button object
func (g *GoTerm2) NewButton(x, y int, content ...string) *Button {
	button := &Button{
		X:       x,
		Y:       y,
		content: content,
		inherit: g,
	}

	g.elements = append(g.elements, button)
	return button
}

// OnClick is how you can register an event within the handler
func (b *Button) OnClick(f func() bool) {
	b.onClick = f
}

// clicked represents an event where the button has been clicked
func (b *Button) clicked() bool {
	return b.onClick()
}

// draw implements the functionality of writing a component
func (b *Button) draw(matrix [][]string) {
	for p := range b.content {
		if b.Y+p-1 > b.inherit.Y {
			break
		}

		feed := Split(b.content[p])
		for data := range feed {
			if b.X+data-1 >= len(matrix[b.Y]) {
				continue
			}

			matrix[b.Y+p-1][b.X:][data] = feed[data]
		}
	}
}

// coordinates will calculate the coordinates of the button
func (b *Button) coordinates() map[int][]int {
	product := make(map[int][]int)
	for pos := range b.content {
		product[b.Y+pos] = make([]int, 0)
		for directional := range Split(b.content[pos]) {
			product[b.Y+pos-1] = append(product[b.Y+pos-1], b.X+directional)
		}
	}

	return product
}
