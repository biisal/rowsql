// Package database provides utilities for connecting to a PostgreSQL database.
package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDBPool(ctx context.Context, dbString string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dbString)
}
