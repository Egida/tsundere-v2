package commands

type ArgumentType int

const (
	ArgumentString ArgumentType = iota
	ArgumentInteger
	ArgumentBoolean
	ArgumentUser
)

type Argument struct {
	Type ArgumentType
	Name string
}

type ParsedArgument struct {
	Type  ArgumentType
	Value string
}
