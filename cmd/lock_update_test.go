package main

import (
	"sync"
	"testing"
)

func TestUpdatesWithLock(t *testing.T) {
	db, err := connectDatabase()

	if err != nil {
		t.Fatal(err)
	}
	userUpdate := &LockWithMutexUpdate{db: db, locker: &sync.Mutex{}}

	t.Run("updates balance given lock using Mutex", func(t *testing.T) {
		// Arrange
		truncateTable(t, db)

		userId := 999
		balance := 100
		delta := -20
		expected := 0

		insertUser(t, db, userId, balance)

		// Act
		c := make(chan any)
		doAsync(c, 10, func() {
			err = userUpdate.update(userId, delta)
			verifyError(t, err)
		})
		awaitChannel(c, 10)

		// Assert
		assertUserBalance(t, db, userId, expected)
	})
}
