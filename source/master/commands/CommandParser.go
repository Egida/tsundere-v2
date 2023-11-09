package commands

func ParseArguments(cmd *Command, arguments []string) (*CommandContext, error) {
	var ctx = new(CommandContext)
	ctx.arguments = make(map[string]*ParsedArgument)

	for i, argument := range cmd.Arguments {
		var arg = arguments[i]

		ctx.arguments[argument.Name] = &ParsedArgument{
			Type:  argument.Type,
			Value: arg,
		}
	}

	return ctx, nil
}
