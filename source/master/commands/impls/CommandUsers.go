package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
)

func init() {
	var AddCommand = &commands.Command{
		Aliases:     []string{"add", "create"},
		Description: "Insert a user into the database.",
		Arguments: []*commands.Argument{
			commands.NewArgument("username", nil, commands.ArgumentString, true),
			commands.NewArgument("password", "changeme", commands.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
			username, err := ctx.String("username")
			if err != nil {
				return err
			}

			password, err := ctx.String("password")
			if err != nil {
				return err
			}

			return session.Println(username, " ", password)
		},
	}

	commands.Create(&commands.Command{
		Aliases: []string{"users"},
		Roles:   []string{"admin"},
		SubCommands: []*commands.Command{
			commands.SubCommandListCommand,
			AddCommand,
		},
		Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
			return session.Println("fookin table here that shows users")
		},
	})
}
