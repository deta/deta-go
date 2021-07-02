package deta

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// ErrBadRequest bad request
	ErrBadRequest = errors.New("bad request")
	// ErrUnauthorized aunauthorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
	// ErrConflict conflict
	ErrConflict = errors.New("conflict")
	// ErrInternalServerError internal server error
	ErrInternalServerError = errors.New("internal server error")
	// internal error
	// invalid auth type
	errInvalidAuthType = errors.New("invalid auth type")
)

// auth info for requests
type authInfo struct {
	authType    string // auth type
	headerKey   string // header key
	headerValue string // header value
}

// client that talks with deta apis
type detaClient struct {
	rootEndpoint string
	client       *http.Client
	authInfo     *authInfo
}

// returns a pointer to a new deta client
func newDetaClient(rootEndpoint string, ai *authInfo) *detaClient {
	// only api keys auth for now
	/*
		if i.Auth.Type != "api-key" {
			return nil, errInvalidAuthType
		}
	*/
	return &detaClient{
		rootEndpoint: rootEndpoint,
		authInfo:     ai,
		client:       &http.Client{},
	}
}

// error response
type errorResp struct {
	StatusCode int      `json:"-"`
	Errors     []string `json:"errors"`
}

// returns appropriate errors from the error response
func (c *detaClient) errorRespToErr(e *errorResp) error {
	var errorMsg string
	if len(e.Errors) >= 1 {
		errorMsg = e.Errors[0]
	}

	switch e.StatusCode {
	case 400:
		return fmt.Errorf("%w: %s", ErrBadRequest, errorMsg)
	case 401:
		// does not require wrapping
		return ErrUnauthorized
	case 404:
		// does not require wrapping
		return ErrNotFound
	case 409:
		return fmt.Errorf("%w: %s", ErrConflict, errorMsg)
	default:
		// default internal server error for other error status codes
		// does not require wrapping
		return ErrInternalServerError
	}
}

// input to request method
type requestInput struct {
	Path           string
	Method         string
	Headers        map[string]string
	QueryParams    map[string]string
	Body           interface{}
	RawBody        []byte
	ContentType    string
	ShouldReadBody bool
}

// output of request function
type requestOutput struct {
	Status         int
	Body           []byte
	BodyReadCloser io.ReadCloser
	Header         http.Header
	Error          *errorResp
}

func (c *detaClient) request(i *requestInput) (*requestOutput, error) {
	marshalled := []byte("")
	if i.Body != nil {
		// set default content-type to application/json
		if i.ContentType == "" {
			i.ContentType = "application/json"
		}
		var err error
		marshalled, err = json.Marshal(&i.Body)
		if err != nil {
			return nil, err
		}
	}

	if i.RawBody != nil {
		marshalled = i.RawBody
	}

	url := fmt.Sprintf("%s%s", c.rootEndpoint, i.Path)
	req, err := http.NewRequest(i.Method, url, bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, err
	}

	// headers
	if i.ContentType != "" {
		req.Header.Set("Content-type", i.ContentType)
	}
	for k, v := range i.Headers {
		req.Header.Set(k, v)
	}

	// auth
	if c.authInfo != nil {
		// set auth value in specified header key in the request headers
		req.Header.Set(c.authInfo.headerKey, c.authInfo.headerValue)
	}

	// query params
	q := req.URL.Query()
	for k, v := range i.QueryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// send the request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// request output
	o := &requestOutput{
		Status: res.StatusCode,
		Header: res.Header,
	}

	if i.ShouldReadBody && res.StatusCode >= 200 && res.StatusCode <= 299 {
		o.BodyReadCloser = res.Body
		return o, nil
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		o.Body = b
		return o, nil
	}

	// errors
	var er errorResp
	// json unmarshal json error responses
	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if err = json.Unmarshal(b, &er); err != nil {
			return nil, err
		}
	}

	er.StatusCode = res.StatusCode
	return nil, c.errorRespToErr(&er)
}