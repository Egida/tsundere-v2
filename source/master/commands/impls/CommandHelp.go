package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
)

func init() {
	commands.Create(&commands.Command{
		Aliases: []string{"help", "?"},
		Commands: []*commands.Command{
			{
				Aliases: []string{"admin"},
				Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
					return session.Println("this is a fucking command named admin in a fucking help command")
				},
			},
		},
		Arguments: []*commands.Argument{
			commands.NewArgument("query", nil, commands.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
			if v, err := ctx.String("query"); err == nil {
				return session.Println("Retrieved search query: ", v)
			}

			return session.ExecuteBranding(make(map[string]any), "commands", "help.tfx")
		},
	})
}
