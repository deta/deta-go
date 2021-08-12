package deta

import (
	"errors"
)

var (
	// ErrBadProjectKey bad project key
	ErrBadProjectKey = errors.New("bad project key")
	
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

	// ErrEmptyDetaInstance empty deta instance
	ErrEmptyDetaInstance = errors.New("empty deta instance")

	// ErrBadBaseName bad base name
	ErrBadBaseName = errors.New("bad base name")
	// ErrTooManyItems too many items
	ErrTooManyItems = errors.New("too many items")
	// ErrBadDestination bad destination
	ErrBadDestination = errors.New("bad destination")
	// ErrBadItem bad item/items
	ErrBadItem = errors.New("bad item/items")

	// ErrBadDriveName bad drive name
	ErrBadDriveName = errors.New("bad drive name")
	// ErrEmptyName empty name
	ErrEmptyName = errors.New("name is empty")
	// ErrEmptyNames empty names
	ErrEmptyNames = errors.New("names is empty")
	// ErrTooManyNames too many names
	ErrTooManyNames = errors.New("too many names")
	// ErrEmptyData no data
	ErrEmptyData = errors.New("no data provided")
)
