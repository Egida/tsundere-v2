package commands

type ArgumentType int

const (
	ArgumentString ArgumentType = iota
	ArgumentInteger
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
