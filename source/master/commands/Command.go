package commands

import (
	"errors"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
	"strings"
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

	SubCommandListCommand = &Command{
		Aliases:     []string{"help", "?"},
		Description: "Provides a list of sub-commands.",
		Executor: func(session *sessions.Session, ctx *CommandContext) error {
			if err := session.Println("List of available sub-commands from ", strconv.Quote(ctx.Parent().Aliases[0]), ":"); err != nil {
				return err
			}

			for _, command := range ctx.Parent().SubCommands {
				err := session.Printf(" - %s: %s\r\n", strings.Join(command.Aliases, ", "), command.Description)
				if err == nil {
					continue
				}

				return err
			}
			return nil
		},
	}

	Commands = make([]*Command, 0)
)

// ExecutorFunc is the function for commands
type ExecutorFunc func(session *sessions.Session, ctx *CommandContext) error

// ErrorFunc is the function for errors to execute brandings..etc
type ErrorFunc func(session *sessions.Session) error

// Command is a simple command
type Command struct {
	Aliases     []string
	Roles       []string
	Arguments   []*Argument
	SubCommands []*Command

	Description string

	// Executor is the executor of the command
	Executor ExecutorFunc
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
func Parse(profile *database.UserProfile, args ...string) (parent *Command, command *Command, index int, err error) {
	// find parent command based on the first argument lol
	parent = CommandByName(Commands, args[0])
	if parent == nil {
		return nil, nil, 0, ErrCommandNotRegistered
	}

	// checks if user has enough permissions for parent command
	if !parent.HasNeededRoles(profile) {
		return nil, nil, 0, ErrCommandNotEnoughPermissions
	}

	// check if parent command has subcommands
	if len(parent.SubCommands) == 0 || len(args) == 1 {
		return parent, parent, 1, nil
	}

	var actualParent = parent

	for pos, arg := range args[1:] {
		// find subcommand based on second argument...
		child := CommandByName(parent.SubCommands, arg)
		if child == nil {
			return parent, parent, pos + 1, nil
		}

		// checks if user has enough permissions for sub command
		if !child.HasNeededRoles(profile) {
			return nil, nil, pos, ErrCommandNotEnoughPermissions
		}

		parent = child
	}

	// return child command and index
	return actualParent, parent, 2, nil
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
