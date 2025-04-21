package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

func NewAPIError(status int, code, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

var (
	ErrInvalidURL        = NewAPIError(http.StatusBadRequest, "INVALID_URL", "The provided URL is invalid")
	ErrCodeInUse         = NewAPIError(http.StatusConflict, "CODE_IN_USE", "The custom short code is already in use")
	ErrShortCodeNotFound = NewAPIError(http.StatusNotFound, "NOT_FOUND", "Short code does not exist")
	ErrInternal          = NewAPIError(http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
)

// IsAPIError helps to unwrap and detect custom errors
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// ExtractAPIError tries to unwrap and return a *APIError
func ExtractAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return ErrInternal.WithDetails(err.Error())
}
