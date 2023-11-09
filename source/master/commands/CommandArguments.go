package commands

type ArgumentType int

const (
	ArgumentString ArgumentType = iota
	ArgumentInteger
	ArgumentBoolean
	ArgumentUser
)

// Argument is a simple command argument which contains the name & type.
type Argument struct {
	Type ArgumentType
	Name string
}

// ParsedArgument is a parsed argument which contains the value & type.
type ParsedArgument struct {
	Type  ArgumentType
	Value string
}
