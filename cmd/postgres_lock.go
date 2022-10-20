package main

import (
	"context"
	"database/sql"
	"errors"
	"rccsilva/go-concurrent-updates/cmd/retry"
	"time"
)

type PostgresLock struct {
	db *sql.DB
}

func (p *PostgresLock) TryWithPgLock(
	ctx context.Context,
	tx *sql.Tx,
	key int,
	runnable func() error,
) error {
	err := p.acquireLock(ctx, tx, key)
	if err != nil {
		return err
	}
	return runnable()
}

func (p *PostgresLock) acquireLock(ctx context.Context, tx *sql.Tx, key int) error {
	return retry.Retry(
		ctx,
		10,
		100*time.Millisecond,
		func() error {
			row := tx.QueryRow("select pg_try_advisory_xact_lock($1)", key)
			var result bool
			err := row.Scan(&result)

			if err != nil {
				return err
			}

			if !result {
				return errors.New("unable to acquire lock")
			}
			return nil
		},
	)
}
