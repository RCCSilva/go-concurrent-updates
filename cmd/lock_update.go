package main

import (
	"database/sql"
	"sync"
)

type UserLockUpdate struct{}

func (UserNaiveUpdate *UserLockUpdate) update(l sync.Locker, db *sql.DB, userId, delta int) error {
	l.Lock()
	defer l.Unlock()

	row := db.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

	var balance int
	err := row.Scan(&balance)

	if err != nil {
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
