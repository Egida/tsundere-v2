package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
)

func init() {
	commands.Create(&commands.Command{
		Aliases: []string{"echo"},
		Roles:   make([]string, 0),
		Arguments: []*commands.Argument{
			{
				Type: commands.ArgumentString,
				Name: "message",
			},
		},
		Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
			s, err := ctx.String("message")
			if err != nil {
				return err
			}

			return session.Println(s)
		},
	})
}
