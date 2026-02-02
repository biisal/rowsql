// Package apperr contains error constants and functions for the application.
package apperr

import (
	"errors"
	"strconv"
)

var (
	ErrorInvalidColumn           = errors.New("invalid column name")
	ErrorInvalidDriver           = errors.New("invalid driver provided")
	ErrorEmptyTableName          = errors.New("table name cannot be empty")
	ErrorInvalidTableName        = errors.New("invalid table name")
	ErrorInvalidPagination       = errors.New("invalid limit or offset, limit must be > 0 and offset must be >= 0")
	ErrorNoValueProvided         = errors.New("no value provided")
	ErrorNoValuesProvided        = errors.New("no values provided")
	ErrorDuplicateColumn         = errors.New("duplicate column name")
	ErrorInvalidJSON             = errors.New("invalid JSON")
	ErrorInvalidPlaceHolderIndex = errors.New("invalid placeholder provided! should be grather than 0")
	ErrorNotSameRowColsSize      = errors.New("cols and rows aren't same in length")
)

func ErrorLimitTooLarge(max int) error {
	return errors.New("limit cannot be greater than " + strconv.Itoa(max))
}
