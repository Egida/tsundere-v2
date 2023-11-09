package views

import (
	"log"
	"strings"
	"tsundere/packages/customization/goterm2lite"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/master/views/styles"
)

func Welcome(terminal *sshd.Terminal) {
	terminal.Channel.Write([]byte("\u001B]0;tsundere.dev | Login\007"))
	terminal.Channel.Write([]byte("\x1bc"))

	// initialize goterm2lite instance for login page
	var term = goterm2lite.New(terminal.Channel, 81, 20)

	term.NewButton(3, 3, strings.Split(styles.SmallButton.Render("Test"), "\n")...).OnClick(func() bool {
		log.Println("button clicked :(")
		return false
	})

	// run goterm2lite
	if err := term.Run(); err != nil {
		log.Println(err)
		return
	}
}
