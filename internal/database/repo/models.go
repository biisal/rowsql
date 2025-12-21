package repo

type ListDataCol struct {
	IsUnique   bool   `json:"isUnique"`
	Value      any    `json:"value"`
	ColumnName string `json:"columnName"`
	DataType   string `json:"dataType"`
	InputType  string `json:"inputType"`
}

type ListDataProps struct {
	TableName string `json:"tableName"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
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
