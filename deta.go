package deta

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	baseEndpoint = "https://database.deta.sh/v1"
)

var (
	// ErrBadProjectKey bad project key
	ErrBadProjectKey = errors.New("bad project key")
	// ErrBadBaseName bad base name
	ErrBadBaseName = errors.New("bad base name")
)

// Deta xx
type Deta struct {
	projectKey string
}

// NewDeta returns a pointer to a new `Deta` instance
func NewDeta(projectKey string) (*Deta, error) {
	if projectKey == "" {
		projectKey = os.Getenv("DETA_PROJECT_KEY")
	}
	// verify project id
	if len(strings.Split(projectKey, "_")) != 2 {
		return nil, ErrBadProjectKey
	}
	return &Deta{
		projectKey: projectKey,
	}, nil
}

// NewBase returns a pointer to a new 'Base' instance
func (d *Deta) NewBase(baseName string) (*Base, error) {
	if baseName == "" {
		return nil, fmt.Errorf("%w: base name is empty", ErrBadBaseName)
	}
	rootEndpoint := os.Getenv("DETA_BASE_ROOT_ENDPOINT")
	if rootEndpoint == "" {
		rootEndpoint = baseEndpoint
	}
	return newBase(d.projectKey, baseName, rootEndpoint), nil
}
