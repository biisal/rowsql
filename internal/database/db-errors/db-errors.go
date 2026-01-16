// Package dberrors contains errors related to database operations.
package dberrors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Common database errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrInvalidInput        = errors.New("invalid input provided")
	ErrConnectionFailed    = errors.New("database connection failed")
	ErrTimeout             = errors.New("database operation timed out")
	ErrConstraintViolation = errors.New("database constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrTransactionFailed   = errors.New("transaction failed")
	ErrInvalidOperation    = errors.New("invalid database operation")
)

// DBError represents a detailed database error with context
type DBError struct {
	Operation   string // e.g., "insert", "update", "delete", "select"
	Resource    string // e.g., "user", "table", "database"
	Err         error  // underlying error
	Details     string // additional details for developers
	UserMessage string // user-friendly message
	HTTPStatus  int    // HTTP status code
}

func (e *DBError) Error() string {
	if e.UserMessage != "" {
		return e.UserMessage
	}
	if e.Details != "" {
		return fmt.Sprintf("%s %s: %s - %s", e.Operation, e.Resource, e.Err.Error(), e.Details)
	}
	return fmt.Sprintf("%s %s: %s", e.Operation, e.Resource, e.Err.Error())
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError creates a new database error with context
func NewDBError(operation, resource string, err error) *DBError {
	return &DBError{
		Operation:  operation,
		Resource:   resource,
		Err:        err,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// WithDetails adds developer details to the error
func (e *DBError) WithDetails(details string) *DBError {
	e.Details = details
	return e
}

// WithUserMessage sets a user-friendly message
func (e *DBError) WithUserMessage(msg string) *DBError {
	e.UserMessage = msg
	return e
}

// WithHTTPStatus sets the HTTP status code
func (e *DBError) WithHTTPStatus(status int) *DBError {
	e.HTTPStatus = status
	return e
}

// UserFriendlyError converts a database error into an HTTP status code and user-friendly error
func UserFriendlyError(err error) (int, error) {
	if err == nil {
		return http.StatusOK, nil
	}

	// Check if it's already a DBError
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		if dbErr.UserMessage != "" {
			return dbErr.HTTPStatus, errors.New(dbErr.UserMessage)
		}
		return dbErr.HTTPStatus, dbErr
	}

	// Handle standard SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, errors.New("the requested resource was not found")
	}

	if errors.Is(err, sql.ErrConnDone) {
		return http.StatusServiceUnavailable, errors.New("database connection is closed")
	}

	if errors.Is(err, sql.ErrTxDone) {
		return http.StatusInternalServerError, errors.New("database transaction has already been committed or rolled back")
	}

	// Check error message for common patterns
	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// Connection errors
	if strings.Contains(errMsgLower, "connection refused") ||
		strings.Contains(errMsgLower, "connection failed") ||
		strings.Contains(errMsgLower, "cannot connect") {
		return http.StatusServiceUnavailable, errors.New("unable to connect to database - please try again later")
	}

	// Timeout errors
	if strings.Contains(errMsgLower, "timeout") ||
		strings.Contains(errMsgLower, "deadline exceeded") {
		return http.StatusGatewayTimeout, errors.New("database operation timed out - please try again")
	}

	// Constraint violations (PostgreSQL, MySQL, SQLite)
	if strings.Contains(errMsgLower, "unique constraint") ||
		strings.Contains(errMsgLower, "duplicate key") ||
		strings.Contains(errMsgLower, "duplicate entry") {
		return http.StatusConflict, errors.New("a record with this information already exists")
	}

	// Foreign key violations
	if strings.Contains(errMsgLower, "foreign key constraint") ||
		strings.Contains(errMsgLower, "violates foreign key") {
		return http.StatusBadRequest, errors.New("cannot complete operation - referenced record does not exist or is in use")
	}

	// Syntax errors
	if strings.Contains(errMsgLower, "syntax error") ||
		strings.Contains(errMsgLower, "near") {
		return http.StatusBadRequest, errors.New("invalid request format")
	}

	// Permission errors
	if strings.Contains(errMsgLower, "permission denied") ||
		strings.Contains(errMsgLower, "access denied") {
		return http.StatusForbidden, errors.New("insufficient permissions to access this resource")
	}

	// Not null violations
	if strings.Contains(errMsgLower, "not null constraint") ||
		strings.Contains(errMsgLower, "cannot be null") {
		return http.StatusBadRequest, errors.New("required field is missing")
	}

	// Check constraint violations
	if strings.Contains(errMsgLower, "check constraint") {
		return http.StatusBadRequest, errors.New("invalid data provided - does not meet requirements")
	}

	// Generic fallback
	return http.StatusInternalServerError, errors.New("an unexpected database error occurred - please try again later")
}

// Helper functions for common error scenarios

// NotFound creates a "not found" error for a specific resource
func NotFound(resource string) *DBError {
	return &DBError{
		Operation:   "query",
		Resource:    resource,
		Err:         ErrNotFound,
		UserMessage: fmt.Sprintf("%s not found", resource),
		HTTPStatus:  http.StatusNotFound,
	}
}

// AlreadyExists creates an "already exists" error for a specific resource
func AlreadyExists(resource string) *DBError {
	return &DBError{
		Operation:   "insert",
		Resource:    resource,
		Err:         ErrAlreadyExists,
		UserMessage: fmt.Sprintf("%s already exists", resource),
		HTTPStatus:  http.StatusConflict,
	}
}

// InvalidInput creates an "invalid input" error with details
func InvalidInput(resource, details string) *DBError {
	return &DBError{
		Operation:   "validate",
		Resource:    resource,
		Err:         ErrInvalidInput,
		Details:     details,
		UserMessage: fmt.Sprintf("invalid %s: %s", resource, details),
		HTTPStatus:  http.StatusBadRequest,
	}
}

// ConnectionError creates a connection failure error
func ConnectionError(database string, err error) *DBError {
	return &DBError{
		Operation:   "connect",
		Resource:    database,
		Err:         err,
		UserMessage: "cannot connect to database. please check your connection settings",
		HTTPStatus:  http.StatusServiceUnavailable,
	}
}

// TransactionError creates a transaction failure error
func TransactionError(operation string, err error) *DBError {
	return &DBError{
		Operation:   operation,
		Resource:    "transaction",
		Err:         err,
		UserMessage: "transaction failed. please try again",
		HTTPStatus:  http.StatusInternalServerError,
	}
}

// WrapError wraps an error with database operation context
func WrapError(operation, resource string, err error) error {
	if err == nil {
		return nil
	}

	// If it's already a DBError, just update operation/resource
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		if dbErr.Operation == "" {
			dbErr.Operation = operation
		}
		if dbErr.Resource == "" {
			dbErr.Resource = resource
		}
		return dbErr
	}

	// Create new DBError and let UserFriendlyError handle status/message
	status, friendlyErr := UserFriendlyError(err)
	return &DBError{
		Operation:   operation,
		Resource:    resource,
		Err:         err,
		UserMessage: friendlyErr.Error(),
		HTTPStatus:  status,
	}
}

