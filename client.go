package appoptics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"fmt"

	"time"

	log "github.com/sirupsen/logrus"
)

// Version number of this package.
const (
	MajorVersion = 0
	MinorVersion = 2
	PatchVersion = 3
)

const (
	// MeasurementPostMaxBatchSize defines the max number of Measurements to send to the API at once
	MeasurementPostMaxBatchSize = 1000
	// DefaultPersistenceErrorLimit sets the number of errors that will be allowed before persistence shuts down
	DefaultPersistenceErrorLimit = 5
	defaultBaseURL               = "https://api.appoptics.com/v1/"
	defaultMediaType             = "application/json"
	clientIdentifier             = "appoptics-api-go"
)

var (
	// Version is the current version of this httpClient

	regexpIllegalNameChars = regexp.MustCompile("[^A-Za-z0-9.:_-]") // from https://www.AppOptics.com/docs/api/#measurements
	// ErrBadStatus is returned if the AppOptics API returns a non-200 error code.
	ErrBadStatus = errors.New("Received non-OK status from AppOptics POST")
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)
}

// ServiceAccessor defines an interface for talking to  via domain-specific service constructs
type ServiceAccessor interface {
	// MeasurementsService implements an interface for dealing with  Measurements
	MeasurementsService() MeasurementsCommunicator
	// SpacesService implements an interface for dealing with  Spaces
	SpacesService() SpacesCommunicator
}

// ErrorResponse represents the response body returned when the API reports an error
type ErrorResponse struct {
	// Errors holds the error information from the API
	Errors interface{} `json:"errors"`
}

// RequestErrorMessage represents the error schema for request errors
// TODO: add API reference URLs here
type RequestErrorMessage map[string][]string

// ParamErrorMessage represents the error schema for param errors
// TODO: add API reference URLs here
type ParamErrorMessage []map[string]string

// Client implements ServiceAccessor
type Client struct {
	// baseURL is the base endpoint of the remote  service
	baseURL *url.URL
	// httpClient is the http.Client singleton used for wire interaction
	httpClient *http.Client
	// token is the private part of the API credential pair
	token string
	// measurementsService embeds the httpClient and implements access to the Measurements API
	measurementsService MeasurementsCommunicator
	// spacesService embeds the httpClient and implements access to the Spaces API
	spacesService SpacesCommunicator
	// callerUserAgentFragment is placed in the User-Agent header
	callerUserAgentFragment string
}

// ClientOption provides functional option-setting behavior
type ClientOption func(*Client) error

// NewClient returns a new AppOptics API client. Optional arguments UserAgentClientOption and BaseURLClientOption can be provided.
func NewClient(token string, opts ...func(*Client) error) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		token:   token,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 4,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
	c.measurementsService = &MeasurementsService{c}
	c.spacesService = &SpacesService{c}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewRequest standardizes the request being sent
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	requestURL := c.baseURL.ResolveReference(rel)

	var buffer io.ReadWriter

	if body != nil {
		buffer = &bytes.Buffer{}
		encodeErr := json.NewEncoder(buffer).Encode(body)
		if encodeErr != nil {
			log.Println(encodeErr)
		}
	}
	req, err := http.NewRequest(method, requestURL.String(), buffer)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("token", c.token)
	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)
	req.Header.Set("User-Agent", c.completeUserAgentString())

	return req, nil
}

// UserAgentClientOption is a config function allowing setting of the User-Agent header in requests
func UserAgentClientOption(userAgentString string) ClientOption {
	return func(c *Client) error {
		c.callerUserAgentFragment = userAgentString
		return nil
	}
}

// BaseURLClientOption is a config function allowing setting of the base URL the API is on
func BaseURLClientOption(urlString string) ClientOption {
	return func(c *Client) error {
		var altURL *url.URL
		var err error
		if altURL, err = url.Parse(urlString); err != nil {
			return err
		}
		c.baseURL = altURL
		return nil
	}
}

// MeasurementsService represents the subset of the API that deals with AppOptics Measurements
func (c *Client) MeasurementsService() MeasurementsCommunicator {
	return c.measurementsService
}

// SpacesService represents the subset of the API that deals with  Measurements
func (c *Client) SpacesService() SpacesCommunicator {
	return c.spacesService
}

// Error makes ErrorResponse satisfy the error interface and can be used to serialize error responses back to the httpClient
func (e *ErrorResponse) Error() string {
	errorData, _ := json.Marshal(e)
	return string(errorData)
}

// Do performs the HTTP request on the wire, taking an optional second parameter for containing a response
func (c *Client) Do(req *http.Request, respData interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)

	// error in performing request
	if err != nil {
		return resp, err
	}

	// request response contains an error
	if err = checkError(resp); err != nil {
		return resp, err
	}

	defer resp.Body.Close()
	if respData != nil {
		if writer, ok := respData.(io.Writer); ok {
			_, err := io.Copy(writer, resp.Body)
			return resp, err
		}
		err = json.NewDecoder(resp.Body).Decode(respData)
	}

	return resp, err
}

// completeUserAgentString returns the string that will be placed in the User-Agent header.
// It ensures that any caller-set string has the client name and version appended to it.
func (c *Client) completeUserAgentString() string {
	if c.callerUserAgentFragment == "" {
		return clientVersionString()
	}
	return fmt.Sprintf("%s:%s", c.callerUserAgentFragment, clientVersionString())
}

// clientVersionString returns the canonical name-and-version string
func clientVersionString() string {
	return fmt.Sprintf("%s-v%s", clientIdentifier, Version())
}

// checkError creates an ErrorResponse from the http.Response.Body
func checkError(resp *http.Response) error {
	var errResponse ErrorResponse
	if resp.StatusCode >= 299 {
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&errResponse)
		log.Printf("error: %+v\n", errResponse)
		return &errResponse
	}
	return nil
}

// dumpResponse is a debugging function which dumps the HTTP response to stdout
func dumpResponse(resp *http.Response) {
	buf := new(bytes.Buffer)
	fmt.Printf("response status: %s\n", resp.Status)
	if resp.Body != nil {
		buf.ReadFrom(resp.Body)
		fmt.Printf("response body: %s\n\n", string(buf.Bytes()))
	}
}
