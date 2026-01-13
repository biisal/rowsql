// Package models conatains the models of database qureies and forms
package models

import "time"

type ListDataCol struct {
	IsUnique         bool   `json:"isUnique"`
	Value            any    `json:"value"`
	ColumnName       string `json:"columnName"`
	DataType         string `json:"dataType"`
	InputType        string `json:"inputType"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
	HasDefault       bool   `json:"hasDefault"`
}

type ListDataProps struct {
	TableName string `json:"tableName"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Column    string `json:"column"`
	Order     string `json:"order"`
}

type ListDataRow []any

type QueryParts struct {
	Columns      string
	Placeholders string
	Args         []any
}
type FormValue struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type InsertDataProps struct {
	TableName string
	Values    map[string]FormValue
}

type ListTablesRow struct {
	TableSchema string `json:"tableSchema"`
	TableName   string `json:"tableName"`
}

type History struct {
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
