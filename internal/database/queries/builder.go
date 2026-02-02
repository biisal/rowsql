// Package queries provides a set of SQL queries based on the driver
package queries

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
)

type Builder interface {
	CheckTableExitsQuery(tableName string) (string, []any)
	ColumnsList(tableName string) (string, []any)
	ListTables() string
	ListRows(tableName, orderCol, orderBy string, limit, offset int) (string, []any, error)
	DeleteTable(tableName string) string
	GetRows(tableName string, limit, offset int) (string, []any, error)
	DeleteRow(string, []models.ListDataCol, []any, int) (string, []any, error)
	UpdateRow(tableName string, form []models.RowItem, columns []models.ListDataCol, row []any) (string, []any, error)
	InsertRow(tableName string, form []models.RowItem) (string, []any, error)
	CreateTable(tableName string, inputs []database.Input) (string, error)
}

type builder struct {
	driver   configs.Driver
	maxLimit int
	dialect  Dialect
}

func NewBuilder(driver configs.Driver, maxLimit int) (Builder, error) {
	dialect, err := GetDialect(driver)
	if err != nil {
		return nil, err
	}
	return &builder{
		driver:   driver,
		maxLimit: maxLimit,
		dialect:  dialect,
	}, nil
}

func (b *builder) CheckTableExitsQuery(tableName string) (string, []any) {
	return b.dialect.CheckTableExists(tableName)
}

func (b *builder) ColumnsList(tableName string) (string, []any) {
	return b.dialect.ListColumns(tableName)
}

func (b *builder) ListTables() string {
	return b.dialect.ListTables()
}

func (b *builder) ListRows(tableName, orderCol, orderBy string, limit, offset int) (string, []any, error) {
	if tableName == "" {
		return "", nil, apperr.ErrorEmptyTableName
	}
	if limit < 0 || offset < 0 {
		return "", nil, apperr.ErrorInvalidPagination
	}
	if limit > b.maxLimit {
		return "", nil, apperr.ErrorLimitTooLarge(b.maxLimit)
	}
	tableName = b.dialect.QuoteTableName(tableName)
	if tableName == "" {
		return "", nil, apperr.ErrorInvalidTableName
	}

	parts := []string{fmt.Sprintf("SELECT * FROM %s", tableName)}
	if orderCol != "" {
		order := "ASC"
		if strings.ToLower(orderBy) == "desc" {
			order = "DESC"
		}
		parts = append(parts, fmt.Sprintf("ORDER BY %s %s", orderCol, order))
	}

	args := []any{}
	if limit > 0 {
		ph, err := b.dialect.PlaceHolder(1)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("LIMIT %s", ph))
		args = append(args, limit)
	}
	if offset > 0 {
		placeholder, err := b.dialect.PlaceHolder(1)
		if err != nil {
			return "", nil, err
		}
		if limit > 0 {
			placeholder, err = b.dialect.PlaceHolder(2)
			if err != nil {
				return "", nil, err
			}
		}
		parts = append(parts, fmt.Sprintf("OFFSET %s", placeholder))
		args = append(args, offset)
	}

	return strings.Join(parts, " "), args, nil
}

