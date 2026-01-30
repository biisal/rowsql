// Package apperr contains error constants and functions for the application.
package apperr

import (
	"errors"
	"strconv"
)

var (
	ErrorInvalidColumn     = errors.New("invalid column name")
	ErrorInvalidDriver     = errors.New("invalid driver provided")
	ErrorEmptyTableName    = errors.New("table name cannot be empty")
	ErrorInvalidPagination = errors.New("invalid limit or offset, must be > 0")
)

func ErrorLimitTooLarge(max int) error {
	return errors.New("limit cannot be greater than " + strconv.Itoa(max))
}
