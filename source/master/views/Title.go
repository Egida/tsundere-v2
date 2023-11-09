package views

import (
	"fmt"
	"time"
	"tsundere/source/database"
	"tsundere/source/master/sessions"
	"tsundere/source/master/sessions/swashengine"
)

func titleWorker(session *sessions.Session, swashEngine *swashengine.SwashEngine, cancel chan struct{}) {
	// new ticker
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			var err error

			// update user profile from DB
			session.UserProfile, err = database.UserFromName(session.Username)
			if err != nil || time.Now().After(session.UserProfile.Expiry) {
				// tell user that their access has been terminated
				if err := session.Println("Your access has been terminated."); err != nil {
					return
				}

				// cancel title worker
				cancel <- struct{}{}
				return
			}

			if n, err := session.Channel.Write([]byte(
				fmt.Sprintf("\033]0;%s\007",
					swashEngine.ExecuteString(
						"title.tfx",
						swashEngine.Elements(nil),
					),
				),
			)); err != nil || n <= 0 {
				cancel <- struct{}{}
				return
			}
		case <-cancel:
			session.Remove()
			return
		}
	}
}
