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

type RowItem struct {
	ColumnName string `json:"columnName"`
	Value      string `json:"value"`
	Type       string `json:"type"`
}
type InsertRowRequest struct {
	TableName string    `json:"tableName"`
	Data      []RowItem `json:"data"`
}
type InsertDataProps struct {
	TableName string
	Values    []RowItem
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

type ColValues struct {
	Name          string `json:"colName"`
	Value         any    `json:"value"`
	Type          string `json:"type"`
	IsNull        bool   `json:"isNull"`
	IsPK          bool   `json:"isPk"`
	IsUnique      bool   `json:"isUnique"`
	Default       any    `json:"default"`
	Size          int    `json:"size,omitempty"`
	AutoIncrement bool   `json:"autoIncrement,omitempty"`
}

type ColMetaData struct {
	Name             string `json:"colName"`
	Type             string `json:"type"`
	HasSize          bool   `json:"hasSize"`
	HasValues        bool   `json:"hasValues,omitempty"`
	HasDigit         bool   `json:"hasDigit,omitempty"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
	HasDefault       any    `json:"hasDefault"`
	IsUnique         bool   `json:"isUnique"`
}
