package commands

import (
	"errors"
	"golang.org/x/exp/slices"
	"log"
	"tsundere/source/master/sessions"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
	ErrArgumentNotRegistered    = errors.New("argument not registered")
	ErrArgumentInvalidType      = errors.New("tried to get argument of invalid type")

	Commands = make(map[int]*Command)
)

type Command struct {
	Aliases   []string
	Roles     []string
	Arguments []*Argument

	// Executor is the executor of the command
	Executor func(session *sessions.Session, args []string, ctx *CommandContext)
}

func Create(command *Command) {
	if CommandByName(command.Aliases[0]) != nil {
		log.Println("Failed to register command: ", ErrCommandAlreadyRegistered)
		return
	}

	Commands[len(Commands)+1] = command
}

// CommandByName gets a command by its name
func CommandByName(value string) *Command {
	for _, command := range Commands {
		if slices.Contains(command.Aliases, value) {
			return command
		}

	}

	return nil
}
