package main

import (
	"context"
	"database/sql"
)

type PostgresAdvisoryUpdate struct {
	db     *sql.DB
	pgLock *PostgresLock
}

func (p *PostgresAdvisoryUpdate) update(ctx context.Context, userId, delta int) error {
	tx, err := p.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return err
	}

	return p.pgLock.TryWithPgLock(
		ctx,
		tx,
		userId,
		func() error {
			row := tx.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

			var balance int
			err := row.Scan(&balance)

			if err != nil {
				return err
			}

			newBalance := balance + delta

			if newBalance < 0 {
				return nil
			}

			_, err = tx.Exec(
				"UPDATE balance SET balance = $1 WHERE user_id = $2",
				newBalance,
				userId,
			)

			if err != nil {
				return err
			}

			return tx.Commit()
		},
	)
}
