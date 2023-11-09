package views

import (
	"fmt"
	"log"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/database"
)

func Prompt(terminal *sshd.Terminal) {
	profile, err := database.UserFromName(terminal.Conn.Name())
	if err != nil {
		log.Println("Failed to get user profile: ", err)
		return
	}

	terminal.SetPrompt(fmt.Sprintf("\x1b[92m%s@xubuntu\u001B[97m:\u001B[94m~\u001B[97m$ \u001B[0m", profile.Username))

	for {
		line, err := terminal.ReadLine()
		if err != nil {
			return
		}

		fmt.Println(line)
	}
}
