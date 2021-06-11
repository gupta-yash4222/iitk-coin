package model

import( jwt"github.com/dgrijalva/jwt-go" )

type User struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type JWTclaims struct {
	Rollno int    `json:"rollno"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

type Session struct {
	Rollno     int    `json:"rollno"`
	Token      string `json:"token"`
	IsLoggedIn bool   `json:"isloggedin"`
}