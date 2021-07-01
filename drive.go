package deta

import (
	"encoding/json"
	"errors"
	"fmt"
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

type ListOutput struct {
	Paging *paging  `json:"paging"`
	Names  []string `json:"names"`
}

// List file names from the Drive.
//
// List is paginated, returns the last name fetched, and the size if further pages are left.
// Provide the last name in the subsequent list operation to list remaining pages.
func (d *Drive) List(limit int, prefix, last string) (*ListOutput, error) {
	url := fmt.Sprintf("/files?limit=%d", limit)
	if prefix != "" {
		url = url + fmt.Sprintf("&prefix=%s", prefix)
	}
	if last != "" {
		url = url + fmt.Sprintf("&last=%s", last)
	}
	o, err := d.client.request(&requestInput{
		Path:   url,
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}
	var lr ListOutput
	err = json.Unmarshal(o.Body, &lr)

	if err != nil {
		return nil, err
	}
	return &lr, nil
}
