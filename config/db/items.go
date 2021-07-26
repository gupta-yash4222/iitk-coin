package db

import (
	_ "database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gupta-yash4222/iitk-coin/model"
	_ "github.com/mattn/go-sqlite3"
)

func FindItem(itemId string) (model.Item, error) {

	var data model.Item

	err := Database.QueryRow("SELECT itemId, itemDescription, quantity, price FROM Items WHERE itemId = ?", itemId).Scan(&data.ItemId, &data.ItemDescription, &data.Quantity, &data.Price)
	if err != nil {
		return data, err
	}

	return data, nil
}

func AddItem(itemData model.Item) error {

	_, err := FindItem(itemData.ItemId)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			stmt, err := Database.Prepare("INSERT INTO Items (itemId, itemDescription, quantity, price) VALUES (?, ?, ?, ?)")
			if err != nil {
				return err
			}

			stmt.Exec(itemData.ItemId, itemData.ItemDescription, itemData.Quantity, itemData.Price)
			return nil
		}

		return err
	}

	result, err := Database.Exec("UPDATE Items SET quantity = quantity + ? WHERE itemId = ?", itemData.Quantity, itemData.ItemId)
	rowsAffected, _ := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Database Could not be updated. Try Again.")
	}

	return errors.New("Database not maintained properly")
}

func FetchItems(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := Database.Query("SELECT itemId, itemDescription, quantity, price FROM Items")
	if err != nil {
		log.Fatal(err)
		return
	}

	var data model.Item
	for rows.Next() {
		err = rows.Scan(&data.ItemId, &data.ItemDescription, &data.Quantity, &data.Price)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Fprintf(w, "Item ID.: %s, Item Description: %s, Quantity: %d, Price: %d\n", data.ItemId, data.ItemDescription, data.Quantity, data.Price)
	}
}
