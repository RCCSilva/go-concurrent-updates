package main

import (
	"database/sql"
	"fmt"
)

type NaiveUpdate struct{}

func (UserNaiveUpdate *NaiveUpdate) update(db *sql.DB, userId, delta int) error {
	row := db.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

	var balance int
	err := row.Scan(&balance)

	if err != nil {
		fmt.Println("---FAILED---")
		return err
	}

	newBalance := balance + delta

	if newBalance < 0 {
		return nil
	}

	_, err = db.Exec("UPDATE balance SET balance = $1 WHERE user_id = $2", newBalance, userId)

	if err != nil {
		return err
	}

	return nil
}
