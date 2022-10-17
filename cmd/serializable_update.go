package main

import (
	"context"
	"database/sql"

	"github.com/avast/retry-go/v4"
)

type SerializableUpdate struct{}

func (u *SerializableUpdate) update(db *sql.DB, userId, delta int) error {
	return retry.Do(func() error {
		tx, err := db.BeginTx(
			context.TODO(),
			&sql.TxOptions{Isolation: sql.LevelSerializable},
		)
		defer tx.Rollback()
		if err != nil {
			return err
		}

		row := tx.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

		var balance int
		err = row.Scan(&balance)

		if err != nil {
			return err
		}

		newBalance := balance + delta

		if newBalance < 0 {
			return nil
		}

		_, err = tx.Exec("UPDATE balance SET balance = $1 WHERE user_id = $2", newBalance, userId)

		if err != nil {
			return err
		}

		return tx.Commit()
	}, retry.DelayType(retry.BackOffDelay))
}
