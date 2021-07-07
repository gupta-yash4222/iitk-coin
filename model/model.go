package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Rollno       int    `json:"rollno"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Coins        int    `json:"coins"`
	IsAdmin      int    `json:"isAdmin"`
	IsinCoreTeam int    `json:"isinCoreTeam"`
}

type Response struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type JWTclaims struct {
	Rollno       int    `json:"rollno"`
	Name         string `json:"name"`
	Admin        bool   `json:"admin"`
	IsinCoreTeam bool   `json:"isinCoreTeam"`
	jwt.StandardClaims
}

type Session struct {
	Rollno     int    `json:"rollno"`
	Token      string `json:"token"`
	IsLoggedIn bool   `json:"isloggedin"`
}

type TransferDetails struct {
	SenderRollno   int `json:"senderRollno"`
	ReceiverRollno int `json:"receiverRollno"`
	Coins          int `json:"coins"`
}

type TransactionDetails struct {
	Time            string `json:"time"`
	TransactionType string `json:"transactionType"`
	SenderRollno    int    `json:"senderRollno"`
	ReceiverRollno  int    `json:"receiverRollno"`
	Coins           int    `json:"coins"`
	Remarks         string `json:"remarks"`
}
