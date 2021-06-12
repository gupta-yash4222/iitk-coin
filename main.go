package main

import(
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/gupta-yash4222/iitk-coin/server"
	"github.com/gupta-yash4222/iitk-coin/config/db"
	
)

func main(){
	/*
	var data model.User

	data.Rollno = 190998
	data.Name = "Yash Gupta"

	err := db.AddUserData(data)

	if err != nil{
		fmt.Println(err.Error())
	}
	*/

	r := mux.NewRouter()

	r.HandleFunc("/hihi", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Aa gaye meri maut ka tamasha dekhne")
	})

	r.HandleFunc("/signup", server.RegisterUser).Methods("POST")

	r.HandleFunc("/getUsers", db.FetchUserDataServer).Methods("GET")

	r.HandleFunc("/login", server.LoginUser).Methods("POST")

	r.HandleFunc("/secretpage", server.WelcomeUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))

	/*

	database, _ := sql.Open("sqlite3", "user_details.db")

	stmt, _ := database.Prepare( "CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT, password TEXT )")
	stmt.Exec()

	stmt, _ = database.Prepare("INSERT INTO User (rollno, name, password) VALUES (?, ?, ?)")

	stmt.Exec(500998, "Psycic Clown", "jester_dominion")

	db.FetchUserData()

	*/

	db.FetchUserDataTerminal()

	


}