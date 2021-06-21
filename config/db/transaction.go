package db 

import (
	_"database/sql"
	"fmt"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

// Add given number of coins in the specified user's account
func AddCoins(rollno int, coins int) model.Response {

	var res model.Response

	_, err := FindUser(rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", rollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	result, err := Database.Exec("UPDATE User SET coins = coins + ? WHERE rollno = ?", coins, rollno)
	rowsAffected, _ := result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		
		if err != nil {
			res.Error = err.Error()
			res.Result = "Transaction aborted"
			return res
		}

		if rowsAffected != 1 {
			res.Error = "Unexpected error"
			res.Result = "Transaction aborted"
			return res
		}
	}

	res.Result = "Transaction successful"
	return res
}

// Transfer given number of coins, if possible, from sender's account to the receiver's account
func TransferCoins(data model.TransferDetails) model.Response {

	var res model.Response

	_, err := FindUser(data.SenderRollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", data.SenderRollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	_, err = FindUser(data.ReceiverRollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", data.ReceiverRollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	tx, err := Database.Begin()
	if err != nil {
		res.Error = err.Error()
		res.Result = "Transaction aborted"
		return res
	}

	result, err := tx.Exec("UPDATE User SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0", data.Coins, data.SenderRollno, data.Coins)
	rowsAffected, _ := result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			res.Error = err.Error()
			res.Error = "Transaction aborted"
			return res
		}

		if rowsAffected == 0 {
			res.Error = "Insufficient balance"
		} else {
			res.Error = "Unexpected error"
		}
		res.Result = "Transaction aborted"
		return res
	}

	result, err = tx.Exec("UPDATE User SET coins = coins + ? WHERE rollno = ?", data.Coins, data.ReceiverRollno)
	rowsAffected, _ = result.RowsAffected()

	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			res.Error = err.Error()
			res.Error = "Transaction aborted"
			return res
		}

		if rowsAffected == 0 {
			res.Error = "Insufficient balance"
		} else {
			res.Error = "Unexpected error"
		}
		res.Result = "Transaction aborted"
		return res
	}

	err = tx.Commit()
	if err != nil {
		res.Error = err.Error()
		res.Result = "Transaction aborted"
		return res
	}

	res.Result = "Transaction successful"
	return res

}