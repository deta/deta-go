package deta

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyName    = errors.New("name is empty")
	ErrEmptyNames   = errors.New("names is empty")
	ErrTooManyNames = errors.New("more than 1000 names to delete")
)

type Drive struct {
	// deta api client
	client *detaClient

	// auth info for authenticating requests
	auth *authInfo
}

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

type deleteRequest struct {
	Names []string `json:names`
}

type deleteResponse struct {
	Deleted []string          `json:deleted`
	Failed  map[string]string `json:failed`
}

func (d *Drive) DeleteMany(names []string) (*deleteResponse, error) {
	if names == nil {
		return nil, ErrEmptyNames
	}

	if len(names) > 1000 {
		return nil, ErrTooManyNames
	}
	o, err := d.client.request(&requestInput{
		Path:   "/files",
		Method: "DELETE",
		Body: &deleteRequest{
			Names: names,
		},
	})

	if err != nil {
		return nil, err
	}

	var dr deleteResponse
	err = json.Unmarshal(o.Body, &dr)

	if err != nil {
		return nil, err
	}
	return &dr, nil
}

// Delete a file from the drive
//
// If the file does not exist, a nil error is returned
func (d *Drive) Delete(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	payload := make([]string, 1)
	payload = append(payload, name)
	dr, err := d.DeleteMany(payload)
	if err != nil {
		return nil
	}

	_, ok := dr.Failed[name]
	if !ok {
		return fmt.Errorf("failed to delete %s", name)
	}

	return nil
}

type listResponse struct {
	Paging *paging       `json:paging`
	Names  []interface{} `json:names`
}

func (d *Drive) List(limit int, prefix, last string) (*listResponse, error) {
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
	var lr listResponse
	err = json.Unmarshal(o.Body, &lr)

	if err != nil {
		return nil, err
	}
	return &lr, nil
}
