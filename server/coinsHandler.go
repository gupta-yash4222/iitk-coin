package server

import (
	"database/sql"
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

type TransferDetails struct {
	SenderRollno int `json:"senderRollno"`
	ReceiverRollno int `json:"receiverRollno"`
	Coins int `json:"coins"`
}

func RewardCoins(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var data RewardDetails
	var res model.Response

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

	_, err = db.FindUser(data.Rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", data.Rollno)
			res.Result = "Transaction aborted"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	database, err := sql.Open("sqlite3", "user_details.db")
	if err != nil {
		//log.Fatal(err)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	result, err := database.Exec("UPDATE User SET coins = coins + ? WHERE rollno = ?", data.Coins, data.Rollno)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(result.RowsAffected())

	res.Result = "Transaction successful"
	json.NewEncoder(w).Encode(res)

}

func FetchUserBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//var inputData map[string]interface{}
	var res model.Response

	/*
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		err = json.Unmarshal(body, &inputData)

		rollno := int(inputData["rollno"].(float64))
	*/

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

	//rollno := 500998

	data, err := db.FindUser(rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = "User not registered"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Fprintf(w, "You currently have %d coins", data.Coins)
}

