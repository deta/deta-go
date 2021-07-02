package deta

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	uploadChunkSize = 1024 * 1024 * 10
)

var (
	// ErrEmptyName empty name 
	ErrEmptyName    = errors.New("name is empty")
	// ErrEmptyNames empty names 
	ErrEmptyNames   = errors.New("names is empty")
	// ErrTooManyNames too many items 
	ErrTooManyNames = errors.New("more than 1000 files to delete")
	// ErrEmptyData no data
	ErrEmptyData = errors.New("no data provided")
)

// Drive is a Deta Drive service client that offers the API to make requests to Deta Drive
type Drive struct {
	// deta api client
	client *detaClient

	// auth info for authenticating requests
	auth *authInfo
}

// NewDrive returns a pointer to new Drive
func newDrive(projectKey, driveName, rootEndpoint string) *Drive {
	parts := strings.Split(projectKey, "_")
	projectID := parts[0]

	// root endpoint for the base
	rootEndpoint = fmt.Sprintf("%s/%s/%s", rootEndpoint, projectID, driveName)

	return &Drive{
		client: newDetaClient(rootEndpoint, &authInfo{
			authType:    "api-key",
			headerKey:   "X-API-Key",
			headerValue: projectKey,
		}),
	}
}

// Get a file from the Drive.
//
// Returns a io.ReadCloser for the file.
func (d *Drive) Get(name string) (io.ReadCloser, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	url := "/files/download"
	queryParams := map[string]string{"name": name}
	o, err := d.client.request(&requestInput{
		Path: url,
		QueryParams: queryParams,
		Method: "GET",
		ShouldReadBody: true,
	})
	if err != nil {
		return nil, err
	}

	return o.BodyReadCloser, nil
}
