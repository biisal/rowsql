// Package queries builds the SQL used to browse and edit tables for each supported database.
package queries

import (
	"encoding/json"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/rowkey"
	"github.com/biisal/rowsql/internal/database/sqlb"
)

const PhysicalRowIDAlias = "_rowsql_row_id"

type ListRowsOptions struct {
	OrderColumn          string
	Descending           bool
	Limit                int
	Offset               int
	DefaultOrderColumns  []string
	IncludePhysicalRowID bool
}

type Builder interface {
	CheckTableExists(tableName string) (string, []any)
	ListColumns(tableName string) (string, []any)
	ListTables() string
	ListRows(tableName string, options ListRowsOptions) (string, []any, error)
	CountRows(tableName string) (string, error)
	SelectRow(tableName string, key rowkey.Key) (string, []any, error)
	InsertRow(tableName string, values []models.ColValue) (string, []any, error)
	UpdateRow(tableName string, changes []models.ColValue, key rowkey.Key) (string, []any, error)
	DeleteRow(tableName string, key rowkey.Key) (string, []any, error)
	CreateTable(tableName string, columns []models.ColValue) (string, error)
	DeleteTable(tableName string) string
	QuoteName(name string) string
	PhysicalRowID() (PhysicalRowID, bool)
}

type builder struct {
	maxLimit int
	dialect  Dialect
}

func NewBuilder(driver configs.Driver, maxLimit int) (Builder, error) {
	dialect, err := GetDialect(driver)
	if err != nil {
		return nil, err
	}
	return &builder{maxLimit: maxLimit, dialect: dialect}, nil
}

func (b *builder) CheckTableExists(tableName string) (string, []any) {
	return b.dialect.CheckTableExists(tableName)
}

func (b *builder) ListColumns(tableName string) (string, []any) {
	return b.dialect.ListColumns(tableName)
}

func (b *builder) ListTables() string {
	return b.dialect.ListTables()
}

func (b *builder) QuoteName(name string) string {
	return b.dialect.QuoteIdentifier(name)
}

func (b *builder) PhysicalRowID() (PhysicalRowID, bool) {
	return b.dialect.PhysicalRowID()
}

func (b *builder) statement() *sqlb.Statement {
	return sqlb.New(b.dialect)
}

func validateName(name string, emptyError error) error {
	if strings.TrimSpace(name) == "" {
		return emptyError
	}
	if strings.ContainsRune(name, 0) {
		return apperr.ErrorInvalidTableName
	}
	return nil
}

func (b *builder) ListRows(tableName string, options ListRowsOptions) (string, []any, error) {
	if err := validateName(tableName, apperr.ErrorEmptyTableName); err != nil {
		return "", nil, err
	}
	if options.Limit < 0 || options.Offset < 0 {
		return "", nil, apperr.ErrorInvalidPagination
	}
	if options.Limit > b.maxLimit {
		return "", nil, apperr.ErrorLimitTooLarge(b.maxLimit)
	}

	statement := b.statement().Raw("SELECT ")
	if options.IncludePhysicalRowID {
		physicalRowID, supported := b.dialect.PhysicalRowID()
		if !supported {
			return "", nil, apperr.ErrorPhysicalRowIDUnsupported
		}
		statement.Raw(physicalRowID.SelectExpr + " AS ").Identifier(PhysicalRowIDAlias).Raw(", ")
	}
	statement.Raw("* FROM ").Identifier(tableName)
	b.writeOrderBy(statement, options)
	statement.LimitOffset(options.Limit, options.Offset)

	query, arguments := statement.Build()
	return query, arguments, nil
}

func (b *builder) writeOrderBy(statement *sqlb.Statement, options ListRowsOptions) {
	orderColumns := make([]string, 0, len(options.DefaultOrderColumns)+1)
	if options.OrderColumn != "" {
		orderColumns = append(orderColumns, options.OrderColumn)
	}
	for _, column := range options.DefaultOrderColumns {
		if column != options.OrderColumn {
			orderColumns = append(orderColumns, column)
		}
	}
	if len(orderColumns) == 0 {
		return
	}

	statement.Raw(" ORDER BY ")
	for index, column := range orderColumns {
		if index > 0 {
			statement.Raw(", ")
		}
		statement.Identifier(column)
		if index == 0 && options.OrderColumn != "" && options.Descending {
			statement.Raw(" DESC")
		}
	}
}

func (b *builder) CountRows(tableName string) (string, error) {
	if err := validateName(tableName, apperr.ErrorEmptyTableName); err != nil {
		return "", err
	}
	query, _ := b.statement().Raw("SELECT COUNT(*) FROM ").Identifier(tableName).Build()
	return query, nil
}

