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

const (
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = 0

	// MeasurementPostMaxBatchSize defines the max number of Measurements to send to the API at once
	MeasurementPostMaxBatchSize = 1000
	defaultBaseURL              = "https://api.appoptics.com/v1/"
	defaultMediaType            = "application/json"
)

var (
	// Version is the current version of this httpClient
	Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

	regexpIllegalNameChars = regexp.MustCompile("[^A-Za-z0-9.:_-]") // from https://www.AppOptics.com/docs/api/#measurements
	ErrBadStatus           = errors.New("Received non-OK status from AppOptics POST")
)

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

// TODO: add API reference URLs here
// RequestErrorMessage represents the error schema for request errors
type RequestErrorMessage map[string][]string

// TODO: add API reference URLs here
// ParamErrorMessage represents the error schema for param errors
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
}

// TODO: make this take an optional URL parameter so it can be used against Librato API
func NewClient(token string) *Client {
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
			dumpMeasurements(body)
			return nil, encodeErr
		}

	}
	req, err := http.NewRequest(method, requestURL.String(), buffer)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("token", c.token)
	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)

	return req, nil
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
		} else {
			err = json.NewDecoder(resp.Body).Decode(respData)
		}
	}

	return resp, err
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

func dumpBody(body interface{}) {
	jsonData, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(jsonData))
}
