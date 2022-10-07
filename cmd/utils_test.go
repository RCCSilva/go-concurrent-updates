package main

import (
	"database/sql"
	"testing"
)

func doAsync(channel chan any, times int, function func()) {
	for i := 0; i < 10; i++ {
		go func() {
			function()
			channel <- nil
		}()
	}
}

func awaitChannel(channel chan any, times int) {
	for i := 0; i < times; i++ {
		<-channel
	}
}

func verifyError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func truncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE balance")

	if err != nil {
		t.Fatal("failed to truncate table")
	}
}

func insertUser(t *testing.T, db *sql.DB, userId, balance int) {
	_, err := db.Exec("INSERT INTO balance (user_id, balance) VALUES ($1, $2)", userId, balance)

	if err != nil {
		t.Fatalf("failed to insert user (id: %d, balance: %d)", userId, balance)
	}
}

func assertUserBalance(t *testing.T, db *sql.DB, userId, expectedBalance int) {
	row := db.QueryRow("SELECT balance FROM balance WHERE user_id = $1", userId)

	var balance int
	err := row.Scan(&balance)
	verifyError(t, err)

	if balance != expectedBalance {
		t.Fatalf("expected %d, got %d", expectedBalance, balance)
	}
}
