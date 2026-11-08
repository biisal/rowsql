package repo

import (
	"sort"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/rowkey"
)

type rowKeyPlan struct {
	kind    rowkey.Kind
	columns []string
}

func (q *Queries) planRowKey(columns []models.ColValue) rowKeyPlan {
	primaryKeyColumns := make([]models.ColValue, 0)
	for _, column := range columns {
		if column.ColumnType.PkPosition > 0 {
			primaryKeyColumns = append(primaryKeyColumns, column)
		}
	}
	if len(primaryKeyColumns) > 0 {
		sort.Slice(primaryKeyColumns, func(i, j int) bool {
			return primaryKeyColumns[i].ColumnType.PkPosition < primaryKeyColumns[j].ColumnType.PkPosition
		})
		return rowKeyPlan{kind: rowkey.KindPrimary, columns: columnNames(primaryKeyColumns)}
	}

	if physicalRowID, supported := q.queryBuilder.PhysicalRowID(); supported {
		return rowKeyPlan{kind: rowkey.KindPhysical, columns: []string{physicalRowID.Column}}
	}
	return rowKeyPlan{kind: rowkey.KindAllColumns, columns: columnNames(columns)}
}

func (p rowKeyPlan) defaultOrderColumns() []string {
	if p.kind == rowkey.KindAllColumns {
		return nil
	}
	return p.columns
}

func (p rowKeyPlan) buildKey(valuesByColumn map[string]any, physicalRowID any) rowkey.Key {
	key := rowkey.Key{Kind: p.kind, Columns: p.columns, Values: make([]any, len(p.columns))}
	for index, column := range p.columns {
		if p.kind == rowkey.KindPhysical {
			key.Values[index] = rowkey.NormalizeValue(physicalRowID)
			continue
		}
		key.Values[index] = rowkey.NormalizeValue(valuesByColumn[column])
	}
	return key
}

func columnNames(columns []models.ColValue) []string {
	names := make([]string, len(columns))
	for index, column := range columns {
		names[index] = column.ColumnName
	}
	return names
}

func alignWithSchema(submitted []models.ColValue, schema []models.ColValue) ([]models.ColValue, error) {
	schemaByName := make(map[string]models.ColValue, len(schema))
	for _, column := range schema {
		schemaByName[column.ColumnName] = column
	}

	aligned := make([]models.ColValue, 0, len(submitted))
	for _, column := range submitted {
		known, found := schemaByName[column.ColumnName]
		if !found {
			return nil, apperr.ErrorInvalidColumn
		}
		known.Value = column.Value
		aligned = append(aligned, known)
	}
	return aligned, nil
}
