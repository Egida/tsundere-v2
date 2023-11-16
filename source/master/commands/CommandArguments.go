package commands

import (
	"strconv"
	"tsundere/source/database"
)

type ArgumentType int

const (
	ArgumentString ArgumentType = iota
	ArgumentInteger
	ArgumentBoolean
	ArgumentUser
)

// Argument is a simple command argument which contains the name & type.
type Argument struct {
	Type     ArgumentType
	Name     string
	Default  any
	Required bool
}

// ParsedArgument is a parsed argument which contains the value & type.
type ParsedArgument struct {
	Type  ArgumentType
	Value any
}

// NewArgument creates a new instance of an argument
func NewArgument(Name string, Default any, Type ArgumentType, Required bool) *Argument {
	return &Argument{
		Type:     Type,
		Name:     Name,
		Default:  Default,
		Required: Required,
	}
}

// Literal2Go converts a literal to a corresponding go type
func (a *Argument) Literal2Go(literal string) (any, error) {
	switch a.Type {
	case ArgumentString:
		return literal, nil
	case ArgumentUser:
		return database.UserFromName(literal)
	case ArgumentInteger:
		return strconv.Atoi(literal)
	case ArgumentBoolean:
		return strconv.ParseBool(literal)
	}

	return nil, nil
}
