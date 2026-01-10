package errors

import "errors"

// Domain-level sentinel errors for consistent error handling across layers
// These errors are used to decouple infrastructure errors (like GORM) from domain logic

var (
	// ErrNotFound indicates that a requested resource was not found in the repository
	// Handlers should map this to HTTP 404 Not Found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates that input validation failed
	// Handlers should map this to HTTP 400 Bad Request
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal indicates an unexpected internal error occurred
	// Handlers should map this to HTTP 500 Internal Server Error
	ErrInternal = errors.New("internal error")
)
