package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"


	"github.com/gupta-yash4222/iitk-coin/model"
)

func WelcomeUser(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//tokenString := r.Header.Get("Authorization")
	//splitToken := strings.Split(tokenString, "Bearer ")
	//tokenString = splitToken[1]

	var res model.Response

	cookie, err := r.Cookie("token")
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	tokenString := cookie.Value

	tokenClaims := &model.JWTclaims{}

	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error){
		// Validating the algorithm used in producing the JWT
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		return jwtKey, nil
	})

	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if token.Valid{

		if !tokenClaims.Admin{
			res.Result = "Access denied"
			json.NewEncoder(w).Encode(res)
			return
		}

		fmt.Fprintf(w, "Everyone welcome %s. We hope you have brought pizzas.", tokenClaims.Name)

		res.Result = "Access granted"
		json.NewEncoder(w).Encode(res)
		return

	} else if ve, ok := err.(*jwt.ValidationError); ok{
		
		if ve.Errors & jwt.ValidationErrorMalformed != 0 {
			res.Error = "The token is malformed!"
			json.NewEncoder(w).Encode(res)
			return
		}

		if ve.Errors & jwt.ValidationErrorExpired != 0{
			res.Error = "The token has expired. Log in your account again"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = "Couldn't handle the token: " + err.Error()
		json.NewEncoder(w).Encode(res)
		return

	}else{
		res.Error = "Couldn't handle the token: " + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}