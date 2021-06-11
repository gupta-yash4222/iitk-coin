package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/gupta-yash4222/iitk-coin/config/db"
	"github.com/gupta-yash4222/iitk-coin/model"
)

var jwtKey = []byte("CROWmium")
var validDuration int = 10


func LoginUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	var res model.Response

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	data, err := db.FindUser(user.Rollno)

	if err != nil{
		if err.Error() == "sql: no rows in result set"{
			res.Result = "User not registered"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if user.Name != data.Name{
		res.Result = "Invalid Username"  // given username doesn't matches with the registered username
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password))

	if err != nil{
		res.Result = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "User authenticated"
	//json.NewEncoder(w).Encode(res)

	var result model.Session

	expirationTime := time.Now().Add(25 * time.Second)

	userClaims := model.JWTclaims{
		Rollno: user.Rollno,
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil{
		res.Error = "Error while generating JSON Web Token. Please try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Rollno = user.Rollno
	result.Token = tokenString
	result.IsLoggedIn = true

	//json.NewEncoder(w).Encode(res)
	json.NewEncoder(w).Encode(result)

}