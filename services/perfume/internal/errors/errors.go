package errors

import "fmt"

type ServiceError interface {
	error
	HTTPStatus() int
	Type() ErrorType
}

type ErrorType string

const (
	ErrorTypeDB         ErrorType = "database_error"
	ErrorTypeValidation ErrorType = "validation_error"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeAuth       ErrorType = "authentication_error"
)

type DBError struct {
	Message string
	Err     error
}

func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("database error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("database error: %s", e.Message)
}

func (e *DBError) HTTPStatus() int {
	return 500
}

func (e *DBError) Type() ErrorType {
	return ErrorTypeDB
}

func NewDBError(message string, err error) *DBError {
	return &DBError{
		Message: message,
		Err:     err,
	}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) HTTPStatus() int {
	return 400
}

func (e *ValidationError) Type() ErrorType {
	return ErrorTypeValidation
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

func (e *NotFoundError) HTTPStatus() int {
	return 404
}

func (e *NotFoundError) Type() ErrorType {
	return ErrorTypeNotFound
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e *AuthError) HTTPStatus() int {
	return 403
}

func (e *AuthError) Type() ErrorType {
	return ErrorTypeAuth
}

func NewAuthError(message string) *AuthError {
	return &AuthError{
		Message: message,
	}
}
