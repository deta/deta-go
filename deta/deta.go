/*
Package deta is the official Deta SDK for Go

Example:
	package main

	import (
		"fmt"

		"github.com/deta/deta-go/deta"
		"github.com/deta/deta-go/service/base"
	)

	type User struct {
		Key      string `json:"key"` // json struct tag key to denote the key
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	func main() {
		d, err := deta.New(deta.WithProjectKey("project_key"))
		if err != nil {
			fmt.Println("failed to init new Deta instance:", err)
			return
		}

		db := base.New(d, "users")

		u := &User{
			Key:      "abasd",
			Username: "jimmy",
			Email:    "jimmy@deta.sh",
		}
		key, err := db.Put(u)
		if err != nil {
			fmt.Println("failed to put item:", err)
			return
		}
		fmt.Println("successfully put item with key", key)
	}

More examples on https://docs.deta.sh/docs/base/sdk
*/
package deta

import (
	"os"
	"strings"
)

// Deta is a top-level deta service instance
type Deta struct {
	ProjectKey string
}

type ConfigOption func(*Deta)

// Func returns a function that updates the ProjectKey
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
