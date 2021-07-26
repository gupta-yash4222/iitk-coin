package db

import (
	_ "database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

func AddTransactionDetails(data model.TransactionDetails) error {

	stmt, err := Database.Prepare("INSERT INTO TransferLog (time, senderRollno, receiverRollno, coins, remarks) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(data.Time, data.SenderRollno, data.ReceiverRollno, data.Coins, data.Remarks)
	return nil
}

func AddRewardDetails(data model.RewardDetails) error {
	stmt, err := Database.Prepare("INSERT INTO RewardLog (time, receiverRollno, coins, remarks) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}

	stmt.Exec(data.Time, data.ReceiverRollno, data.Coins, data.Remarks)
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

	// checking whether the awardee is an admin or a council member
	if user.IsAdmin == 1 || user.CanEarn == 0 {
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

	_, err = Database.Exec("UPDATE User SET noOfEvents = noOfEvents + 1 WHERE rollno = ?", rollno)
	if err != nil {
		log.Fatal(err)
	}

	reward := model.RewardDetails{
		Time:           time.Now().String(),
		ReceiverRollno: rollno,
		Coins:          coins,
		Remarks:        "Coins rewarded successfully",
	}

	err = AddRewardDetails(reward)

	res.Result = "Transaction successful"
	return res
}

// Transfer given number of coins, if possible, from sender's account to the receiver's account
func TransferCoins(data model.TransferDetails) model.Response {

	var res model.Response

	sender, err := FindUser(data.SenderRollno)
	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", data.SenderRollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	// checking whether the sender has participated in adequate number of events
	if (sender.NoOfEvents | sender.IsAdmin | sender.IsinCoreTeam) == 0 {
		res.Result = "User has not participated in adequate number of events, thus cannot exchange coins"
		return res
	}

	receiver, err := FindUser(data.ReceiverRollno)
	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			res.Error = fmt.Sprint("No user with Rollno", data.ReceiverRollno)
			res.Result = "Transaction aborted"
			return res
		}

		res.Error = err.Error()
		return res
	}

	// checking whether the receiver has participated in adequate number of events
	if (receiver.NoOfEvents | receiver.IsAdmin | receiver.IsinCoreTeam) == 0 {
		res.Result = "User has not participated in adequate number of events, thus cannot exchange coins"
		return res
	}

	// checking whether the reciever is an admin or a council member
	if receiver.IsAdmin == 1 || receiver.CanEarn == 0 {
		res.Error = "User not allowed to receive coins"
		res.Result = "Transaction aborted"
		return res
	}

	// imposing the required tax on the transaction
	var coins int
	if sender.Batch == receiver.Batch {
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
			res.Result = "Transaction aborted"
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
			res.Result = "Transaction aborted"
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
		Time:           time.Now().String(),
		SenderRollno:   data.SenderRollno,
		ReceiverRollno: data.ReceiverRollno,
		Coins:          data.Coins,
		Remarks:        "Coins transferred successfully",
	}

	err = AddTransactionDetails(transaction)

	res.Result = "Transaction successful"
	return res

}