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

func RewardCoins(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	res = db.AddCoins(data.Rollno, data.Coins)
	json.NewEncoder(w).Encode(res)

}

func FetchUserBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var res model.Response

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

	res = db.TransferCoins(inputData)
	json.NewEncoder(w).Encode(res)

}

