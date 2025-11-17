package errors

import (
	"errors"
	"testing"
)

func TestDBError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
		wantMsg string
	}{
		{
			name:    "with error",
			message: "connection failed",
			err:     errors.New("network error"),
			wantMsg: "database error: connection failed: network error",
		},
		{
			name:    "without error",
			message: "connection failed",
			err:     nil,
			wantMsg: "database error: connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbErr := NewDBError(tt.message, tt.err)
			if dbErr.Error() != tt.wantMsg {
				t.Errorf("DBError.Error() = %q, want %q", dbErr.Error(), tt.wantMsg)
			}
			if dbErr.HTTPStatus() != 500 {
				t.Errorf("DBError.HTTPStatus() = %d, want 500", dbErr.HTTPStatus())
			}
			if dbErr.Type() != ErrorTypeDB {
				t.Errorf("DBError.Type() = %q, want %q", dbErr.Type(), ErrorTypeDB)
			}
			if dbErr.Message != tt.message {
				t.Errorf("DBError.Message = %q, want %q", dbErr.Message, tt.message)
			}
			if dbErr.Err != tt.err {
				t.Errorf("DBError.Err = %v, want %v", dbErr.Err, tt.err)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	message := "invalid input"
	validationErr := NewValidationError(message)

	if validationErr.Error() != "validation error: "+message {
		t.Errorf("ValidationError.Error() = %q, want %q", validationErr.Error(), "validation error: "+message)
	}
	if validationErr.HTTPStatus() != 400 {
		t.Errorf("ValidationError.HTTPStatus() = %d, want 400", validationErr.HTTPStatus())
	}
	if validationErr.Type() != ErrorTypeValidation {
		t.Errorf("ValidationError.Type() = %q, want %q", validationErr.Type(), ErrorTypeValidation)
	}
	if validationErr.Message != message {
		t.Errorf("ValidationError.Message = %q, want %q", validationErr.Message, message)
	}
}

func TestNotFoundError(t *testing.T) {
	message := "perfume not found"
	notFoundErr := NewNotFoundError(message)

	if notFoundErr.Error() != "not found: "+message {
		t.Errorf("NotFoundError.Error() = %q, want %q", notFoundErr.Error(), "not found: "+message)
	}
	if notFoundErr.HTTPStatus() != 404 {
		t.Errorf("NotFoundError.HTTPStatus() = %d, want 404", notFoundErr.HTTPStatus())
	}
	if notFoundErr.Type() != ErrorTypeNotFound {
		t.Errorf("NotFoundError.Type() = %q, want %q", notFoundErr.Type(), ErrorTypeNotFound)
	}
	if notFoundErr.Message != message {
		t.Errorf("NotFoundError.Message = %q, want %q", notFoundErr.Message, message)
	}
}

func TestAuthError(t *testing.T) {
	message := "invalid token"
	authErr := NewAuthError(message)

	if authErr.Error() != "authentication error: "+message {
		t.Errorf("AuthError.Error() = %q, want %q", authErr.Error(), "authentication error: "+message)
	}
	if authErr.HTTPStatus() != 403 {
		t.Errorf("AuthError.HTTPStatus() = %d, want 403", authErr.HTTPStatus())
	}
	if authErr.Type() != ErrorTypeAuth {
		t.Errorf("AuthError.Type() = %q, want %q", authErr.Type(), ErrorTypeAuth)
	}
	if authErr.Message != message {
		t.Errorf("AuthError.Message = %q, want %q", authErr.Message, message)
	}
}

func TestServiceErrorInterface(t *testing.T) {
	tests := []struct {
		name       string
		err        ServiceError
		wantType   ErrorType
		wantStatus int
	}{
		{
			name:       "DBError",
			err:        NewDBError("test", nil),
			wantType:   ErrorTypeDB,
			wantStatus: 500,
		},
		{
			name:       "ValidationError",
			err:        NewValidationError("test"),
			wantType:   ErrorTypeValidation,
			wantStatus: 400,
		},
		{
			name:       "NotFoundError",
			err:        NewNotFoundError("test"),
			wantType:   ErrorTypeNotFound,
			wantStatus: 404,
		},
		{
			name:       "AuthError",
			err:        NewAuthError("test"),
			wantType:   ErrorTypeAuth,
			wantStatus: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type() != tt.wantType {
				t.Errorf("ServiceError.Type() = %q, want %q", tt.err.Type(), tt.wantType)
			}
			if tt.err.HTTPStatus() != tt.wantStatus {
				t.Errorf("ServiceError.HTTPStatus() = %d, want %d", tt.err.HTTPStatus(), tt.wantStatus)
			}
			if tt.err.Error() == "" {
				t.Error("ServiceError.Error() should not be empty")
			}
		})
	}
}
