package commands

import (
	"errors"
	"golang.org/x/exp/slices"
	"log"
	"tsundere/source/database"
	"tsundere/source/master/sessions"
)

var (
	ErrCommandAlreadyRegistered    = errors.New("command already registered")
	ErrCommandNotRegistered        = errors.New("command not registered")
	ErrCommandNotEnoughPermissions = errors.New("not enough permissions to run command")

	ErrArgumentNotRegistered = errors.New("argument not registered")
	ErrArgumentInvalidType   = errors.New("tried to get argument of invalid type")
	ErrArgumentRequired      = errors.New("missing required argument")

	Commands = make([]*Command, 0)
)

// Command is a simple command
type Command struct {
	Aliases   []string
	Roles     []string
	Arguments []*Argument
	Commands  []*Command

	// Executor is the executor of the command
	Executor func(session *sessions.Session, ctx *CommandContext) error
}

// Create adds a command to the registry
func Create(command *Command) *Command {
	if CommandByName(Commands, command.Aliases[0]) != nil {
		log.Println("Failed to register command: ", ErrCommandAlreadyRegistered)
		return nil
	}

	Commands = append(Commands, command)
	return command
}

// CommandByName gets a command by its name
func CommandByName(list []*Command, value string) *Command {
	for _, command := range list {
		if slices.Contains(command.Aliases, value) {
			return command
		}

	}

	return nil
}

// Parse will try to get a command within an argument array. (Parses sub-commands too)
func Parse(profile *database.UserProfile, args ...string) (command *Command, index int, err error) {
	// find parent command based on the first argument lol
	parent := CommandByName(Commands, args[0])
	if parent == nil {
		return nil, 0, ErrCommandNotRegistered
	}

	// checks if user has enough permissions for parent command
	if !parent.HasNeededRoles(profile) {
		return nil, 0, ErrCommandNotEnoughPermissions
	}

	// check if parent command has subcommands
	if len(parent.Commands) == 0 || len(args) == 1 {
		return parent, 1, nil
	}

	// find subcommand based on second argument...
	child := CommandByName(parent.Commands, args[1])
	if child == nil {
		return parent, 1, nil
	}

	// checks if user has enough permissions for sub command
	if !child.HasNeededRoles(profile) {
		return nil, 0, ErrCommandNotEnoughPermissions
	}

	// return child command and index
	return child, 2, nil
}

// HasNeededRoles checks if the user has a needed role for the command
func (c *Command) HasNeededRoles(profile *database.UserProfile) bool {
	if len(c.Roles) < 1 {
		return true
	}

	// iterate through command roles
	for _, literal := range c.Roles {
		if !profile.HasRole(literal) {
			continue
		}

		return true
	}

	return false
}
