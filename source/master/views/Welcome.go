package views

import (
	"log"
	"tsundere/packages/customization/termengine"
	"tsundere/packages/customization/termengine/mouse"
	"tsundere/packages/utilities/sshd"
)

func Welcome(terminal *sshd.Terminal) {
	terminal.Channel.Write([]byte("\u001B]0;tsundere.dev | Login\007"))
	terminal.Channel.Write([]byte("\x1bc"))

	var term = termengine.New(terminal.Channel, 81, 20)

	term.DefaultButton(3, 3, func(event *mouse.Event) bool {
		return false
	}, "test")

	if err := term.Run(); err != nil {
		log.Println(err)
		return
	}
}
