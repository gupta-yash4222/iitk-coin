package db

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func OpenDatabase() *sql.DB {
	database, err := sql.Open("sqlite3", "user_details.db")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := database.Prepare("CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT, password TEXT, coins INTEGER )")
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	return database
}
