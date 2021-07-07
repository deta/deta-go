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
	// ErrTooManyNames too many names
	ErrTooManyNames  = errors.New("too many names")
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

// Represents input structure for DELETE request 
type deleteManyRequest struct {
	Names []string `json:"names"`
}

// Represents output structure of DeleteMany 
type DeleteManyOutput struct {
	Deleted []string          `json:"deleted"`
	Failed  map[string]string `json:"failed"`
}

// DeleteMany deletes multiple files in a Drive.
//
// Deletes at most 1000 files in a single request.
// The file names should be a string slice.
// Returns a pointer to DeleteManyOutput.
func (d *Drive) DeleteMany(names []string) (*DeleteManyOutput, error) {
	if len(names) == 0 {
		return nil, ErrEmptyNames
	}

	if len(names) > 1000 {
		return nil, errors.New("more than 1000 files to delete")
	}
	o, err := d.client.request(&requestInput{
		Path:   "/files",
		Method: "DELETE",
		Body: &deleteManyRequest{
			Names: names,
		},
	})

	if err != nil {
		return nil, err
	}

	var dr DeleteManyOutput
	err = json.Unmarshal(o.Body, &dr)

	if err != nil {
		return nil, err
	}
	return &dr, nil
}

// Delete a file from a Drive.
//
// Returns name of file deleted (even if the file does not exist)
func (d *Drive) Delete(name string) (string, error) {
	if name == "" {
		return name, ErrEmptyName
	}
	payload := []string{name}
	dr, err := d.DeleteMany(payload)
	if err != nil {
		return name, err
	}

	msg, ok := dr.Failed[name]
	if ok {
		return name, fmt.Errorf("failed to delete %s: %v", name, msg)
	}

	return name, nil
}
