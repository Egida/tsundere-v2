package termengine

import (
	"strings"
	"tsundere/packages/customization/termengine/mouse"
)

type Text struct {
	content []string
	x, y    int
}

func (t *TermEngine) Text(x int, y int, content ...string) {
	t.Elements = append(t.Elements, &Text{
		content: content,
		x:       x,
		y:       y,
	})
}

func (t Text) draw() (string, error) {
	return strings.Join(t.content, "\n"), nil
}

func (Text) click(event *mouse.Event) bool {
	return false
}

func (t Text) X() int {
	return t.x
}

func (t Text) Y() int {
	return t.y
}
