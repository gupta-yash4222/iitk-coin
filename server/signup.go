package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gupta-yash4222/iitk-coin/config/db"
	"github.com/gupta-yash4222/iitk-coin/model"
)

type input struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputData input
	var res model.Response

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

	batch := "Y" + strconv.Itoa(inputData.Rollno)[:2]

	data := model.User{
		Rollno:       inputData.Rollno,
		Name:         inputData.Name,
		Password:     inputData.Password,
		Coins:        0,
		Batch:        batch,
		IsAdmin:      0,
		IsinCoreTeam: 0,
		CanEarn:      1,
		NoOfEvents:   0,
	}

	_, err = db.FindUser(data.Rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			if len(data.Password) == 0 {
				res.Error = "Password must contain at least one character. Please enter a valid password."
				res.Result = "Registration unsuccessful"
				json.NewEncoder(w).Encode(res)
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

			if err != nil {
				res.Error = "Error while generating hash. Please try again."
				res.Result = "Registration unsuccessful"
				json.NewEncoder(w).Encode(res)
				return
			}

			data.Password = string(hash)

			err = db.AddUserData(data)
			if err != nil {
				res.Error = "Error while adding user in the database. Please try again"
				res.Result = "Registration unsuccessful"
				json.NewEncoder(w).Encode(res)
				return
			}

			res.Result = "User successfully registered"
			json.NewEncoder(w).Encode(res)
			return

		}

		res.Error = err.Error()
		res.Result = "Registration unsuccessful"
		json.NewEncoder(w).Encode(res)
		return

	}

	res.Error = "User already exists"
	json.NewEncoder(w).Encode(res)
	return

}
