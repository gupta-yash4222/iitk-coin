package server

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"net/http"

	"github.com/gupta-yash4222/iitk-coin/config/db"
	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

func AddItem(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputData model.Item
	var res model.Response

	claims, err := ValidateUser(r)

	if err != nil {
		if err.Error() == "Token is either expired or not active yet" {
			res.Error = "Session expired. Login again."
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		res.Result = "Verification aborted"
		json.NewEncoder(w).Encode(res)
		return
	}

	if !claims.Admin {
		res.Error = "User not authorized to add items into the inventory"
		res.Result = "Permission denied"
		json.NewEncoder(w).Encode(res)
		return 
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &inputData)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = db.AddItem(inputData)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Item added into the inventory successfully"
	json.NewEncoder(w).Encode(res)
	return
}