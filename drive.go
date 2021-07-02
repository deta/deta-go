package deta

import (
	"encoding/json"
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

type startUploadResponse struct {
	UploadID  string `json:"upload_id"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	DriveName string `json:"drive_name"`
}

// Initializes a chuncked file upload.
// 
// If successful, returns the UploadID
func (d *Drive) startUpload(name string) (string, error) {
	url := "/uploads"
	queryParams := map[string]string{"name": name}
	o, err := d.client.request(&requestInput{
		Path: url,
		QueryParams: queryParams,
		Method: "POST",
	})
	if err != nil {
		return "", err
	}

	var sr startUploadResponse 
	err = json.Unmarshal(o.Body, &sr)
	if err != nil {
		return "", err
	}

	return sr.UploadID, nil
}

// End a chuncked upload.
func (d *Drive) finishUpload(name, uploadId string) error {
	url := fmt.Sprintf("/uploads/%s", uploadId)
	queryParams := map[string]string{"name": name}
	_, err := d.client.request(&requestInput{
		Path: url,
		QueryParams: queryParams,
		Method: "PATCH",
	})
	return err
}

// Abort a chunked upload. 
func (d *Drive) abortUpload(name, uploadId string) error {
	url := fmt.Sprintf("/uploads/%s", uploadId)
	queryParams := map[string]string{"name": name}
	_, err := d.client.request(&requestInput{
		Path: url,
		QueryParams: queryParams,
		Method: "DELETE",
	})
	return err
}

// Uploads a chunked part.
func (d *Drive) uploadPart(name string, chunk []byte, uploadId string, part int, contentType string) error {
	url := fmt.Sprintf("/uploads/%s/parts", uploadId)
	queryParams := map[string]string{"name": name, "part": fmt.Sprintf("%d", part)}
	_, err := d.client.request(&requestInput{
		Path: url,
		QueryParams: queryParams,
		Method: "POST",
		RawBody: chunk,
		ContentType: contentType,
	})

	if err != nil {
		return err
	}

	return nil
}


// Represents input structure for Put 
type PutInput struct {
	// name of file
	Name string
	// io.Reader with file content
	Body io.Reader
	// content type of file 
	ContentType string
}

// Put a file in the Drive.
//
// Returns the name of file that was put in the drive.
func (d *Drive) Put(i *PutInput) (string, error) {
	if i.Name == "" {
		return i.Name, ErrEmptyName
	}

	if i.Body == nil {
		return i.Name, ErrEmptyData
	}

	// start upload
	uploadId, _ := d.startUpload(i.Name)
	contentStream := i.Body
	part := 1

	for {
		chunk := make([]byte, uploadChunkSize)
		n, err := contentStream.Read(chunk)
		chunk = chunk[:n]

		if err == io.EOF {
			err = d.finishUpload(i.Name, uploadId)
			if err != nil {
				return i.Name, err
			}
			return i.Name, nil
		}

		err = d.uploadPart(i.Name, chunk, uploadId, part, i.ContentType)
		part = part + 1
		
		if err != nil {
			err = d.abortUpload(i.Name, uploadId)
			return i.Name, err
		}
	}
}