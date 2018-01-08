package appoptics

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

// Measurement wraps the corresponding API construct: https://docs.appoptics.com/api/#measurements
// Each Measurement represents a single timeseries value for an associated Metric. If AppOptics receives a Measurement
// with a Name field that doesn't correspond to an existing Metric, a new Metric will be created.
type Measurement struct {
	Name       string                 `json:"name"`
	Tags       map[string]string      `json:"tags,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	Time       int64                  `json:"time"`
	Count      interface{}            `json:"count,omitempty"`
	Sum        interface{}            `json:"sum,omitempty"`
	Min        interface{}            `json:"min,omitempty"`
	Max        interface{}            `json:"max,omitempty"`
	Last       interface{}            `json:"last,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// MeasurementsCommunicator defines an interface for communicating with the Measurements portion of the AppOptics API
type MeasurementsCommunicator interface {
	Create(*MeasurementsBatch) (*http.Response, error)
}

// MeasurementsService implements MeasurementsCommunicator
type MeasurementsService struct {
	client *Client
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
