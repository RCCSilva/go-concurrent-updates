package main

import (
	"sync"
	"testing"
)

func TestUpdatesUserBalanceWithLockApproach(t *testing.T) {
	db, err := connectDatabase()

	if err != nil {
		t.Fatal(err)
	}

	userUpdate := &UserLockUpdate{}

	t.Run("updates balance given optimistic", func(t *testing.T) {
		// Arrange
		var l sync.Mutex
		truncateTable(t, db)

		userId := 999
		balance := 100
		delta := -20
		expected := 0

		insertUser(t, db, userId, balance)

		// Act
		c := make(chan any)
		doAsync(c, 10, func() {
			err = userUpdate.update(&l, db, userId, delta)
			verifyError(t, err)
		})
		awaitChannel(c, 10)

		// Assert
		assertUserBalance(t, db, userId, expected)
	})
}
