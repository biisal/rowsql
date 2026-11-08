// Package apperr contains error constants and functions for the application.
package apperr

import (
	"errors"
	"strconv"
)

var (
	ErrorInvalidColumn            = errors.New("invalid column name")
	ErrorInvalidDriver            = errors.New("invalid driver provided")
	ErrorEmptyTableName           = errors.New("table name cannot be empty")
	ErrorInvalidTableName         = errors.New("invalid table name")
	ErrorInvalidPagination        = errors.New("invalid limit or offset, limit must be > 0 and offset must be >= 0")
	ErrorNoValuesProvided         = errors.New("no values provided")
	ErrorDuplicateColumn          = errors.New("duplicate column name")
	ErrorInvalidJSON              = errors.New("invalid JSON")
	ErrorNoColumnsInArgs          = errors.New("no colums in args")
	ErrorInvalidAutoIncrement     = errors.New("auto-increment can only be set on primary key columns")
	ErrorPhysicalRowIDUnsupported = errors.New("this database has no physical row id")
	ErrorRowNotFound              = errors.New("row not found, it may have been changed or deleted")
	ErrorInvalidCreateTableName   = errors.New("invalid table name! Only alphanumeric and _ are allowed")
	ErrorUnsupportedDataType      = errors.New("unsupported data type for this database")
	ErrorUnsupportedDefault       = errors.New("unsupported default value")
	ErrorVerificationFailed       = errors.New("failed to verify")
	ErrorNoConfigsFound           = errors.New("no configs found")
	ErrorInvalidJSONFile          = errors.New("invlaid json file! fix the config.json")
)

func ErrorLimitTooLarge(max int) error {
	return errors.New("limit cannot be greater than " + strconv.Itoa(max))
}
