package appoptics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"fmt"
)

const(
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = 0
)

var (
	// Version is the current version of this client
	Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

	regexpIllegalNameChars = regexp.MustCompile("[^A-Za-z0-9.:_-]") // from https://www.AppOptics.com/docs/api/#measurements
	ErrBadStatus           = errors.New("Received non-OK status from AppOptics POST")
)

type Client interface {
	Post(batch *MeasurementsBatch) error
}

type SimpleClient struct {
	httpClient *http.Client
	URL        string
	Token      string
}

func NewClient(url, token string) Client {
	return &SimpleClient{
		URL:   url,
		Token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 4,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

type MeasurementsBatch struct {
	Time         int64             `json:"time"`
	Period       int64             `json:"period,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	Measurements []Measurement     `json:"measurements,omitempty"`
}

type Measurement struct {
	Name       string                 `json:"name"`
	Tags       map[string]string      `json:"tags,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	Count      interface{}            `json:"count,omitempty"`
	Sum        interface{}            `json:"sum,omitempty"`
	Min        interface{}            `json:"min,omitempty"`
	Max        interface{}            `json:"max,omitempty"`
	Last       interface{}            `json:"last,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Used for the collector's internal metrics
func (c *SimpleClient) Post(batch *MeasurementsBatch) error {
	if c.Token == "" {
		return errors.New("AppOptics client not authenticated")
	}
	return c.post(batch, c.Token)
}

func (c *SimpleClient) post(batch *MeasurementsBatch, token string) error {
	json, err := json.Marshal(batch)
	if err != nil {
		log.Error("Error marshaling AppOptics measurements", "err", err)
		return err
	}

	log.Debug("POSTing measurements to AppOptics", "body", string(json))
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(json))
	if err != nil {
		log.Error("Error POSTing measurements to AppOptics", "err", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("Error reading response to AppOptics measurements request", "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error reading response body from AppOptics", "statusCode", resp.StatusCode, "err", err)
			return err
		}

		log.Error("Error POSTing measurements to AppOptics", "statusCode", resp.StatusCode, "respBody", string(body))
		return ErrBadStatus
	}

	log.Debug("Finished uploading AppOptics measurements")

	return nil
}
