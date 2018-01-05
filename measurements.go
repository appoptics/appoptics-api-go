package appoptics

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

// Measurement wraps the corresponding API construct: https://docs.appoptics.com/api/#measurements
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

// MeasurementsBatch is a collection of Measurements persisted to the API at the same time.
// It can optionally have tags that are applied to all contained Measurements.
type MeasurementsBatch struct {
	// Measurements is the collection of timeseries entries being sent to the server
	Measurements []Measurement `json:"measurements,omitempty"`
	// Period is a slice of time measured in seconds, used in service-side aggregation
	Period int64 `json:"period,omitempty"`
	// Time is a Unix epoch timestamp used to align a group of Measurements on a time boundary
	Time int64 `json:"time"`
	// Tags are key-value identifiers that will be applied to all Measurements in the batch
	Tags *map[string]string `json:"tags,omitempty"`
}

// MeasurementsCommunicator defines an interface for communicating with the Measurements portion of the AppOptics API
type MeasurementsCommunicator interface {
	Create(*MeasurementsBatch) (*http.Response, error)
}

// MeasurementsService implements MeasurementsCommunicator
type MeasurementsService struct {
	client *Client
}

func NewMeasurementsBatch(m []Measurement, tags *map[string]string) *MeasurementsBatch {
	return &MeasurementsBatch{
		Time:         time.Now().UTC().Unix(),
		Measurements: m,
		Tags:         tags,
	}
}

// Create persists the given MeasurementCollection to AppOptics
func (ms *MeasurementsService) Create(batch *MeasurementsBatch) (*http.Response, error) {
	req, err := ms.client.NewRequest("POST", "measurements", batch)

	if err != nil {
		log.Println("error creating request:", err)
		return nil, err
	}
	return ms.client.Do(req, nil)
}

// dumpMeasurements is used for debugging
func dumpMeasurements(measurements interface{}) {
	ms := measurements.(MeasurementsBatch)
	for i, measurement := range ms.Measurements {
		floatValue, ok := measurement.Value.(float64)
		if !ok {
			continue
		}
		if math.IsNaN(floatValue) {
			fmt.Println("Found at index ", i)
			fmt.Printf("found in '%s'", measurement.Name)
		}
	}
}
