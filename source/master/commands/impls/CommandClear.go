package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
)

func init() {
	commands.Create(&commands.Command{
		Aliases: []string{"clear", "cls"},
		Roles:   make([]string, 0),
		Executor: func(session *sessions.Session, ctx *commands.CommandContext) error {
			if err := session.Clear(); err != nil {
				return err
			}

			return session.ExecuteBranding(make(map[string]any), "banner.tfx")
		},
	})
}
