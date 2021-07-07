package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gupta-yash4222/iitk-coin/config/db"
	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

type RewardDetails struct {
	Rollno int `json:"rollno"`
	Coins  int `json:"coins"`
}

// Give the specified user given number of coins as a reward when participated in an event
func RewardCoins(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var res model.Response

	// validating the request if it is made by an admin or not
	claims, err := ValidateUser(r)
	if err != nil {

		if err.Error() == "Token is either expired or not active yet" {
			res.Error = err.Error()
			res.Result = "Login again"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		res.Result = "Transaction aborted"
		json.NewEncoder(w).Encode(res)
		return
	}

	if !claims.Admin {
		res.Error = "User not authorized for the action"
		res.Result = "Action denied"
		json.NewEncoder(w).Encode(res)
		return
	}

	var data RewardDetails

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		//log.Fatal(err)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		//log.Fatal(err)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res = db.AddCoins(data.Rollno, data.Coins)
	json.NewEncoder(w).Encode(res)

}

// Fetch the specified user balance from the database
func FetchUserBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rolls, ok := r.URL.Query()["rollno"]
	if !ok || len(rolls[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rollno, err := strconv.Atoi(rolls[0])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := db.FindUser(rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			fmt.Fprintf(w, "User with rollno %d is not registered.", rollno)
			return
		}

		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "You currently have %d coins ", data.Coins)
}

// Transfer given number of coins from the sender's account to the receiver's account
func TransferCoins(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var inputData model.TransferDetails
	var res model.Response

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &inputData)

	// validating the sender if he/she is logged in or not
	claims, err := ValidateUser(r)
	if err != nil {

		if err.Error() == "Token is either expired or not active yet" {
			res.Error = err.Error()
			res.Result = "Login again"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		res.Result = "Transaction aborted"
		json.NewEncoder(w).Encode(res)
		return
	}

	if claims.Rollno != inputData.SenderRollno {
		res.Error = "User not authorized for the action"
		res.Result = "Action denied"
		json.NewEncoder(w).Encode(res)
		return
	}

	res = db.TransferCoins(inputData)
	json.NewEncoder(w).Encode(res)

}

