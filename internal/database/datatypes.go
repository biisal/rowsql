package database

type NumericDataType struct {
	Type             string `json:"type"`
	HasSize          bool   `json:"hasSize"`
	HasDigit         bool   `json:"hasDigit"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
}

type StringDataType struct {
	Type      string `json:"type"`
	HasSize   bool   `json:"hasSize"`
	HasValues bool   `json:"hasValues"`
}
