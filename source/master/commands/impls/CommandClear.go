package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
	"tsundere/source/master/sessions/swashengine"
)

func init() {
	commands.Create(&commands.Command{
		Aliases: []string{"clear", "cls"},
		Roles:   make([]string, 0),
		Executor: func(session *sessions.Session, engine *swashengine.SwashEngine, ctx *commands.CommandContext) error {
			if err := session.Clear(); err != nil {
				return err
			}

			return engine.Execute("banner.tfx", true, engine.Elements(nil))
		},
	})
}