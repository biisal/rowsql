package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, driver string, maxItemsPerPage int) *Queries {
	return &Queries{db, driver, NewRowCache(100), maxItemsPerPage, nil}
}

type Queries struct {
	db              sqlx.ExtContext
	driver          string
	cache           *RowCache
	maxItemsPerPage int
	Tables          []ListTablesRow
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func (q *Queries) Init(ctx context.Context) (err error) {
	q.Tables, err = q.ListTables(ctx)
	if err != nil {
		return err
	}
	return q.CreateHistoryTable(ctx)
}