// HandleError is the main entry point for error handling.
// Pass any error and get back a status code and user-friendly error message.
// It automatically detects database errors and provides better messages.
//
// Usage:
//
//	status, err := dberrors.HandleError(someError)
//	http.Error(w, err.Error(), status)
func HandleError(err error) (int, error) {
	if err == nil {
		return http.StatusOK, nil
	}

	// Check if it's already a DBError with user message
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		if dbErr.UserMessage != "" {
			return dbErr.HTTPStatus, errors.New(dbErr.UserMessage)
		}
		return dbErr.HTTPStatus, dbErr
	}

	// Check for standard SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, errors.New("no data found")
	}

	if errors.Is(err, sql.ErrConnDone) {
		return http.StatusServiceUnavailable, errors.New("database connection closed. please reconnect")
	}

	if errors.Is(err, sql.ErrTxDone) {
		return http.StatusInternalServerError, errors.New("transaction already completed. please start a new one")
	}

	// Check error message for database-related patterns
	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// Connection errors
	if strings.Contains(errMsgLower, "connection refused") ||
		strings.Contains(errMsgLower, "connection failed") ||
		strings.Contains(errMsgLower, "cannot connect") ||
		strings.Contains(errMsgLower, "dial tcp") {
		return http.StatusServiceUnavailable, errors.New("cannot connect to database. please check your connection settings")
	}

	// Timeout errors
	if strings.Contains(errMsgLower, "timeout") ||
		strings.Contains(errMsgLower, "deadline exceeded") ||
		strings.Contains(errMsgLower, "context deadline") {
		return http.StatusGatewayTimeout, errors.New("query timed out. try simplifying your query or increasing timeout")
	}

	// Unique constraint violations
	if strings.Contains(errMsgLower, "unique constraint") ||
		strings.Contains(errMsgLower, "duplicate key") ||
		strings.Contains(errMsgLower, "duplicate entry") ||
		strings.Contains(errMsgLower, "uniqueness violation") {
		return http.StatusConflict, errors.New("duplicate value. this record already exists")
	}

	// Foreign key violations
	if strings.Contains(errMsgLower, "foreign key constraint") ||
		strings.Contains(errMsgLower, "violates foreign key") ||
		strings.Contains(errMsgLower, "foreign key violation") {
		return http.StatusBadRequest, errors.New("foreign key constraint failed. referenced record doesn't exist or record is in use")
	}

	// Syntax errors
	if strings.Contains(errMsgLower, "syntax error") ||
		strings.Contains(errMsgLower, "sql syntax") {
		return http.StatusBadRequest, errors.New("SQL syntax error. please check your query")
	}

	// Permission errors
	if strings.Contains(errMsgLower, "permission denied") ||
		strings.Contains(errMsgLower, "access denied") ||
		strings.Contains(errMsgLower, "unauthorized") {
		return http.StatusForbidden, errors.New("permission denied. you don't have access to this resource")
	}

	// Not null violations
	if strings.Contains(errMsgLower, "not null constraint") ||
		strings.Contains(errMsgLower, "cannot be null") ||
		strings.Contains(errMsgLower, "null value") {
		return http.StatusBadRequest, errors.New("required field cannot be empty")
	}

	// Check constraint violations
	if strings.Contains(errMsgLower, "check constraint") {
		return http.StatusBadRequest, errors.New("value violates check constraint. please verify your data")
	}

	// Table/column doesn't exist
	if strings.Contains(errMsgLower, "no such table") ||
		strings.Contains(errMsgLower, "doesn't exist") ||
		strings.Contains(errMsgLower, "unknown column") ||
		strings.Contains(errMsgLower, "relation") && strings.Contains(errMsgLower, "does not exist") {
		return http.StatusNotFound, errors.New("that table or column doesn't exist. please check the name and try again")
	}

	// If it doesn't look like a DB error, return as-is with 500
	// This handles general application errors
	return http.StatusInternalServerError, err
}
