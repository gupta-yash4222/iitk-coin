package db 

import (
	_"database/sql"
	"fmt"
	"strconv"
	"math"
	"time"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

func AddTransactionDetails(data model.TransactionDetails) error {

	stmt, err := Database.Prepare("INSERT INTO Transaction_Log (time, transactionType, senderRollno, receiverRollno, coins, remarks) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(data.Time, data.TransactionType, data.SenderRollno, data.ReceiverRollno, data.Coins, data.Remarks)
	return nil
}

// Add given number of coins in the specified user's account
func AddCoins(rollno int, coins int) model.Response {

	var res model.Response

	user, err := FindUser(rollno)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", rollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	if user.IsAdmin == 1 || user.IsinCoreTeam == 1 {
		res.Error = "User not allowed to receive a reward"
		res.Result = "Transaction aborted"
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

	transaction := model.TransactionDetails{
		Time: time.Now().String(),
		TransactionType: "Reward",
		SenderRollno: 0,
		ReceiverRollno: rollno,
		Coins: coins,
		Remarks: "Coins rewarded successfully",
	}

	err = AddTransactionDetails(transaction)

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

	roll1, err := strconv.Atoi(strconv.Itoa(data.SenderRollno)[:2])
	if err != nil {
		fmt.Println(err.Error())
		res.Error = err.Error()
		res.Result = "Transaction aborted"
		return res
	}

	roll2, err := strconv.Atoi(strconv.Itoa(data.ReceiverRollno)[:2])
	if err != nil {
		fmt.Println(err.Error())
		res.Error = err.Error()
		res.Result = "Transaction aborted"
		return res
	}

	// imposing the required tax on the transaction
	var coins int
	if roll1 == roll2 {
		coins = int(math.Round(0.98 * float64(data.Coins)))
	} else {
		coins = int(math.Round(0.77 * float64(data.Coins)))
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

	result, err = tx.Exec("UPDATE User SET coins = coins + ? WHERE rollno = ?", coins, data.ReceiverRollno)
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

	transaction := model.TransactionDetails{
		Time: time.Now().String(),
		TransactionType: "Coin Transfer",
		SenderRollno: data.SenderRollno,
		ReceiverRollno: data.ReceiverRollno,
		Coins: data.Coins,
		Remarks: "Coins transferred successfully",
	}

	err = AddTransactionDetails(transaction)

	res.Result = "Transaction successful"
	return res

}