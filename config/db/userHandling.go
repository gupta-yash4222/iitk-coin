package db

import (
	_"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

// Finds a user with the given rollno
func FindUser(roll int) (model.User, error) {  

	var data model.User
	
	err := Database.QueryRow("SELECT rollno, name, password, coins FROM User WHERE rollno = ?", roll).Scan(&data.Rollno, &data.Name, &data.Password, &data.Coins)
	if err != nil{
		return data, err
	}

	return data, nil
}

// Adds a user with the given credentials if not already present in the database
func AddUserData(data model.User) error { 

	_, err := FindUser(data.Rollno)

	if err != nil{

		if err.Error() == "sql: no rows in result set" {
			stmt, err := Database.Prepare("INSERT INTO User (rollno, name, password, coins) VALUES (?, ?, ?, ?)")
			if err != nil{
				return err
			}

			stmt.Exec(data.Rollno, data.Name, data.Password, data.Coins);
			return nil
		}

		return err
	}

	return errors.New("User already present")

}

// Fetch all registered users from the database and output them on the terminal
func FetchUserDataTerminal() error {   

	rows, err := Database.Query("SELECT rollno, name FROM User")
	if err != nil{
		return err
	}

	var data model.User
	for rows.Next(){
		err = rows.Scan(&data.Rollno, &data.Name)
		if err != nil{
			return err
		}
		fmt.Printf("Rollno.: %d, Name: %s\n", data.Rollno, data.Name)
	}

	return nil
}

// Fetch all registered users from the database and output them as a response
func FetchUserDataServer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := Database.Query("SELECT rollno, name FROM User")
	if err != nil{
		log.Fatal(err)
		return
	}

	var data model.User
	for rows.Next(){
		err = rows.Scan(&data.Rollno, &data.Name)
		if err != nil{
			log.Fatal(err)
			return
		}
		fmt.Fprintf(w, "Rollno.: %d, Name: %s\n", data.Rollno, data.Name)
	}

}

// Delete a user from the database 
func DeleteUser(roll int) error {

	_, err := FindUser(roll)
	if err != nil{
		log.Fatal("User not in the database")
		return err
	}

	_, err = Database.Exec("DELETE FROM User WHERE rollno = ?", roll)
	if err != nil{
		return err
	}

	fmt.Println("User deleted successfully")

	return nil

}
