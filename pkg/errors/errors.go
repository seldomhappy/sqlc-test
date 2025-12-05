package errors

import "fmt"

// DomainError represents a domain-level error.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NotFoundError represents a "not found" error.
var NotFoundError = NewDomainError("NOT_FOUND", "resource not found", nil)

// ValidationError represents a validation error.
func ValidationError(message string) *DomainError {
	return NewDomainError("VALIDATION_ERROR", message, nil)
}

// DatabaseError represents a database error.
func DatabaseError(err error) *DomainError {
	return NewDomainError("DATABASE_ERROR", "database operation failed", err)
}
