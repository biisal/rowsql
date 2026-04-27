package database

type VarianDataType struct {
	Type             string `json:"type"`
	HasSize          bool   `json:"hasSize"`
	HasDigit         bool   `json:"hasDigit"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
	HasValues        bool   `json:"hasValues"`
	IsNumeric        bool   `json:"isNumeric"`
}
