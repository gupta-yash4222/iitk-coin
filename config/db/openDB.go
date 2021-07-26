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

	stmt, err := database.Prepare("CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT, password TEXT, coins INTEGER, batch TEXT, isAdmin INTEGER, isinCoreTeam INTEGER, canEarn INTEGER, noOfEvents INTEGER )")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	stmt, err = database.Prepare("CREATE TABLE IF NOT EXISTS Items ( itemId TEXT PRIMARY KEY, itemDescription TEXT, quantity INTEGER, price INTEGER )")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	stmt, err = database.Prepare("CREATE TABLE IF NOT EXISTS TransferLog ( time TEXT, senderRollno INTEGER, receiverRollno INTEGER, coins INTEGER, remarks TEXT )")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	stmt, err = database.Prepare("CREATE TABLE IF NOT EXISTS RewardLog ( time TEXT, receiverRollno INTEGER, coins INTEGER, remarks TEXT )")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	stmt, err = database.Prepare("CREATE TABLE IF NOT EXISTS RedeemLog ( id INTEGER PRIMARY KEY AUTOINCREMENT, time TEXT, rollno INTEGER, itemId TEXT, coins INTEGER, status TEXT, remarks TEXT )")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	return database
}
