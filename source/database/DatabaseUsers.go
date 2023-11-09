package database

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func CreateUser(user *UserProfile, existsCheck bool) error {
	if existsCheck {
		exists, err := Exists(user.Username)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if exists {
			return ErrKnownUser
		}

		return createUser(user)
	}
	return createUser(user)
}

func (u *UserProfile) Remove() error {
	exists, err := Exists(u.Username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUnknownUser
	}
	_, err = Instance.Exec("DELETE FROM users WHERE username=?", u.Username)
	return err
}

/* ----------------------------------------------------------------------------------- */

func UserFromName(name string) (*UserProfile, error) {
	rows, err := Instance.Query("SELECT * FROM users WHERE username=?", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var user UserProfile
		var expiry int64
		var roles string

		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Theme,
			&user.Concurrents,
			&user.Cooldown,
			&user.MaxTime,
			&user.MaxSessions,
			&expiry,
			&roles,
			&user.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		user.Roles = strings.Split(roles, ",")
		user.Expiry = time.Unix(expiry, 0)

		return &user, nil
	}

	return nil, errors.New("user does not exist")
}

func UserFromID(id int) (*UserProfile, error) {
	rows, err := Instance.Query("SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var user UserProfile
		var expiry int64
		var roles string

		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Theme,
			&user.Concurrents,
			&user.Cooldown,
			&user.MaxTime,
			&user.MaxSessions,
			&expiry,
			&roles,
			&user.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		user.Roles = strings.Split(roles, ",")
		user.Expiry = time.Unix(expiry, 0)

		return &user, nil
	}

	return nil, errors.New("user does not exist")
}

func Users() ([]*UserProfile, error) {
	var users []*UserProfile

	rows, err := Instance.Query("SELECT id FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		user, err := UserFromID(id)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func Exists(name string) (bool, error) {
	result, err := Instance.Query("SELECT * FROM users WHERE username=?", name)
	if err != nil {
		return false, err
	}

	defer result.Close()
	return result.Next(), nil
}

func VerifyCredentials(name, password string) error {
	user, err := UserFromName(name)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUnknownUser
	}

	if user.Password != password {
		return ErrInvalidCredentials
	}

	return nil
}

func (u *UserProfile) LastFlood() (*Flood, error) {
	rows, err := Instance.Query("SELECT id FROM logs WHERE rowid = (SELECT MAX(rowid) FROM logs WHERE user_id = ?)", u.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		flood, err := FloodByID(id)
		if err != nil {
			return nil, err
		}

		return flood, nil
	}

	return nil, nil
}

func (u *UserProfile) Floods() ([]*Flood, error) {
	var floods []*Flood

	rows, err := Instance.Query("SELECT id FROM logs WHERE user_id = ?", u.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		flood, err := FloodByID(id)
		if err != nil {
			return nil, err
		}

		floods = append(floods, flood)
	}

	return floods, nil
}

func (u *UserProfile) RunningAttacks() ([]*Flood, error) {
	var runningFloods []*Flood

	floods, err := u.Floods()
	if err != nil {
		return nil, err
	}

	for _, flood := range floods {
		if !flood.End.After(time.Now()) {
			continue
		}

		runningFloods = append(runningFloods, flood)
	}

	return runningFloods, nil
}

func createUser(user *UserProfile) error {
	_, err := Instance.Exec("INSERT INTO users (username, password, theme, concurrents, cooldown, max_time, max_sessions, expiry, roles, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Username, user.Password, user.Theme, user.Concurrents, user.Cooldown, user.MaxTime, user.MaxSessions, user.Expiry.Unix(), strings.Join(user.Roles, ","), user.CreatedBy)
	return err
}

func SetUser(user *UserProfile) error {
	_, err := Instance.Exec("UPDATE users SET username = ?, password = ?, theme = ?, concurrents = ?, cooldown = ?, max_time = ?, max_sessions = ?, expiry = ?, roles = ?, created_by = ? WHERE username=?", user.Username, user.Password, user.Theme, user.Concurrents, user.Cooldown, user.MaxTime, user.MaxSessions, user.Expiry.Unix(), strings.Join(user.Roles, ","), user.CreatedBy, user.Username)
	return err
}
