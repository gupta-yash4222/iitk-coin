package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gupta-yash4222/iitk-coin/config/db"
	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

type RedeemRequest struct {
	ItemId string `json:"itemId"`
	Coins int `json:"coins"`
}

// Make a redeem request with item ID and no. of coins specified
func RedeemCoins(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var res model.Response

	claims, err := ValidateUser(r)

	if err != nil {
		if err.Error() == "Token is either expired or not active yet" {
			res.Error = "Session expired. Login again."
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		res.Result = "Redeem request aborted"
		json.NewEncoder(w).Encode(res)
		return
	}

	//fmt.Println(claims.Name)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var inputData RedeemRequest
	err = json.Unmarshal(body, &inputData)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res = db.RedeemHandler(claims.Rollno, inputData.ItemId, inputData.Coins)

	json.NewEncoder(w).Encode(res)
	return
}

// Admin verifies a redeem request by specifying its ID and then either approving it or rejecting it, or having it on a pending status
func VerifyRedeemRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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
		res.Error = "User not authorized for the action"
		res.Result = "Permission denied"
		json.NewEncoder(w).Encode(res)
		return 
	}

	var inputData map[string]interface{}

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

	redeemId := int( inputData["redeemId"].(float64) )
	action := inputData["action"].(string)

	res = db.RedeemRequestVerification(redeemId, action)

	json.NewEncoder(w).Encode(res)

}