func (b *builder) SelectRow(tableName string, key rowkey.Key) (string, []any, error) {
	if err := validateName(tableName, apperr.ErrorEmptyTableName); err != nil {
		return "", nil, err
	}
	statement := b.statement().Raw("SELECT * FROM ").Identifier(tableName)
	if err := b.writeKeyCondition(statement, key); err != nil {
		return "", nil, err
	}
	statement.LimitOffset(1, 0)
	query, arguments := statement.Build()
	return query, arguments, nil
}

func (b *builder) writeKeyCondition(statement *sqlb.Statement, key rowkey.Key) error {
	if len(key.Columns) == 0 || len(key.Columns) != len(key.Values) {
		return rowkey.ErrInvalidKey
	}

	if key.Kind == rowkey.KindPhysical {
		physicalRowID, supported := b.dialect.PhysicalRowID()
		if !supported || key.Values[0] == nil {
			return rowkey.ErrInvalidKey
		}
		statement.Raw(" WHERE ").
			Identifier(physicalRowID.Column).
			Raw(" = ").
			Argument(key.Values[0]).
			Raw(physicalRowID.BindSuffix)
		return nil
	}

	statement.WhereAllEqual(key.Columns, key.Values)
	return nil
}

func (b *builder) limitWriteToOneRow(statement *sqlb.Statement, key rowkey.Key) {
	if b.dialect.LimitsWritesToOneRow() && key.Kind == rowkey.KindAllColumns {
		statement.Raw(" LIMIT 1")
	}
}

func (b *builder) InsertRow(tableName string, values []models.ColValue) (string, []any, error) {
	if err := validateName(tableName, apperr.ErrorEmptyTableName); err != nil {
		return "", nil, err
	}
	statement := b.statement().Raw("INSERT INTO ").Identifier(tableName)
	if len(values) == 0 {
		statement.Raw(" " + b.dialect.InsertDefaultValuesClause())
		query, arguments := statement.Build()
		return query, arguments, nil
	}

	if err := rejectDuplicateColumns(values); err != nil {
		return "", nil, err
	}

	statement.Raw(" (")
	for index, value := range values {
		if index > 0 {
			statement.Raw(", ")
		}
		statement.Identifier(value.ColumnName)
	}
	statement.Raw(") VALUES (")
	for index, value := range values {
		if index > 0 {
			statement.Raw(", ")
		}
		bindable, err := bindableValue(value)
		if err != nil {
			return "", nil, err
		}
		statement.Argument(bindable)
	}
	statement.Raw(")")

	query, arguments := statement.Build()
	return query, arguments, nil
}

func (b *builder) UpdateRow(tableName string, changes []models.ColValue, key rowkey.Key) (string, []any, error) {
	if err := validateName(tableName, apperr.ErrorInvalidTableName); err != nil {
		return "", nil, err
	}
	if len(changes) == 0 {
		return "", nil, apperr.ErrorNoValuesProvided
	}
	if err := rejectDuplicateColumns(changes); err != nil {
		return "", nil, err
	}

	statement := b.statement().Raw("UPDATE ").Identifier(tableName).Raw(" SET ")
	for index, change := range changes {
		if index > 0 {
			statement.Raw(", ")
		}
		bindable, err := bindableValue(change)
		if err != nil {
			return "", nil, err
		}
		statement.Identifier(change.ColumnName).Raw(" = ").Argument(bindable)
	}
	if err := b.writeKeyCondition(statement, key); err != nil {
		return "", nil, err
	}
	b.limitWriteToOneRow(statement, key)

	query, arguments := statement.Build()
	return query, arguments, nil
}

func (b *builder) DeleteRow(tableName string, key rowkey.Key) (string, []any, error) {
	if err := validateName(tableName, apperr.ErrorEmptyTableName); err != nil {
		return "", nil, err
	}
	statement := b.statement().Raw("DELETE FROM ").Identifier(tableName)
	if err := b.writeKeyCondition(statement, key); err != nil {
		return "", nil, err
	}
	b.limitWriteToOneRow(statement, key)

	query, arguments := statement.Build()
	return query, arguments, nil
}

func (b *builder) DeleteTable(tableName string) string {
	query, _ := b.statement().Raw("DROP TABLE ").Identifier(tableName).Build()
	return query
}

func rejectDuplicateColumns(values []models.ColValue) error {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if value.ColumnName == "" {
			return apperr.ErrorInvalidColumn
		}
		if seen[value.ColumnName] {
			return apperr.ErrorDuplicateColumn
		}
		seen[value.ColumnName] = true
	}
	return nil
}

func isJSONColumn(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "json", "jsonb":
		return true
	default:
		return false
	}
}

func bindableValue(column models.ColValue) (any, error) {
	if !isJSONColumn(column.ColumnType.DataType) || column.Value == nil {
		return column.Value, nil
	}
	text, isText := column.Value.(string)
	if !isText || !json.Valid([]byte(text)) {
		return nil, apperr.ErrorInvalidJSON
	}
	return text, nil
}