func (b *builder) InsertRow(tableName string, form []models.RowItem) (string, []any, error) {
	if tableName == "" {
		return "", nil, apperr.ErrorEmptyTableName
	}
	columns := make([]string, 0, len(form))
	placeholders := make([]string, 0, len(form))
	args := make([]any, 0, len(form))

	paramIndex := 1

	seen := make(map[string]bool)

	for _, field := range form {
		if _, ok := seen[field.ColumnName]; ok {
			return "", nil, apperr.ErrorDuplicateColumn
		}
		seen[field.ColumnName] = true
		columns = append(columns, field.ColumnName)

		ph, err := b.dialect.PlaceHolder(paramIndex)
		if err != nil {
			return "", nil, err
		}
		placeholders = append(placeholders, ph)

		if field.Type == "json" {
			var jsonVal any
			if err := json.Unmarshal([]byte(field.Value), &jsonVal); err != nil {
				logger.Errorln(err)
				var syntaxErr *json.SyntaxError
				if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
					return "", nil, apperr.ErrorInvalidJSON
				}
				return "", nil, err

			}
			logger.Info("jsonVal: %v", jsonVal)
			args = append(args, jsonVal)
		} else {
			args = append(args, field.Value)
		}

		paramIndex++
	}

	qColumns := ""
	qPlaceholders := ""

	if len(columns) > 0 {
		qColumns = strings.Join(columns, ", ")
		qPlaceholders = strings.Join(placeholders, ", ")
	}

	var query string
	tableName = b.dialect.QuoteTableName(tableName)

	if qColumns == "" {
		query = b.dialect.InsertDefaultValues(tableName)
	} else {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, qColumns, qPlaceholders)
	}
	logger.Info("Query: %s", query)
	return query, args, nil
}

func (b *builder) GetRows(tableName string, limit, offset int) (string, []any, error) {
	if limit <= 0 || offset < 0 {
		return "", nil, apperr.ErrorInvalidPagination
	}

	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return "", nil, apperr.ErrorInvalidTableName
	}

	args := []any{limit}
	ph, err := b.dialect.PlaceHolder(1)
	if err != nil {
		return "", nil, err
	}
	parts := []string{fmt.Sprintf("SELECT * FROM %s LIMIT %s", tableName, ph)}
	if offset > 0 {
		ph, err = b.dialect.PlaceHolder(2)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("OFFSET %s", ph))
		args = append(args, offset)
	}

	return strings.Join(parts, " "), args, nil
}

func (b *builder) DeleteRow(tableName string, columns []models.ListDataCol, rows []any, argIdx int) (string, []any, error) {
	clause, args, err := b.dialect.WhereCluse(columns, rows, argIdx)
	logger.Info("Clause: %s", clause)
	if err != nil {
		return "", nil, err
	}
	query := b.dialect.DeleteRow(tableName, clause)
	if query == "" {
		return "", nil, apperr.ErrorInvalidDriver
	}
	return query, args, nil
}

func (b *builder) UpdateRow(tableName string, form []models.RowItem, columns []models.ListDataCol, row []any) (string, []any, error) {
	parts := make([]string, 0, len(form))
	index := 1
	args := make([]any, 0, len(form))
	for _, v := range form {
		ph, err := b.dialect.PlaceHolder(index)
		if err != nil {
			return "", nil, err
		}

		parts = append(parts, fmt.Sprintf("%s=%s", v.ColumnName, ph))
		args = append(args, v.Value)
		index++
	}
	updateQuery := strings.Join(parts, ",")
	whereClause, whereClauseArgs, err := b.dialect.WhereCluse(columns, row, len(args)+1)
	if err != nil {
		return "", nil, err
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, updateQuery, b.dialect.FilterOneRowClause(tableName, whereClause))
	return query, append(args, whereClauseArgs...), nil
}

func (b *builder) CreateTable(tableName string, inputs []database.Input) (string, error) {
	logger.Info("Building create table query")
	columnDefs := make([]string, 0, len(inputs))
	for _, input := range inputs {
		if input.ColName == "" {
			continue
		}
		formattedColDef, err := b.formatColumnDefinition(input)
		if err != nil {
			return "", err
		}
		columnDefs = append(columnDefs, formattedColDef)
	}
	logger.Info("Query: %s", strings.Join(columnDefs, ", "))

	validName := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if !validName.MatchString(tableName) {
		logger.Error("invalid table name! Only alphanumeric and _ are allowed table_name=%s", tableName)
		return "", errors.New("invalid table name! Only alphanumeric and _ are allowed")
	}
	parts := strings.Join(columnDefs, ", ")

	query := fmt.Sprintf("CREATE TABLE %s (%s) ;", tableName, parts)
	return query, nil
}

func (b *builder) DeleteTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}
