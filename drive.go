package deta


import (
	"encoding/json"
	"fmt"
	"strings"
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

// List file names from drive.
//
func (d *Drive) list(limit int, prefix, last string) error {
	url := fmt.Sprintf("/files?limit=%d", limit)
	if prefix != "" {
		url = url + fmt.Sprintf("&prefix=%s", prefix)
	}
	if last != "" {
		url = url + fmt.Sprintf("&last=%s", last)
	}
	fmt.Println(url)
	return nil
}