package database

type NumericDataType struct {
	Type     string `json:"type"`
	HasSize  bool   `json:"hasSize"`
	HasDigit bool   `json:"hasDigit"`
}

type StringDataType struct {
	Type      string `json:"type"`
	HasSize   bool   `json:"hasSize"`
	HasValues bool   `json:"hasValues"`
}

type Input struct {
	ColName  string   `json:"colName"`
	IsNull   bool     `json:"isNull"`
	IsPK     bool     `json:"isPk"`
	IsUnique bool     `json:"isUnique"`
	DataType DataType `json:"dataType"`
}

type DataType struct {
	Type      string `json:"type"`
	HasSize   bool   `json:"hasSize"`
	HasValues bool   `json:"hasValues,omitempty"`
	Size      int    `json:"size,omitempty"`
	HasDigit  bool   `json:"hasDigit,omitempty"`
}
