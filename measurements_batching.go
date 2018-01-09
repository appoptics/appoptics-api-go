package appoptics

import (
	"fmt"
	"time"

	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/solarwinds/prometheus2appoptics/config"
)

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

// BatchPersister implements persistence to AppOptics and enforces error limits
type BatchPersister struct {
	// ms is the MeasurementsCommunicator used to talk to the AppOptics API
	ms MeasurementsCommunicator
	// errorLimit is the number of persistence errors that will be tolerated
	errorLimit int
	// prepChan is a channel of Measurements slices
	prepChan chan []Measurement
	// batchChan is used to create MeasurementsBatches for persistence to AppOptics
	batchChan chan *MeasurementsBatch
	// stopChan is used to cease persisting MeasurementsBatches to AppOptics
	stopChan chan bool
	// errorChan is used to tally errors that occur in batching/persisting
	errorChan chan error
	// maximumPushInterval is the max time (in milliseconds) to wait before pushing a batch whether its length is equal
	// to the MeasurementPostMaxBatchSize or not
	maximumPushInterval int
}

// NewBatchPersister sets up a new instance of batched persistence capabilites using the provided MeasurementsCommunicator
func NewBatchPersister(ms MeasurementsCommunicator) *BatchPersister {
	return &BatchPersister{
		ms:                  ms,
		errorLimit:          DefaultPersistenceErrorLimit,
		prepChan:            make(chan []Measurement),
		batchChan:           make(chan *MeasurementsBatch),
		stopChan:            make(chan bool),
		errorChan:           make(chan error),
		maximumPushInterval: 2000,
	}
}

func NewMeasurementsBatch(m []Measurement, tags *map[string]string) *MeasurementsBatch {
	return &MeasurementsBatch{
		Measurements: m,
		Tags:         tags,
		Time:         time.Now().UTC().Unix(),
	}
}

// MeasurementsSink gives calling code write-only access to the Measurements prep channel
func (bp *BatchPersister) MeasurementsSink() chan<- []Measurement {
	return bp.prepChan
}

// MeasurementsStopChannel gives calling code write-only access to the Measurements control channel
func (bp *BatchPersister) MeasurementsStopChannel() chan<- bool {
	return bp.stopChan
}

// MeasurementsErrorChannel gives calling code write-only access to the Measurements error channel
func (bp *BatchPersister) MeasurementsErrorChannel() chan<- error {
	return bp.errorChan
}

// MaxiumumPushIntervalMilliseconds returns the number of milliseconds the system will wait before pushing any
// accumulated Measurements to AppOptics
func (bp *BatchPersister) MaximumPushInterval() int {
	return bp.maximumPushInterval
}

// SetMaximumPushInterval sets the number of milliseconds the system will wait before pushing any accumulated
// Measurements to AppOptics
func (bp *BatchPersister) SetMaximumPushInterval(ms int) {
	bp.maximumPushInterval = ms
}

// batchMeasurements reads slices of Measurements off a channel and packages them into batches conforming to the
// limitations imposed by the API. If Measurements are arriving slowly, collected Measurements will be pushed on an
// interval defined by maximumPushIntervalMilliseconds
func (bp *BatchPersister) batchMeasurements() {
	var currentMeasurements = []Measurement{}
	pushBatch := &MeasurementsBatch{}
	ticker := time.NewTicker(time.Millisecond * time.Duration(bp.maximumPushInterval))
LOOP:
	for {
		select {
		case receivedMeasurements := <-bp.prepChan:
			currentMeasurements = append(currentMeasurements, receivedMeasurements...)
			if len(currentMeasurements) >= MeasurementPostMaxBatchSize {
				pushBatch.Measurements = currentMeasurements[:MeasurementPostMaxBatchSize]
				bp.batchChan <- pushBatch
				currentMeasurements = currentMeasurements[MeasurementPostMaxBatchSize:]
			}
		case <-ticker.C:
			if len(currentMeasurements) > 0 {
				if len(currentMeasurements) >= MeasurementPostMaxBatchSize {
					pushBatch.Measurements = currentMeasurements[:MeasurementPostMaxBatchSize]
					bp.batchChan <- pushBatch
					currentMeasurements = currentMeasurements[MeasurementPostMaxBatchSize:]
				} else {
					bp.batchChan <- pushBatch
					currentMeasurements = []Measurement{}
				}
			}
		case <-bp.stopChan:
			break LOOP
		}
	}
}

// BatchAndPersistMeasurementsForever continually packages up Measurements from the channel returned by MeasurementSink()
// and persists them to the AppOptics backend
func (bp *BatchPersister) BatchAndPersistMeasurementsForever() {
	go bp.batchMeasurements()
	go bp.persistBatches()
	go bp.managePersistenceErrors()
}

// persistBatches reads maximal slices of AppOptics.Measurement types off a channel and persists them to the remote AppOptics
// API. Errors are placed on the error channel.
func (bp *BatchPersister) persistBatches() {
	ticker := time.NewTicker(time.Millisecond * 500)
LOOP:
	for {
		select {
		case <-ticker.C:
			batch := <-bp.batchChan
			err := bp.persistBatch(batch)
			if err != nil {
				bp.errorChan <- err
			}
		case <-bp.stopChan:
			ticker.Stop()
			break LOOP
		}
	}
}

// managePersistenceErrors tracks errors on the provided channel and sends a stop signal if the ErrorLimit is reached
func (bp *BatchPersister) managePersistenceErrors() {
	var errors []error
LOOP:
	for {
		select {
		case err := <-bp.errorChan:
			errors = append(errors, err)
			if len(errors) > bp.errorLimit {
				bp.stopChan <- true
				break LOOP
			}
		}

	}
}

// persistBatch sends to the remote AppOptics endpoint unless config.SendStats() returns false, when it prints to stdout
func (bp *BatchPersister) persistBatch(batch *MeasurementsBatch) error {
	if config.SendStats() {
		log.Printf("persisting %d Measurements to AppOptics\n", len(batch.Measurements))
		resp, err := bp.ms.Create(batch)
		if resp == nil {
			fmt.Println("response is nil")
			return err
		}
		dumpResponse(resp)
	} else {
		printMeasurements(batch.Measurements)
	}
	return nil
}

// printMeasurements pretty-prints the supplied measurements to stdout
func printMeasurements(data []Measurement) {
	for _, measurement := range data {
		fmt.Printf("\nMetric name: '%s' \n", measurement.Name)
		fmt.Printf("\t value: %d \n", measurement.Value)
		fmt.Printf("\t\tTags: ")
		for k, v := range measurement.Tags {
			fmt.Printf("\n\t\t\t%s: %s", k, v)
		}
	}
}

func dumpResponse(resp *http.Response) {
	buf := new(bytes.Buffer)
	fmt.Printf("response status: %s\n", resp.Status)
	if resp.Body != nil {
		buf.ReadFrom(resp.Body)
		fmt.Printf("response body: %s\n\n", string(buf.Bytes()))
	}
}
