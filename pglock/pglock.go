package pglock

import (
	"context"
	"database/sql"
	"errors"
	"rccsilva/go-concurrent-updates/retry"
	"time"
)

type PostgresLock struct {
}

const startDuration = 100 * time.Millisecond
const attempts = 10
const maxDuration = 3 * time.Second

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
	runnable := func() error {
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
	}

	return retry.Retry(
		ctx,
		attempts,
		startDuration,
		maxDuration,
		runnable,
	)
}
