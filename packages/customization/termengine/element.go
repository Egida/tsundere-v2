package termengine

import "tsundere/packages/customization/termengine/mouse"

type Element interface {
	draw() (string, error)
	click(event *mouse.Event) bool
	X() int
	Y() int
}
