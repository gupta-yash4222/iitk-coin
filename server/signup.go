package server

import(
	"net/http"
	"encoding/json"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"

	"github.com/gupta-yash4222/iitk-coin/model"
	"github.com/gupta-yash4222/iitk-coin/config/db"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data model.User
	var res model.Response

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 
	}

	err = json.Unmarshal(body, &data)
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 
	}

	_, err = db.FindUser(data.Rollno)

	if err != nil{

		if err.Error() == "sql: no rows in result set"{

			hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

			if err != nil{
				res.Error = "Error while generating hash. Please try again."
				json.NewEncoder(w).Encode(res)
				return 
			}

			data.Password = string(hash)

			err = db.AddUserData(data)
			if err != nil{
				res.Error = "Error while adding user in the database. Please try again"
				json.NewEncoder(w).Encode(res)
				return 
			}

			res.Result = "User successfully registered"
			json.NewEncoder(w).Encode(res)
			return 

		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 

	}

	res.Error = "User already exists"
	json.NewEncoder(w).Encode(res)
	return 

}