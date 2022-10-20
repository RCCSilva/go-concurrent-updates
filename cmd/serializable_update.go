package main

import (
	"context"
	"database/sql"
	"rccsilva/go-concurrent-updates/retry"
	"time"
)

type SerializableUpdate struct {
	db *sql.DB
}

func (s *SerializableUpdate) update(ctx context.Context, userId, delta int) error {
	runnable := func() error {
		tx, err := s.db.BeginTx(
			ctx,
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
	}

	return retry.Retry(
		ctx,
		3,
		100*time.Millisecond,
		500*time.Millisecond,
		runnable,
	)
}
