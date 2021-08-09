/*
Package deta is the official Deta SDK for Go

Example:
	package main

	import (
		"fmt"
		"github.com/deta/deta-go"
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

		db, err := d.NewBase("base_name")
		if err != nil{
			fmt.Println("failed to init a new Base instance:", err)
			return
		}

		u := &User{
			Key: "abasd",
			Username: "jimmy",
			Email: "jimmy@deta.sh",
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
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/deta/deta-go/service/base"
	"github.com/deta/deta-go/service/drive"
)

const (
	baseEndpoint = "https://database.deta.sh/v1"
	driveEndpoint = "https://drive.deta.sh/v1"
)

var (
	// ErrBadProjectKey bad project key
	ErrBadProjectKey = errors.New("bad project key")
	// ErrBadBaseName bad base name
	ErrBadBaseName = errors.New("bad base name")
	// ErrBadDriveName bad drive name
	ErrBadDriveName = errors.New("bad drive name")

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
func (d *Deta) NewBase(baseName string) (*base.Base, error) {
	if baseName == "" {
		return nil, fmt.Errorf("%w: base name is empty", ErrBadBaseName)
	}
	rootEndpoint := os.Getenv("DETA_BASE_ROOT_ENDPOINT")
	if rootEndpoint == "" {
		rootEndpoint = baseEndpoint
	}
	return base.NewBase(d.projectKey, baseName, rootEndpoint), nil
}

// NewDrive returns a pointer to a new Drive instance
func (d *Deta) NewDrive(driveName string) (*drive.Drive, error) {
	if driveName == "" {
		return nil, fmt.Errorf("%w: drive name is empty", ErrBadDriveName)
	}
	rootEndpoint := os.Getenv("DETA_DRIVE_ROOT_ENDPOINT")
	if rootEndpoint == "" {
		rootEndpoint = driveEndpoint
	}
	return drive.NewDrive(d.projectKey, driveName, rootEndpoint), nil
}