package views

import (
	"github.com/mattn/go-shellwords"
	"log"
	"strings"
	"time"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/database"
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
)

func Prompt(terminal *sshd.Terminal) {
	var cancel = make(chan struct{})

	// add defer
	defer func() {
		cancel <- struct{}{}

		err := terminal.Close()
		if err != nil {
			return
		}
	}()

	// get profile from db properly
	profile, err := database.UserFromName(terminal.Conn.Name())
	if err != nil {
		log.Println("Failed to get user profile: ", err)
		return
	}

	// create new session
	session := sessions.New(&sessions.Session{
		Terminal:    terminal,
		UserProfile: profile,
		Created:     time.Now(),
	})

	// start title worker
	go titleWorker(session, cancel)

	// clear screen
	if err := session.Clear(); err != nil {
		return
	}

	if err := session.ExecuteBranding(make(map[string]any), "banner.tfx"); err != nil {
		return
	}

	for {
		// set prompt properly
		terminal.SetPrompt(session.ExecuteBrandingToStringNoError(make(map[string]any), "prompt.tfx"))

		line, err := terminal.ReadLine()
		if err != nil {
			return
		}

		args, err := shellwords.Parse(line)
		if err != nil {
			return
		}

		if strings.HasPrefix(line, "|") ||
			strings.HasPrefix(line, "&") ||
			strings.HasPrefix(line, "<") ||
			strings.HasPrefix(line, ">") ||
			strings.HasPrefix(line, ";") ||
			len(strings.Trim(line, " ")) < 1 {
			continue
		}

		// get command by name
		parent, command, index, err := commands.Parse(session.UserProfile, args...)
		if err != nil {
			if err1 := session.Println(err); err1 != nil {
				return
			}

			continue
		}

		// create a new command context aka. argument parser
		context, err := commands.NewContext(parent, command, args[index:]...)
		if err != nil {
			if err1 := session.Println(err); err1 != nil {
				return
			}

			continue
		}

		// execute command
		if err := command.Executor(session, context); err != nil {
			if err1 := session.Println(err); err1 != nil {
				return
			}

			continue
		}
	}
}
