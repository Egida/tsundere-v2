package views

import (
	"log"
	"tsundere/packages/customization/termengine"
	"tsundere/packages/customization/termengine/mouse"
	"tsundere/packages/utilities/sshd"
)

func Welcome(terminal *sshd.Terminal) {
	var engine = termengine.New(terminal.Channel, 80, 20)

	engine.DefaultButton(15, 2, func(event *mouse.Event) bool {
		if event.Click == mouse.ScrollClick || event.Click == mouse.ScrollDown || event.Click == mouse.ScrollUp {
			return false
		}

		log.Println("This button has been clicked.", event.Click.String())

		return false
	}, "Button")

	err := engine.Run()
	if err != nil {
		return
	}
}

func Place(terminal *sshd.Terminal, width, height int, text string) {
	defer terminal.Channel.Write([]byte("\x1b[?1000l"))
	if _, err := terminal.Channel.Write([]byte(text + "\033[?1000h\033[?25l")); err != nil {
		return
	}
}
