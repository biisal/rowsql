// Package repo contains the repository for the database
package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type History struct {
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

const historyTableName = "rowsql_history"

func IsTableNotExistError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01"
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1146
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlite3.SQLITE_ERROR &&
			strings.Contains(strings.ToLower(err.Error()), "no such table")
	}

	return false
}

func (q *Queries) InsertHistory(ctx context.Context, message string) {
	var query string

	if q.db.DriverName() == configs.DriverPostgres {
		query = fmt.Sprintf("INSERT INTO %s (message, time) VALUES ($1, NOW()) RETURNING id", historyTableName)
	} else if q.db.DriverName() == configs.DriverMySQL {
		query = fmt.Sprintf("INSERT INTO %s (message, time) VALUES (?, NOW()) RETURNING id", historyTableName)
	} else if q.db.DriverName() == configs.DriverSQLite {
		query = fmt.Sprintf("INSERT INTO %s (message, time) VALUES (?, datetime('now')) RETURNING id", historyTableName)
	}

	_, err := q.db.ExecContext(ctx, query, message)
	if err != nil {
		if IsTableNotExistError(err) {
			if err = q.CreateHistoryTable(ctx); err != nil {
				logger.Errorln(err)
				return
			}
			_, err = q.db.ExecContext(ctx, query, message)
			if err != nil {
				logger.Errorln(err)
				return
			}
		} else {
			logger.Errorln(err)
			return
		}
	}
}

func (q *Queries) DeleteHistory(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", historyTableName)
	_, err := q.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

func (q *Queries) CreateHistoryTable(ctx context.Context) error {
	var query string

	if q.db.DriverName() == configs.DriverPostgres {
		query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, message TEXT, time TIMESTAMP WITH TIME ZONE);", historyTableName)
	} else if q.db.DriverName() == configs.DriverMySQL {
		query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, message TEXT, time TIMESTAMP DEFAULT CURRENT_TIMESTAMP);", historyTableName)
	} else if q.db.DriverName() == configs.DriverSQLite {
		query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, message TEXT, time DATETIME DEFAULT CURRENT_TIMESTAMP);", historyTableName)
	}
	_, err := q.db.ExecContext(ctx, query)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

func (q *Queries) ListHistory(ctx context.Context, limit, offset int) ([]History, error) {
	var query string
	if q.db.DriverName() == configs.DriverPostgres {
		query = fmt.Sprintf("SELECT id, message, time FROM %s ORDER BY id DESC LIMIT $1 OFFSET $2", historyTableName)
	} else if q.db.DriverName() == configs.DriverMySQL {
		query = fmt.Sprintf("SELECT id, message, time FROM %s ORDER BY id DESC LIMIT ? OFFSET ?", historyTableName)
	} else if q.db.DriverName() == configs.DriverSQLite {
		query = fmt.Sprintf("SELECT id, message, time FROM %s ORDER BY id DESC LIMIT ? OFFSET ?", historyTableName)
	}
	rows, err := q.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(&i.ID, &i.Message, &i.Time); err != nil {
			logger.Error("failed to scan rows in list cols: %v", err)
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}
	return items, nil
}
