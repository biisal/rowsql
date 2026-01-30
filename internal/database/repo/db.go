package repo

import (
	"context"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db              sqlx.ExtContext
	driver          configs.Driver
	queryBuilder    *queries.Builder
	cache           *RowCache
	maxItemsPerPage int
}

func New(db *sqlx.DB, driver configs.Driver, queryBuilder *queries.Builder, maxItemsPerPage int) *Queries {
	return &Queries{
		db:              db,
		driver:          driver,
		queryBuilder:    queryBuilder,
		cache:           NewRowCache(100),
		maxItemsPerPage: maxItemsPerPage,
	}
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func (q *Queries) Init(ctx context.Context) (err error) {
	return q.CreateHistoryTable(ctx)
}
