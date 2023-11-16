package views

import (
	"fmt"
	"time"
	"tsundere/source/database"
	"tsundere/source/master/sessions"
)

func titleWorker(session *sessions.Session, cancel chan struct{}) {
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
					session.ExecuteBrandingToStringNoError(
						make(map[string]any),
						"title.tfx",
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
