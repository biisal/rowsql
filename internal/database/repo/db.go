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
	queryBuilder    queries.Builder
	columnCache     *columnMetadataCache
	maxItemsPerPage int
}

func New(db sqlx.ExtContext, driver configs.Driver, queryBuilder queries.Builder, maxItemsPerPage int) *Queries {
	return &Queries{
		db:              db,
		driver:          driver,
		queryBuilder:    queryBuilder,
		columnCache:     newColumnMetadataCache(),
		maxItemsPerPage: maxItemsPerPage,
	}
}

func (q *Queries) Init(ctx context.Context) (err error) {
	return q.CreateHistoryTable(ctx)
}
