package database

import (
	"errors"
	"time"
	"tsundere/source/models/roles"
)

var (
	ErrKnownUser          = errors.New("known user")
	ErrUnknownUser        = errors.New("unknown user")
	ErrInvalidCredentials = errors.New("invalid")
)

type UserProfile struct {
	Id          int       `swash:"id"`
	Username    string    `swash:"username"`
	Password    string    `swash:"password"`
	Theme       string    `swash:"theme"`
	Concurrents int       `swash:"concurrents"`
	Cooldown    int       `swash:"cooldown"`
	MaxTime     int       `swash:"max_time"`
	MaxSessions int       `swash:"max_sessions"`
	Expiry      time.Time `swash:"expiry"`
	Roles       []string  `swash:"roles"`
	CreatedBy   string    `swash:"created_by"`
}

func (u *UserProfile) HasRole(name string) bool {
	for _, role := range u.Roles {
		_, exists := roles.List[name]
		if exists && role == name {
			return true
		}
	}

	return false
}
