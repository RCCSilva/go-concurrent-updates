package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func connectDatabase() (*sql.DB, error) {
	connStr := "postgres://user:password@localhost:5432/main?sslmode=disable"
	return sql.Open("postgres", connStr)
}
