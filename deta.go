/*
Package deta is the official Deta SDK for Go

Example:
	package main

	import (
		"deta"
		"fmt"
	)

	type User struct {
		Key string `json:"key"`
		Username string `json:"username"`
		Email string `json:"email"`
	}

	func main(){
		d, err := deta.New("project_key")
		if err != nil{
			fmt.Println("failed to init a new Deta instance:", err)
			return
		}

		db, err := deta.NewBase("base_name")
		if err != nil{
			fmt.Println("failed to init a new Base instance:", err)
			return
		}

		u := &User{
			Key: "abasd",
			Username: "jimmy",
			Email: "jimmy@deta.sh"
		}
		key, err := db.Put(u)
		if err != nil {
			fmt.Println("failed to put item:", err)
			return
		}
		fmt.Println("successfully put item with key", key)
	}

*/
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

// Deta is a top-level deta service instance
type Deta struct {
	projectKey string
}

// New returns a pointer to a new Deta instance
func New(projectKey string) (*Deta, error) {
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

// NewBase returns a pointer to a new Base instance
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
