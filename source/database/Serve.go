package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var Instance *sql.DB

// Configure will configure the database.
func Configure() error {
	conn, err := sql.Open("sqlite3", "resources/database.db")
	if err != nil {
		return err
	}

	Instance = conn

	err = CreateUserTable()
	if err != nil {
		return err
	}

	err = CreateLogsTable()
	if err != nil {
		return err
	}

	return serve()
}

func serve() error {
	err := CreateUser(&UserProfile{
		Username:    "root",
		Password:    "123",
		Theme:       "default",
		Concurrents: 100,
		Cooldown:    0,
		MaxTime:     60,
		MaxSessions: 10,
		Expiry:      time.Now().Add(99999 * time.Hour),
		Roles:       []string{"admin"},
		CreatedBy:   "SQLite3",
	}, true)
	if err != nil && !errors.Is(err, ErrKnownUser) {
		return err
	}

	log.Printf("Database has been parsed & configured successfully")
	return nil
}
