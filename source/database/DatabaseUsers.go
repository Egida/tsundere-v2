package database

import (
	"database/sql"
	"strings"
	"time"
)

// CreateUser inserts a user in the database and also checks if it already exists
func CreateUser(user *UserProfile) error {
	if exists, err := Exists(user.Username); err != nil || exists {
		if exists {
			return ErrKnownUser
		}

		return err
	}

	return createUser(user)
}

// SetUser sets user information to something else
func SetUser(user *UserProfile) error {
	_, err := Instance.Exec("UPDATE users SET username = ?, password = ?, theme = ?, concurrents = ?, cooldown = ?, max_time = ?, max_sessions = ?, expiry = ?, roles = ?, created_by = ? WHERE username=?", user.Username, user.Password, user.Theme, user.Concurrents, user.Cooldown, user.MaxTime, user.MaxSessions, user.Expiry.Unix(), strings.Join(user.Roles, ","), user.CreatedBy, user.Username)
	return err
}

// Remove will remove the user from the database
func (u *UserProfile) Remove() error {
	// checks if the user is already deleted but this should in theory never happen
	if exists, err := Exists(u.Username); err != nil || !exists {
		if exists {
			return ErrUnknownUser
		}

		return err
	}

	// execute delete query
	_, err := Instance.Exec("DELETE FROM users WHERE username=?", u.Username)
	return err
}

/* ----------------------------------------------------------------------------------- */

// UserFromName gets a user from the database by the name as the method already says
func UserFromName(name string) (*UserProfile, error) {
	// select user from database with username filter
	rows, err := Instance.Query("SELECT * FROM users WHERE username=?", name)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// if there's no next row, we return an unknown user error.
	if !rows.Next() {
		return nil, ErrUnknownUser
	}

	// scan profile into UserProfile pointer
	return ScanUserProfile(rows)
}

// UserFromID gets a user from the database
func UserFromID(id int) (*UserProfile, error) {
	// select user from database with id filter
	rows, err := Instance.Query("SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// if there's no next row, we return an unknown user error.
	if !rows.Next() {
		return nil, ErrUnknownUser
	}

	// scan profile into UserProfile pointer
	return ScanUserProfile(rows)
}

// ScanUserProfile scans a user profile from the Sqlite3 database.
func ScanUserProfile(rows *sql.Rows) (*UserProfile, error) {
	var user UserProfile
	var expiry int64
	var roles string

	// scan user profile
	if err := rows.Scan(
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
	); err != nil {
		return nil, err
	}

	// convert specific units into our format from the UserProfile
	user.Roles = strings.Split(roles, ",")
	user.Expiry = time.Unix(expiry, 0)

	return &user, nil
}

// Users will get all users from the database in a slice
func Users() ([]*UserProfile, error) {
	var users []*UserProfile

	// select all users with all information
	rows, err := Instance.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	// iterate through all rows and add users into a slice
	for rows.Next() {
		user, err := ScanUserProfile(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Exists will check if the user exists in the database
func Exists(name string) (bool, error) {
	// select id by username (even tho we dont need it)
	result, err := Instance.Query("SELECT id FROM users WHERE username=?", name)
	if err != nil {
		return false, err
	}

	defer closeRows(result)

	return result.Next(), nil
}

// VerifyCredentials tries to find a matching username and password from the database
func VerifyCredentials(name, password string) error {
	// get user by username
	user, err := UserFromName(name)
	if err != nil {
		return err
	}

	// if the user is nil it's an unknown user
	if user == nil {
		return ErrUnknownUser
	}

	// compares password with the one from the database
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

	defer closeRows(rows)
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

// createUser will create a user without any checks, this will not be used outside this package
func createUser(user *UserProfile) error {
	_, err := Instance.Exec("INSERT INTO users (username, password, theme, concurrents, cooldown, max_time, max_sessions, expiry, roles, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Username, user.Password, user.Theme, user.Concurrents, user.Cooldown, user.MaxTime, user.MaxSessions, user.Expiry.Unix(), strings.Join(user.Roles, ","), user.CreatedBy)
	return err
}
