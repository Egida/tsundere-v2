package sessions

import (
	"fmt"
	"time"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/database"
)

var (
	sessions = make(map[int]*Session)
)

type Session struct {
	ID int

	*sshd.Terminal
	*database.UserProfile

	Created time.Time
}

func (session *Session) Print(a ...interface{}) error {
	_, err := session.Channel.Write([]byte(fmt.Sprint(a...)))
	return err
}

func (session *Session) Printf(format string, val ...any) error {
	_, err := session.Channel.Write([]byte(fmt.Sprintf(format, val...)))
	return err
}

func (session *Session) Println(a ...interface{}) error {
	_, err := session.Channel.Write([]byte(fmt.Sprint(a...) + "\r\n"))
	return err
}

func (session *Session) Clear() error {
	_, err := session.Channel.Write([]byte("\033c"))
	return err
}

func (session *Session) Close() {
	session.Conn.Close()
	session.Remove()
}

func SessionByName(name string) *Session {
	for _, s := range sessions {
		if s.UserProfile.Username != name {
			continue
		}

		return s
	}

	return nil
}

// Count returns the count of the sessions open
func Count() int {
	return len(sessions)
}

// Clone puts all the sessions into a slice
func Clone() []*Session {
	var list []*Session

	for _, session := range sessions {
		list = append(list, session)
	}

	return list
}
