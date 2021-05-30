package main

import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type user struct{
	rollno int
	name string
}

func AddUserData(db *sql.DB, data user){

	stmt, _ := db.Prepare("INSERT INTO User (rollno, name) VALUES (?, ?)")
	stmt.Exec(data.rollno, data.name);

}

func main(){

	database, _ := sql.Open( "sqlite3", "./user_details.db")

	statement, _ := database.Prepare( "CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT )")
	statement.Exec()

	data := user{rollno: 190998, name: "Yash Gupta"}

	AddUserData(database, data)

}