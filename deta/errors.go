package deta

import (
	"errors"
)

var (
	// ErrBadRequest bad request
	ErrBadRequest = errors.New("bad request")
	// ErrUnauthorized aunauthorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
	// ErrConflict conflict
	ErrConflict = errors.New("conflict")
	// ErrInternalServerError internal server error
	ErrInternalServerError = errors.New("internal server error")
)
