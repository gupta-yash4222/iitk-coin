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
var validDuration time.Duration = 25

func LoginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputData input
	var res model.Response

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &inputData)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	data, err := db.FindUser(inputData.Rollno)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			res.Result = "User not registered"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if inputData.Name != data.Name {
		res.Result = "Invalid Username" // given username doesn't matches with the registered username
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(inputData.Password))

	if err != nil {
		res.Result = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "User authenticated"
	//json.NewEncoder(w).Encode(res)

	var result model.Session

	expirationTime := time.Now().Add(validDuration * time.Second)

	admin := false
	if inputData.Rollno == 190998 {
		admin = true
	}

	userClaims := model.JWTclaims{
		Rollno:         inputData.Rollno,
		Name:           inputData.Name,
		Admin:          admin,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		res.Error = "Error while generating JSON Web Token. Please try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Rollno = inputData.Rollno
	result.Token = tokenString
	result.IsLoggedIn = true

	//json.NewEncoder(w).Encode(res)

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(result)

}
