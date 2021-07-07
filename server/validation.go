package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"


	"github.com/gupta-yash4222/iitk-coin/model"
)

// Validating the JSON Web Token received with the cookie 
func ValidateUser(r *http.Request) (*model.JWTclaims, error) {

	//tokenString := r.Header.Get("Authorization")
	//splitToken := strings.Split(tokenString, "Bearer ")
	//tokenString = splitToken[1]

	cookie, err := r.Cookie("token")
	if err != nil{
		return nil, err
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

	if token.Valid{

		return tokenClaims, nil

	} else if ve, ok := err.(*jwt.ValidationError); ok{
		
		if ve.Errors & jwt.ValidationErrorMalformed != 0 {
			return tokenClaims, errors.New("The token is malformed!")

		} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
			//fmt.Println(err.Error())
			return tokenClaims, errors.New("Token is either expired or not active yet")
			
		} else {
			return tokenClaims, errors.New("Couldn't handle the token: " + err.Error())
		
		}

	}else{
		return tokenClaims, errors.New("Couldn't handle the token: " + err.Error())
		
	}

}

func WelcomeUser(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var res model.Response

	claims, err := ValidateUser(r)
	if err != nil{

		if err.Error() == "Token is either expired or not active yet" {
			res.Error = err.Error()
			res.Result = "Login again"
			json.NewEncoder(w).Encode(res)
		}

		res.Error = err.Error()
		res.Result = "Could not access the page"
		json.NewEncoder(w).Encode(res)
		return
	}

	if !claims.Admin && !claims.IsinCoreTeam {
		res.Error = "User not authorized to access the page"
		res.Result = "Access denied"
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Fprintf(w, "Everyone welcome %s. We hope you have brought pizzas.", claims.Name)

	res.Result = "Access granted"
	json.NewEncoder(w).Encode(res)
	return
	
}