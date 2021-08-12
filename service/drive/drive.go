package drive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/internal/client"
)

const (
	uploadChunkSize = 1024 * 1024 * 10
	driveEndpoint   = "https://drive.deta.sh/v1"
)

// Drive is a Deta Drive service client that offers the API to make requests to Deta Drive
type Drive struct {
	// deta api client
	client *client.DetaClient
}

// NewDrive returns a pointer to new Drive
func New(d *deta.Deta, driveName string) (*Drive, error) {
	if d == nil {
		return nil, deta.ErrEmptyDetaInstance
	}
	if driveName == "" {
		return nil, deta.ErrBadDriveName
	}
	projectKey := d.ProjectKey
	parts := strings.Split(projectKey, "_")
	projectID := parts[0]

	rootEndpoint := os.Getenv("DETA_DRIVE_ROOT_ENDPOINT")
	if rootEndpoint == "" {
		rootEndpoint = driveEndpoint
	}
	rootEndpoint = fmt.Sprintf("%s/%s/%s", rootEndpoint, projectID, driveName)
	return &Drive{
		client: client.NewDetaClient(rootEndpoint, &client.AuthInfo{
			AuthType:    "api-key",
			HeaderKey:   "X-API-Key",
			HeaderValue: projectKey,
		}),
	}, nil
}

// Get a file from the Drive.
//
// Returns a io.ReadCloser for the file.
func (d *Drive) Get(name string) (io.ReadCloser, error) {
	if name == "" {
		return nil, deta.ErrEmptyName
	}

	url := "/files/download"
	queryParams := map[string]string{"name": name}
	o, err := d.client.Request(&client.RequestInput{
		Path:             url,
		QueryParams:      queryParams,
		Method:           "GET",
		ReturnReadCloser: true,
	})
	if err != nil {
		return nil, err
	}

	return o.BodyReadCloser, nil
}

// startUploadResponse response for startUpload operation
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
	o, err := d.client.Request(&client.RequestInput{
		Path:        url,
		QueryParams: queryParams,
		Method:      "POST",
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
	_, err := d.client.Request(&client.RequestInput{
		Path:        url,
		QueryParams: queryParams,
		Method:      "PATCH",
	})
	return err
}

// Abort a chunked upload.
func (d *Drive) abortUpload(name, uploadId string) error {
	url := fmt.Sprintf("/uploads/%s", uploadId)
	queryParams := map[string]string{"name": name}
	_, err := d.client.Request(&client.RequestInput{
		Path:        url,
		QueryParams: queryParams,
		Method:      "DELETE",
	})
	return err
}

// Uploads a chunked part.
func (d *Drive) uploadPart(name string, chunk []byte, uploadId string, part int, contentType string) error {
	url := fmt.Sprintf("/uploads/%s/parts", uploadId)
	queryParams := map[string]string{"name": name, "part": fmt.Sprintf("%d", part)}
	_, err := d.client.Request(&client.RequestInput{
		Path:        url,
		QueryParams: queryParams,
		Method:      "POST",
		RawBody:     chunk,
		ContentType: contentType,
	})

	if err != nil {
		return err
	}

	return nil
}

// PutInput input for Put operation.
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
		return "", deta.ErrEmptyName
	}

	if i.Body == nil {
		return "", deta.ErrEmptyData
	}

	// start upload
	uploadId, err := d.startUpload(i.Name)
	if err != nil {
		return "", err
	}
	contentStream := i.Body
	part := 1

	for {
		chunk := make([]byte, uploadChunkSize)
		n, err := contentStream.Read(chunk)
		chunk = chunk[:n]

		if err == io.EOF {
			err = d.finishUpload(i.Name, uploadId)
			if err != nil {
				return "", err
			}
			return i.Name, nil
		}

		err = d.uploadPart(i.Name, chunk, uploadId, part, i.ContentType)
		part = part + 1

		if err != nil {
			err = d.abortUpload(i.Name, uploadId)
			return "", err
		}
	}
}

type paging struct {
	Size int     `json:"size"`
	Last *string `json:"last"`
}

// ListOutput output for List operation.
type ListOutput struct {
	// Pagination information
	Paging *paging `json:"paging"`
	// list of file names
	Names []string `json:"names"`
}

// List file names from the Drive.
//
// List is paginated, returns the last name fetched, and the size if further pages are left.
// Provide the last name in the subsequent list operation to list remaining pages.
func (d *Drive) List(limit int, prefix, last string) (*ListOutput, error) {
	url := "/files"
	queryParams := make(map[string]string)
	queryParams["limit"] = fmt.Sprintf("%d", limit)
	if prefix != "" {
		queryParams["prefix"] = prefix
	}
	if last != "" {
		queryParams["last"] = last
	}
	o, err := d.client.Request(&client.RequestInput{
		Path:        url,
		QueryParams: queryParams,
		Method:      "GET",
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

// deleteManyRequest input for DELETE request
type deleteManyRequest struct {
	Names []string `json:"names"`
}

// DeleteManyOutput output for DeleteMany operation
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
		return nil, deta.ErrEmptyNames
	}

	if len(names) > 1000 {
		return nil, errors.New("more than 1000 files to delete")
	}
	o, err := d.client.Request(&client.RequestInput{
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
		return "", deta.ErrEmptyName
	}
	payload := []string{name}
	dr, err := d.DeleteMany(payload)
	if err != nil {
		return "", err
	}

	msg, ok := dr.Failed[name]
	if ok {
		return "", fmt.Errorf("failed to delete %s: %v", name, msg)
	}

	return name, nil
}
