package database

import (
	"database/sql"
	"log"
)

func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
