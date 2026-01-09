package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, driver string, maxItemsPerPage int) *Queries {
	return &Queries{
		db:              db,
		driver:          driver,
		cache:           NewRowCache(100),
		maxItemsPerPage: maxItemsPerPage,
	}
}

type Queries struct {
	db              sqlx.ExtContext
	driver          string
	cache           *RowCache
	maxItemsPerPage int
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func (q *Queries) Init(ctx context.Context) (err error) {
	return q.CreateHistoryTable(ctx)
}
