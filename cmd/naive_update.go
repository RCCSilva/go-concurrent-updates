package main

import (
	"database/sql"
)

type UserNaiveUpdate struct{}

func (UserNaiveUpdate *UserNaiveUpdate) naiveUpdateUserBalance(db *sql.DB, userId, delta int) error {
	row := db.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

	var balance int
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE balance SET balance = $1 WHERE user_id = $2", balance+delta, userId)

	if err != nil {
		return err
	}

	return nil
}
