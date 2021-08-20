// Package client provides an internal Deta client
//
// This is an internal package and should not be used by SDK users.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/deta/deta-go/deta"
)

// AuthInfo for requests
type AuthInfo struct {
	AuthType    string // auth type
	HeaderKey   string // header key
	HeaderValue string // header value
}

// DetaClient talks to Deta APIs
type DetaClient struct {
	RootEndpoint string
	Client       *http.Client
	AuthInfo     *AuthInfo
}

// NewDetaClient returns a pointer to a new deta client
func NewDetaClient(rootEndpoint string, ai *AuthInfo) *DetaClient {
	// only api keys auth for now
	/*
		if i.Auth.Type != "api-key" {
			return nil, errInvalidAuthType
		}
	*/
	return &DetaClient{
		RootEndpoint: rootEndpoint,
		AuthInfo:     ai,
		Client:       &http.Client{},
	}
}

// error response
type errorResp struct {
	StatusCode int      `json:"-"`
	Errors     []string `json:"errors"`
}

// returns appropriate errors from the error response
func (c *DetaClient) errorRespToErr(e *errorResp) error {
	var errorMsg string
	if len(e.Errors) >= 1 {
		errorMsg = e.Errors[0]
	}

	switch e.StatusCode {
	case 400:
		return fmt.Errorf("%w: %s", deta.ErrBadRequest, errorMsg)
	case 401:
		// does not require wrapping
		return deta.ErrUnauthorized
	case 404:
		// does not require wrapping
		return deta.ErrNotFound
	case 409:
		return fmt.Errorf("%w: %s", deta.ErrConflict, errorMsg)
	default:
		// default internal server error for other error status codes
		// does not require wrapping
		return deta.ErrInternalServerError
	}
}

// RequestInput to Request method
type RequestInput struct {
	Path             string
	Method           string
	Headers          map[string]string
	QueryParams      map[string]string
	Body             interface{}
	RawBody          []byte
	ContentType      string
	ReturnReadCloser bool
}

// RequestOutput of Request method
type RequestOutput struct {
	Status         int
	Body           []byte
	BodyReadCloser io.ReadCloser
	Header         http.Header
	Error          *errorResp
}

// Request constructs and sends the request
func (c *DetaClient) Request(i *RequestInput) (*RequestOutput, error) {
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

	url := fmt.Sprintf("%s%s", c.RootEndpoint, i.Path)
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
	if c.AuthInfo != nil {
		// set auth value in specified header key in the request headers
		req.Header.Set(c.AuthInfo.HeaderKey, c.AuthInfo.HeaderValue)
	}

	// query params
	q := req.URL.Query()
	for k, v := range i.QueryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// send the request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// request output
	o := &RequestOutput{
		Status: res.StatusCode,
		Header: res.Header,
	}

	if i.ReturnReadCloser && res.StatusCode >= 200 && res.StatusCode <= 299 {
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
