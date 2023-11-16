package commands

import (
	"tsundere/source/database"
)

type CommandContext struct {
	arguments map[string]*ParsedArgument
	rawArgs   []string
	parent    *Command
}

func (ctx *CommandContext) User(name string) (*database.UserProfile, error) {
	value, err := ctx.get(name, ArgumentUser)
	if err != nil {
		return nil, err
	}

	return value.(*database.UserProfile), nil
}

func (ctx *CommandContext) String(name string) (string, error) {
	value, err := ctx.get(name, ArgumentString)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (ctx *CommandContext) Integer(name string) (int, error) {
	value, err := ctx.get(name, ArgumentInteger)
	if err != nil {
		return 0, err
	}

	return value.(int), nil
}

func (ctx *CommandContext) Boolean(name string) (bool, error) {
	value, err := ctx.get(name, ArgumentBoolean)
	if err != nil {
		return false, err
	}

	return value.(bool), nil
}

// get gets value(s) from a name
func (ctx *CommandContext) get(name string, typeToGet ArgumentType) (any, error) {
	parsedArgument, exists := ctx.arguments[name]

	if !exists {
		return "", ErrArgumentNotRegistered
	}

	if parsedArgument.Type != typeToGet {
		return "", ErrArgumentInvalidType
	}

	return parsedArgument.Value, nil
}

func (ctx *CommandContext) ParsedCount() int {
	return len(ctx.arguments)
}

func (ctx *CommandContext) Count() int {
	return len(ctx.rawArgs)
}

func (ctx *CommandContext) Parent() *Command {
	return ctx.parent
}
