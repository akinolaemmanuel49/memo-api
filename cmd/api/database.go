package main

import (
	"context"
	"database/sql"

	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"

	_ "github.com/lib/pq"
)

// openDB returns a new postgres connection pool.
func openDB(config internal.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Db.DSN)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
