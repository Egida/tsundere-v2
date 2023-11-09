package views

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"log"
	"time"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/database"
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
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

		args, err := shellwords.Parse(line)
		if err != nil {
			return
		}

		cmd := commands.CommandByName(args[0])
		if cmd == nil {
			continue
		}

		args = args[1:]

		arguments, err := commands.ParseArguments(cmd, args)
		if err != nil {
			return
		}

		cmd.Executor(&sessions.Session{
			ID:          0,
			Terminal:    terminal,
			UserProfile: nil,
			Created:     time.Time{},
		}, args, arguments)
	}
}
