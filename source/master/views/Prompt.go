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
	"tsundere/source/master/sessions/swashengine"
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

	swashEngine := swashengine.New(session)

	// start title worker
	go titleWorker(session, swashEngine, cancel)

	// clear screen
	if err := session.Clear(); err != nil {
		return
	}

	if err := swashEngine.Execute("banner.tfx", true, swashEngine.Elements(nil)); err != nil {
		return
	}

	for {
		// set prompt properly
		terminal.SetPrompt(swashEngine.ExecuteString("prompt.tfx", swashEngine.Elements(nil)))

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
		command := commands.CommandByName(args[0])
		if command == nil {
			if err := session.Printf("%s: command not found\r\n", args[0]); err != nil {
				return
			}

			continue
		}

		// create a new command context aka. argument parser
		context, err := commands.NewContext(command, args[1:])
		if err != nil {
			if xErr := session.Println(err); xErr != nil {
				return
			}

			continue
		}

		// execute command
		if err := command.Executor(session, swashEngine, context); err != nil {
			if xErr := session.Println(err); xErr != nil {
				return
			}

			continue
		}
	}
}
