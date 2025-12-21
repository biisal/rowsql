package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDbPool(ctx context.Context, dbString string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dbString)
}
