package server

import(
	"net/http"
	"encoding/json"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"

	"github.com/gupta-yash4222/iitk-coin/model"
	"github.com/gupta-yash4222/iitk-coin/config/db"
)

type input struct {
	Rollno int `json:"rollno"`
	Name string `json:"name"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputData input
	var data model.User
	var res model.Response

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 
	}

	err = json.Unmarshal(body, &inputData)
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 
	}

	data.Rollno = inputData.Rollno
	data.Name = inputData.Name
	data.Password = inputData.Password
	data.Coins = 0

	_, err = db.FindUser(data.Rollno)

	if err != nil{

		if err.Error() == "sql: no rows in result set"{

			if len(data.Password) == 0{
				res.Error = "Password must contain at least one character"
				json.NewEncoder(w).Encode(res)
				return
			}

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