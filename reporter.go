package appoptics

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	outputMeasurementsIntervalSeconds = 15
	outputMeasurementsInterval        = outputMeasurementsIntervalSeconds * time.Second
	maxLibratoRetries                 = 3
	maxMeasurementsPerBatch           = 500
)

type Reporter struct {
	measurementSet      *MeasurementSet
	measurementsService *MeasurementsService
	prefix              string

	batchChan             chan *MeasurementsBatch
	measurementSetReports chan *MeasurementSetReport

	globalTags map[string]string
}

func NewReporter(measurementSet *MeasurementSet, ms *MeasurementsService, prefix string) *Reporter {
	r := &Reporter{
		measurementSet:        measurementSet,
		measurementsService:   ms,
		prefix:                prefix,
		batchChan:             make(chan *MeasurementsBatch, 100),
		measurementSetReports: make(chan *MeasurementSetReport, 1000),
	}
	r.initGlobalTags()
	return r
}

func (r *Reporter) Start() {
	go r.postMeasurementBatches()
	go r.flushReportsForever()
}

func (r *Reporter) initGlobalTags() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "na"
	}
	r.globalTags = map[string]string{
		"hostname": hostname + os.Getenv("HOST_SUFFIX"),
	}
}

func (r *Reporter) postMeasurementBatches() {
	for batch := range r.batchChan {
		tryCount := 0
		for {
			log.Debug("Uploading librato measurements batch", "time", time.Unix(batch.Time, 0), "numMeasurements", len(batch.Measurements), "globalTags", r.globalTags)
			_, err := r.measurementsService.Create(batch)
			if err == nil {
				break
			}
			tryCount++
			aborting := tryCount == maxLibratoRetries
			log.Error("Error uploading librato measurements batch", "err", err, "tryCount", tryCount, "aborting", aborting)
			if aborting {
				break
			}
		}
	}
}

func (r *Reporter) flushReport(report *MeasurementSetReport) {
	batchTimeUnixSecs := (time.Now().Unix() / outputMeasurementsIntervalSeconds) * outputMeasurementsIntervalSeconds

	var batch *MeasurementsBatch
	resetBatch := func() {
		batch = &MeasurementsBatch{
			Time:   batchTimeUnixSecs,
			Period: outputMeasurementsIntervalSeconds,
		}
	}
	flushBatch := func() {
		r.batchChan <- batch
	}
	addMeasurement := func(measurement Measurement) {
		batch.Measurements = append(batch.Measurements, measurement)
		// Librato docs advise sending very large numbers of metrics in multiple HTTP requests; so we'll flush
		// batches of 500 measurements at a time.
		if len(batch.Measurements) >= 500 {
			flushBatch()
			resetBatch()
		}
	}
	resetBatch()
	report.Counts["num_librato_measurements"] = int64(len(report.Counts)) + int64(len(report.Gauges)) + 1
	for key, value := range report.Counts {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if value != 0 {
			m.Value = value
		}
		addMeasurement(m)
	}
	// TODO: refactor to use guage methods
	for key, gauge := range report.Gauges {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if gauge.Sum != 0 {
			m.Sum = gauge.Sum
		}
		if gauge.Count != 0 {
			m.Count = gauge.Count
		}
		if gauge.Min != 0 {
			m.Min = gauge.Min
		}
		if gauge.Max != 0 {
			m.Max = gauge.Max
		}
		if gauge.Last != 0 {
			m.Last = gauge.Last
		}
		addMeasurement(m)
	}
	if len(batch.Measurements) > 0 {
		flushBatch()
	}
}

func (r *Reporter) flushReportsForever() {
	// Sleep for a random duration between 0 and outputMeasurementsInterval in order to randomize the counters output cycle.
	time.Sleep(time.Duration(rand.Int63n(int64(outputMeasurementsInterval))))
	report := r.measurementSet.Reset()
	r.flushReport(report)
	// After the initial random sleep, start a regular interval timer. This will output measurements at a consistent time
	// modulo outputMeasurementsInterval.
	for range time.Tick(outputMeasurementsInterval) {
		report := r.measurementSet.Reset()
		r.flushReport(report)
	}
}

func (r *Reporter) mergeGlobalTags(tags map[string]string) map[string]string {
	if tags == nil {
		return r.globalTags
	}

	if r.globalTags == nil {
		return tags
	}

	for k, v := range r.globalTags {
		if _, ok := tags[v]; !ok {
			tags[k] = v
		}
	}

	return tags
}
