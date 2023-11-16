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

func (s *Session) Print(a ...interface{}) error {
	_, err := s.Channel.Write([]byte(fmt.Sprint(a...)))
	return err
}

func (s *Session) Printf(format string, val ...any) error {
	_, err := s.Channel.Write([]byte(fmt.Sprintf(format, val...)))
	return err
}

func (s *Session) Println(a ...interface{}) error {
	_, err := s.Channel.Write([]byte(fmt.Sprint(a...) + "\r\n"))
	return err
}

func (s *Session) Clear() error {
	_, err := s.Channel.Write([]byte("\033c"))
	return err
}

func (s *Session) Close() {
	s.Conn.Close()
	s.Remove()
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
