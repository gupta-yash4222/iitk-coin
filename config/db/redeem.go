package db

import (
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

func AddRedeemDetails(data model.RedeemDetails) error {

	stmt, err := Database.Prepare("INSERT INTO RedeemLog (time, rollno, itemId, coins, status, remarks) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(data.Time, data.Rollno, data.ItemId, data.Coins, data.Status, data.Remarks)
	return nil
}

func FindRedeemRequest(id int) (model.RedeemDetails, error) {

	var data model.RedeemDetails

	err := Database.QueryRow("SELECT time, rollno, itemId, coins, status, remarks FROM RedeemLog WHERE id = ?", id).Scan(&data.Time, &data.Rollno, &data.ItemId, &data.Coins, &data.Status, &data.Remarks)
	if err != nil {
		return data, err
	}

	return data, nil
}

func FetchRedeemRequests(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := Database.Query("SELECT id, time, rollno, itemId, coins, status, remarks FROM RedeemLog")
	if err != nil {
		log.Fatal(err)
		return
	}

	var data model.RedeemDetails
	for rows.Next() {
		err = rows.Scan(&data.Id, &data.Time, &data.Rollno, &data.ItemId, &data.Coins, &data.Status, &data.Remarks)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Fprintf(w, "Id: %d, Time: %v, Rollno: %d, Item ID: %s, Coins, %d, Status: %s, Remarks: %s\n", data.Id, data.Time, data.Rollno, data.ItemId, data.Coins, data.Status, data.Remarks)
	}
}

func RedeemHandler(rollno int, itemId string, coins int) model.Response {

	var res model.Response

	user, err := FindUser(rollno)
	if err != nil {
		res.Error = err.Error()
		return res
	}

	if coins > user.Coins {
		res.Error = "User has insufficient balance"
		res.Result = "Redeem request aborted"
		return res
	}

	item, err := FindItem(itemId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			res.Error = "Item not found in the database. Please enter an available Item ID"
			res.Result = "Redeem request aborted"
			return res
		}
		res.Error = err.Error()
		return res
	}

	if item.Quantity == 0 {
		res.Error = "Entered item is currently not available in the inventory"
		res.Result = "Redeem request aborted"
		return res
	}

	if coins < item.Price {
		res.Error = "Price of the item is more than the entered number of coins. Please enter valid number of coins."
		res.Result = "Redeem request aborted"
		return res
	}

	redeem := model.RedeemDetails{
		Time:    time.Now().String(),
		Rollno:  rollno,
		ItemId:  itemId,
		Coins:   coins,
		Status:  "Pending",
		Remarks: "Redeem Request not accepted yet",
	}

	err = AddRedeemDetails(redeem)

	res.Result = "Redeem request successfully registered. Wait for a couple of days to get it accepted."
	return res
}

func RedeemRequestVerification(id int, action string) model.Response {

	var res model.Response

	redeemData, err := FindRedeemRequest(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			res.Error = "Item not found in the database. Please enter an available Item ID"
			res.Result = "Redeem request aborted"
			return res
		}
		res.Error = err.Error()
		return res
	}

	if redeemData.Status == "Accepted" {
		if action == "Approve" {
			res.Error = "Redeem request is already approved"
		} else {
			res.Error = "Redeem request is already approved. Changes cannot be reverted"
		}
		res.Result = "Verification aborted"
		return res
	}

	var remarks string
	var status string
	if action == "Approve" {
		status = "Accepted"
		remarks = "Redeem request approved"
	} else if action == "Reject" {
		status = "Rejected"
		remarks = "Redeem request rejected"
	} else if action == "Put on Hold" {
		status = "Pending"
		remarks = "Redeem Request not accepted yet"
	} else {
		res.Error = "Invalid type of action"
		res.Result = "Please enter a valid type of action"
	}

	tx, err := Database.Begin()
	if err != nil {
		res.Error = err.Error()
		res.Result = "Verification aborted"
		return res
	}

	result, err := tx.Exec("UPDATE RedeemLog SET status = ?, remarks = ? WHERE id = ?", status, remarks, id)
	rowsAffected, _ := result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			res.Error = err.Error()
			res.Result = "Verification aborted"
			return res
		}

		if rowsAffected == 0 {
			res.Error = "Database could not be updated"
		} else {
			res.Error = "Unexpected error"
		}
		res.Result = "Verification aborted"
		return res
	}

	result, err = tx.Exec("UPDATE Items SET quantity = quantity - 1 WHERE itemId = ? AND quantity - 1 >= 0", redeemData.ItemId)
	rowsAffected, _ = result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			res.Error = err.Error()
			res.Result = "Verification aborted"
			return res
		}

		if rowsAffected == 0 {
			res.Error = "Item not present in sufficient numbers"
		} else {
			res.Error = "Unexpected error"
		}
		res.Result = "Verification aborted"
		return res
	}

	result, err = tx.Exec("UPDATE User SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0", redeemData.Coins, redeemData.Rollno, redeemData.Coins)
	rowsAffected, _ = result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			res.Error = err.Error()
			res.Result = "Verification aborted"
			return res
		}

		if rowsAffected == 0 {
			res.Error = "User has insufficient balance"
		} else {
			res.Error = "Unexpected error"
		}
		res.Result = "Verification aborted"
		return res
	}

	err = tx.Commit()
	if err != nil {
		res.Error = err.Error()
		res.Result = "Verification aborted"
		return res
	}

	res.Result = "Verification successful"
	return res

}
