package termengine

import (
	"strings"
	"tsundere/packages/customization/termengine/mouse"
)

type Button struct {
	content []string
	x, y    int
	onClick func(event *mouse.Event) bool
}

func (t *TermEngine) Button(x int, y int, click func(event *mouse.Event) bool, content ...string) {
	t.Elements = append(t.Elements, &Button{
		content: content,
		x:       x,
		y:       y,
		onClick: click,
	})
}

func (t *TermEngine) DefaultButton(x int, y int, click func(event *mouse.Event) bool, content string) {
	t.Elements = append(t.Elements, &Button{
		content: strings.Split(DefaultButtonStyle.Copy().Width(len(content)).Render(content), "\n"),
		x:       x,
		y:       y,
		onClick: click,
	})
}

func (b Button) draw() (string, error) {
	return strings.Join(b.content, "\n"), nil
}

func (b Button) click(event *mouse.Event) bool {
	return b.onClick(event)
}

func (b Button) X() int {
	return b.x
}

func (b Button) Y() int {
	return b.y
}
