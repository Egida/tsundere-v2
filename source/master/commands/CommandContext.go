package commands

import (
	"strconv"
	"tsundere/source/database"
)

type CommandContext struct {
	arguments map[string]*ParsedArgument
	rawArgs   []string
}

func (ctx *CommandContext) User(name string) (*database.UserProfile, error) {
	value, err := ctx.get(name, ArgumentUser)
	if err != nil {
		return nil, err
	}

	user, err := database.UserFromName(value)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ctx *CommandContext) String(name string) (string, error) {
	value, err := ctx.get(name, ArgumentString)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (ctx *CommandContext) Integer(name string) (int, error) {
	value, err := ctx.get(name, ArgumentInteger)
	if err != nil {
		return -1, err
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func (ctx *CommandContext) Boolean(name string) (bool, error) {
	value, err := ctx.get(name, ArgumentBoolean)
	if err != nil {
		return false, err
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return v, nil
}

// get gets value(s) from a name
func (ctx *CommandContext) get(name string, typeToGet ArgumentType) (string, error) {
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
