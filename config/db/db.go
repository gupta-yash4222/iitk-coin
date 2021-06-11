package db

import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"errors"
	"github.com/gupta-yash4222/iitk-coin/model"
)

func FindUser(roll int) (model.User, error) {

	var data model.User

	db, err := sql.Open( "sqlite3", "user_details.db")
	if err != nil{
		//fmt.Println("Database could not be opened or created")
		return data, err
	}

	err = db.QueryRow("SELECT rollno, name, password FROM User WHERE rollno = ?", roll).Scan(&data.Rollno, &data.Name, &data.Password)

	if err != nil{
		//fmt.Println("User not found")
		return data, err
	}

	//fmt.Printf("Rollno.: %d, Name: %s\n", data.Rollno, data.Name)

	return data, nil
}

func AddUserData(data model.User) error {

	db, err := sql.Open( "sqlite3", "user_details.db")

	if err != nil{
		//fmt.Println("Database could not be opened or created")
		return err
	}

	stmt, err := db.Prepare( "CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT, password TEXT )")
	stmt.Exec()

	_, err = FindUser(data.Rollno)

	if err != nil{

		if err.Error() == "sql: no rows in result set" {
			stmt, err = db.Prepare("INSERT INTO User (rollno, name, password) VALUES (?, ?, ?)")
			if err != nil{
				//fmt.Println("Data could not be inserted")
				return err
			}

			stmt.Exec(data.Rollno, data.Name, data.Password);
			return nil
		}

		return err
	}

	return errors.New("User already present")

}

func FetchUserData() error {
	db, err := sql.Open("sqlite3", "user_details.db")
	
	if err != nil{
		//fmt.Println("Database could not be opened or created")
		return err
	}

	rows, err := db.Query("SELECT rollno, name FROM User")

	if err != nil{
		//fmt.Println("Data could not be fetched")
		return err
	}

	var data model.User
	for rows.Next(){
		err = rows.Scan(&data.Rollno, &data.Name)
		if err != nil{
			//fmt.Println("Rows could not be fetched")
			return err
		}
		fmt.Printf("Rollno.: %d, Name: %s\n", data.Rollno, data.Name)
	}

	return nil
}

func DeleteUser(roll int) error {
	db, err := sql.Open( "sqlite3", "user_details.db")

	if err != nil{
		return err
	}

	_, err = FindUser(roll)
	if err != nil{
		log.Fatal("User not in the database")
		return err
	}

	_, err = db.Exec("DELETE FROM User WHERE rollno = ?", roll)

	if err != nil{
		return err
	}

	fmt.Println("User deleted successfully")

	return nil

}
