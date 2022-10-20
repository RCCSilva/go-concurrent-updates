package main

import (
	"context"
	"testing"
	"time"
)

func TestUpdatesWithSerializableTransactionLevel(t *testing.T) {
	db, err := connectDatabase()

	if err != nil {
		t.Fatal(err)
	}

	userUpdate := &SerializableUpdate{db: db}

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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err = userUpdate.update(ctx, userId, delta)
		verifyError(t, err)
	})
	awaitChannel(c, 10)

	// Assert
	assertUserBalance(t, db, userId, expected)
}
