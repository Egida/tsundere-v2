package sessions

import (
	"time"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/database"
)

type Session struct {
	ID int

	*sshd.Terminal
	*database.UserProfile

	Created time.Time
}
