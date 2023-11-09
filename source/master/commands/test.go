package commands

import "tsundere/source/master/sessions"

func init() {
	Create(&Command{
		Aliases: []string{"test"},
		Roles:   make([]string, 0),
		Arguments: []*Argument{
			{
				Type: ArgumentString,
				Name: "prefix",
			},
		},
		Executor: func(session *sessions.Session, args []string, ctx *CommandContext) {
			session.Terminal.Channel.Write([]byte("test\r\n"))

			s, err := ctx.String("prefix")
			if err != nil {
				return
			}

			session.Terminal.Channel.Write([]byte(s + "\r\n"))

		},
	})
}
