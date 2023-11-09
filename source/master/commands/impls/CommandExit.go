package impls

import (
	"tsundere/source/master/commands"
	"tsundere/source/master/sessions"
	"tsundere/source/master/sessions/swashengine"
)

func init() {
	commands.Create(&commands.Command{
		Aliases: []string{"exit", "quit", "logout"},
		Roles:   make([]string, 0),
		Executor: func(session *sessions.Session, engine *swashengine.SwashEngine, ctx *commands.CommandContext) error {
			session.Close()
			return nil
		},
	})
}
