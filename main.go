package main

import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type user struct{
	rollno int
	name string
}

func AddUserData(stmt *sql.Stmt, data user){
	
	stmt.Exec(data.rollno, data.name);

}

func FetchData(db *sql.DB){
	var data user
	rows, _ := db.Query("SELECT rollno, name FROM User")

	for rows.Next(){
		rows.Scan(&data.rollno, &data.name)
		fmt.Printf("rollno: %d, name: %s\n", data.rollno, data.name)
	}
}

func main(){

	database, _ := sql.Open( "sqlite3", "./user_details.db")

	statement, _ := database.Prepare( "CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT )")
	statement.Exec()

	stmt, _ := database.Prepare("INSERT INTO User (rollno, name) VALUES (?, ?)")

	data := user{rollno: 190998, name: "Yash Gupta"}
	AddUserData(stmt, data)

	FetchData(database)


}