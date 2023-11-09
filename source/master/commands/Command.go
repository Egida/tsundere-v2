package commands

import (
	"errors"
	"golang.org/x/exp/slices"
	"log"
	"tsundere/source/master/sessions"
	"tsundere/source/master/sessions/swashengine"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
	ErrArgumentNotRegistered    = errors.New("argument not registered")
	ErrArgumentInvalidType      = errors.New("tried to get argument of invalid type")
	ErrNotEnoughArguments       = errors.New("not enough arguments")

	Commands = make(map[int]*Command)
)

// Command is a simple command
type Command struct {
	Aliases        []string
	Roles          []string
	Arguments      []*Argument
	ForceArguments bool

	// Executor is the executor of the command
	Executor func(session *sessions.Session, engine *swashengine.SwashEngine, ctx *CommandContext) error
}

// Create adds a command to the registry
func Create(command *Command) *Command {
	if CommandByName(command.Aliases[0]) != nil {
		log.Println("Failed to register command: ", ErrCommandAlreadyRegistered)
		return nil
	}

	Commands[len(Commands)+1] = command
	return command
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
