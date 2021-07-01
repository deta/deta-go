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

type deleteManyRequest struct {
	Names []string `json:"names"`
}

type DeleteManyOutput struct {
	Deleted []string          `json:"deleted"`
	Failed  map[string]string `json:"failed"`
}

// DeleteMany deletes multiple files in the Drive.
//
// Deletes at most 1000 files in a single request.
// The files names should be in a slice.
// Returns a response in the DeleteManyOutput interface format.  
func (d *Drive) DeleteMany(names []string) (*DeleteManyOutput, error) {
	if len(names) == 0 || names == nil {
		return nil, ErrEmptyNames
	}

	if len(names) > 1000 {
		return nil, ErrTooManyNames
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

// Delete a file from the drive
//
// If deleted (even if file with such name does not exist), returns the name
// Else returns an error
func (d *Drive) Delete(name string) (string, error) {
	if name == "" {
		return name, ErrEmptyName
	}
	payload := make([]string, 1)
	payload = append(payload, name)
	dr, err := d.DeleteMany(payload)
	if err != nil {
		return name, err
	}

	_, ok := dr.Failed[name]
	if !ok {
		return name, fmt.Errorf("failed to delete %s", name)
	}

	return name, nil
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

// Get a file from the Drive.
//
// Returns a io.ReadCloser for the file from the Drive.
func (d *Drive) Get(name string) (io.ReadCloser, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	url := fmt.Sprintf("/files/download?name=%s", name)

	o, err := d.client.request(&requestInput{
		Path: url,
		Method: "GET",
		Read: true,
	})
	if err != nil {
		return nil, err
	}

	return o.RawBody, nil
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
	url := fmt.Sprintf("/uploads?name=%s", name)
	o, err := d.client.request(&requestInput{
		Path: url,
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
func (d *Drive) finishUpload(name, uploadId string) {
	url := fmt.Sprintf("/uploads/%s?name=%s", uploadId, name)
	_, _ = d.client.request(&requestInput{
		Path: url,
		Method: "PATCH",
	})
}

// Abort a chunked upload. 
func (d *Drive) abortUpload(name, uploadId string) {
	url := fmt.Sprintf("/uploads/%s?name=%s", uploadId, name)
	_, _ = d.client.request(&requestInput{
		Path: url,
		Method: "DELETE",
	})
}

// Uploads a chunked part.
func (d *Drive) uploadPart(name string, chunk []byte, uploadId string, part int, contentType string) error {
	url := fmt.Sprintf("/uploads/%s/parts?name=%s&part=%d", uploadId, name, part)
	_, err := d.client.request(&requestInput{
		Path: url,
		Method: "POST",
		RawBody: chunk,
		ContentType: contentType,
	})

	if err != nil {
		return err
	}

	return nil
}

type PutInput struct {
	Name string
	Body io.Reader
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
			fmt.Printf("finishing")
			d.finishUpload(i.Name, uploadId)
			return i.Name, nil
		}

		err = d.uploadPart(i.Name, chunk, uploadId, part, i.ContentType)
		part = part + 1
		
		if err != nil {
			d.abortUpload(i.Name, uploadId)
			return i.Name, err
		}
	}
}