package commands

// NewContext creates a new command context
func NewContext(cmd *Command, arguments []string) (*CommandContext, error) {
	// create new command context
	var ctx = new(CommandContext)
	ctx.arguments = make(map[string]*ParsedArgument)
	ctx.rawArgs = arguments

	if len(cmd.Arguments) < 1 {
		return ctx, nil
	}

	// force arguments, shit code but it works
	if !cmd.ForceArguments && len(arguments) < len(cmd.Arguments) {
		for i, arg := range arguments {
			if i > len(cmd.Arguments) {
				continue
			}

			argument := cmd.Arguments[i]

			ctx.arguments[argument.Name] = &ParsedArgument{
				Type:  argument.Type,
				Value: arg,
			}
		}

		return ctx, nil
	}

	if cmd.ForceArguments && len(arguments) < len(cmd.Arguments) {
		return nil, ErrNotEnoughArguments
	}

	for i, argument := range cmd.Arguments {
		var arg string
		if i < len(arguments) {
			arg = arguments[i]
		}

		ctx.arguments[argument.Name] = &ParsedArgument{
			Type:  argument.Type,
			Value: arg,
		}
	}

	return ctx, nil
}
