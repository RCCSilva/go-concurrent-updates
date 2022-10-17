package main

import (
	"database/sql"
	"sync"
)

type LockWithMutexUpdate struct {
	locker sync.Locker
	db     *sql.DB
}

func (l *LockWithMutexUpdate) update(userId, delta int) error {
	l.locker.Lock()
	defer l.locker.Unlock()

	row := l.db.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

	var balance int
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	newBalance := balance + delta

	if newBalance < 0 {
		return nil
	}

	_, err = l.db.Exec("UPDATE balance SET balance = $1 WHERE user_id = $2", newBalance, userId)

	if err != nil {
		return err
	}

	return nil
}
