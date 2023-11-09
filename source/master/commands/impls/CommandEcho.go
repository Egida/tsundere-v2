package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
	"tsundere/source/master/sessions/swashengine"
)

func init() {
	commands.Create(&commands.Command{
		Aliases:        []string{"echo"},
		Roles:          make([]string, 0),
		ForceArguments: true,
		Arguments: []*commands.Argument{
			{
				Type: commands.ArgumentString,
				Name: "message",
			},
		},
		Executor: func(session *sessions.Session, engine *swashengine.SwashEngine, ctx *commands.CommandContext) error {
			s, err := ctx.String("message")
			if err != nil {
				return err
			}

			return session.Println(s)
		},
	})
}
