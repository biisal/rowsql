// Package models conatains the models of database qureies and forms
package models

import "time"

type RowSet struct {
	Columns []ColValue `json:"columns"`
	Key     string     `json:"key"`
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

type ColType struct {
	DataType         string `json:"dataType"`
	HasSize          bool   `json:"hasSize"`
	HasValues        bool   `json:"hasValues,omitempty"`
	HasDigit         bool   `json:"hasDigit,omitempty"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
	HasDefault       bool   `json:"hasDefault"`
	IsUnique         bool   `json:"isUnique"`
	IsPk             bool   `json:"isPk"`
	IsNull           bool   `json:"isNull"`
	PkPosition       int    `json:"-"`
	InputType        string `json:"inputType" enum:"text,number,checkbox,textarea,json,select"`
}
type ColValue struct {
	ColumnName   string  `json:"columnName"`
	Value        any     `json:"value,omitempty"`
	DefaultValue any     `json:"defaultValue,omitempty"`
	Size         int     `json:"size,omitempty"`
	ColumnType   ColType `json:"columnType"`
}

type RunSQLQueryOutput struct {
	Columns   []string `json:"columns"`
	Rows      [][]any  `json:"rows"`
	Truncated bool     `json:"truncated"`
}
