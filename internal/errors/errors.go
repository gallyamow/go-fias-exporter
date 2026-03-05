package errors

import (
	"fmt"
)

var (
	ErrPathRequired     = New("path is required")
	ErrInvalidMode      = New("invalid mode")
	ErrInvalidDBType    = New("invalid database type")
	ErrInvalidBatchSize = New("batch size must be > 0")
	ErrNoFilesFound     = New("no files found")
	ErrInvalidTableName = New("failed to resolve table name")
	ErrInvalidSchema    = New("invalid schema")
	ErrEmptyAttributes  = New("empty attributes")
)

// AppError represents application level errors
type AppError struct {
	Code    string
	Message string
	Cause   error
}

func New(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func Wrap(err error, message string) *AppError {
	return &AppError{
		Message: message,
		Cause:   err,
	}
}

func WithCode(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func WithCodeAndCause(code, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
