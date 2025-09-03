package error

import (
	"encoding/json"
	"fmt"
)

// Category represents the category of an error
type Category string

const (
	CategoryNetwork    Category = "network"
	CategoryValidation Category = "validation"
	CategorySecurity   Category = "security"
	CategoryInternal   Category = "internal"
	CategoryConfig     Category = "config"
)

// Error represents a structured error with context and wrapping support
type Error struct {
	Code     int                 `json:"code"`
	Message  string              `json:"message"`
	Category Category            `json:"category"`
	Cause    error               `json:"cause,omitempty"`
	Context  map[string]any      `json:"context,omitempty"`
	Stack    string              `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.Category, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Category, e.Code, e.Message)
}

// Unwrap implements the unwrap interface for error wrapping
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *Error) WithContext(key string, value any) *Error {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// WithStack adds stack trace information to the error
func (e *Error) WithStack(stack string) *Error {
	e.Stack = stack
	return e
}

// MarshalJSON customizes JSON marshaling for better error representation
func (e *Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		Error string `json:"error"`
		*Alias
	}{
		Error: e.Error(),
		Alias: (*Alias)(e),
	})
}

// New creates a new error with the specified code, category, and message
func New(code int, category Category, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Category: category,
		Context:  make(map[string]any),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code int, category Category, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Category: category,
		Cause:    err,
		Context:  make(map[string]any),
	}
}

// Is checks if the error matches the target error type
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Code == t.Code && e.Category == t.Category
	}
	return false
}

// Common error codes
const (
	// Validation errors
	ErrCodeBadRequest = 40001

	// Authentication/Authorization errors
	ErrCodeUnauthorized = 40101
	ErrCodeForbidden    = 40301

	// Service availability errors
	ErrCodeServiceUnavailable = 50301
	ErrCodeTimeout            = 50401

	// Internal errors
	ErrCodeInternal           = 50001
	ErrCodeEncryptionFailure  = 50002
	ErrCodeDecryptionFailure  = 50003
	ErrCodeSignatureMismatch  = 50004
	ErrCodeCompressionFailure = 50005
)