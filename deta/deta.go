package deta

import (
	"os"
	"strings"
)

// Deta is the top-level Deta service instance
type Deta struct {
	ProjectKey string // deta project key
}

// ConfigOption is a functional config option for Deta
type ConfigOption func(*Deta)

// WithProjectKey config option for setting the project key for Deta
func WithProjectKey(projectKey string) ConfigOption {
	return func(d *Deta) {
		d.ProjectKey = projectKey
	}
}

// New returns a pointer to a new Deta instance
func New(opts ...ConfigOption) (*Deta, error) {
	d := &Deta{
		ProjectKey: os.Getenv("DETA_PROJECT_KEY"),
	}
	for _, opt := range opts {
		opt(d)
	}
	// verify project id
	if len(strings.Split(d.ProjectKey, "_")) != 2 {
		return nil, ErrBadProjectKey
	}
	return d, nil
}
