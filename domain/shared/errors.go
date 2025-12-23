package shared

import (
	"errors"
	"fmt"
)

// ErrorCode classifies domain errors for transport-agnostic mapping.
type ErrorCode string

const (
	CodeValidation ErrorCode = "VALIDATION"
	CodeConflict   ErrorCode = "CONFLICT"
	CodeNotFound   ErrorCode = "NOT_FOUND"
	CodeInternal   ErrorCode = "INTERNAL"
)

// DomainError is a structured domain error with a stable code.
type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e == nil {
		return "domain error is nil"
	}
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewDomainError creates a new domain error without wrapping.
func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// WrapDomainError creates a new domain error with a wrapped cause.
func WrapDomainError(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{Code: code, Message: message, Err: err}
}

// CodeOf extracts the domain error code, defaulting to internal when unknown.
func CodeOf(err error) ErrorCode {
	if err == nil {
		return ""
	}
	var domainErr *DomainError
	if errors.As(err, &domainErr) && domainErr.Code != "" {
		return domainErr.Code
	}
	return CodeInternal
}
