package router

import (
	"errors"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/rowkey"
	"github.com/danielgtaylor/huma/v2"
)

var badRequestErrors = []error{
	rowkey.ErrInvalidKey,
	apperr.ErrorInvalidColumn,
	apperr.ErrorDuplicateColumn,
	apperr.ErrorInvalidJSON,
	apperr.ErrorNoValuesProvided,
	apperr.ErrorNoColumnsInArgs,
	apperr.ErrorUnsupportedDataType,
	apperr.ErrorUnsupportedDefault,
	apperr.ErrorInvalidAutoIncrement,
	apperr.ErrorInvalidCreateTableName,
	apperr.ErrorEmptyTableName,
	apperr.ErrorInvalidTableName,
	apperr.ErrorVerificationFailed,
}

func toHTTPError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, apperr.ErrorRowNotFound) {
		return huma.Error409Conflict(err.Error())
	}
	for _, badRequest := range badRequestErrors {
		if errors.Is(err, badRequest) {
			return huma.Error400BadRequest(err.Error())
		}
	}
	return err
}
