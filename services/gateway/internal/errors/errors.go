package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GatewayError struct {
	Type       string
	StatusCode int
	Message    string
	Err        error
}

func (e *GatewayError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *GatewayError) Unwrap() error {
	return e.Err
}

func (e *GatewayError) WriteHTTP(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	errorResponse := map[string]string{
		"error":   e.Type,
		"message": e.Message,
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

func NewInternalError(err error) *GatewayError {
	return &GatewayError{
		Type:       "INTERNAL_ERROR",
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
		Err:        err,
	}
}

func ErrCORSNotAllowed(origin string) *GatewayError {
	return &GatewayError{
		Type:       "CORS_NOT_ALLOWED",
		StatusCode: http.StatusForbidden,
		Message:    "CORS not allowed",
		Err:        fmt.Errorf("origin not allowed: %s", origin),
	}
}
