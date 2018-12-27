package appoptics

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxMeasurementsPerBatch = 1000
)

// Reporter provides a way to persist data from a set collection of Aggregators and Counters at a regular interval
type Reporter struct {
	*ReporterOpts
	batchChan   chan *MeasurementsBatch
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

// NewReporter returns a reporter for a given MeasurementSet, providing a way to sync metric information
// to AppOptics for a collection of running metrics.
func NewReporter(optFns ...ReporterOptsFn) *Reporter {
	opts := &ReporterOpts{}
	*opts = defaultReporterOpts
	for _, optFn := range optFns {
		optFn(opts)
	}
	return &Reporter{
		ReporterOpts: opts,
		batchChan:    make(chan *MeasurementsBatch, 10),
		stopChan:     make(chan struct{}, 1),
		stoppedChan:  make(chan struct{}),
	}
}

// Start kicks off two goroutines that help batch and report metrics measurements to AppOptics.
func (r *Reporter) Start() {
	go r.postMeasurementBatches()
	go r.flushReportsForever()
}

// Close forces an immediate flush of the metrics and stops further reporting.
func (r *Reporter) Close() error {
	// Notify the flushReportsForever worker that it should exit.
	select {
	case r.stopChan <- struct{}{}:
	default:
	}
	// Wait until the flushReportsForever and postMeasurementBatches workers exit
	<-r.stoppedChan
	return nil
}

func (r *Reporter) postMeasurementBatches() {
	defer close(r.stoppedChan)
	for batch := range r.batchChan {
		tryCount := 0
		for {
			log.Debug("Uploading AppOptics measurements batch", "time", time.Unix(batch.Time, 0), "numMeasurements", len(batch.Measurements), "globalTags", r.globalTags)
			_, err := r.measurementsComm.Create(batch)
			if err == nil {
				break
			}
			tryCount++
			aborting := tryCount == r.maxPostRetries
			log.Error("Error uploading AppOptics measurements batch", "err", err, "tryCount", tryCount, "aborting", aborting)
			if aborting {
				break
			}
		}
	}
}

func (r *Reporter) flushReport(report *MeasurementSetReport, reportTime time.Time) {
	var batch *MeasurementsBatch
	resetBatch := func() {
		batch = &MeasurementsBatch{
			Time:   reportTime.Unix(),
			Period: int64(r.period / time.Second),
		}
	}
	flushBatch := func() {
		r.batchChan <- batch
	}
	addMeasurement := func(measurement Measurement) {
		batch.Measurements = append(batch.Measurements, measurement)
		// AppOptics API docs advise sending very large numbers of metrics in multiple HTTP requests; so we'll limit each
		// request to a batch of 1000 measurements.
		if len(batch.Measurements) >= maxMeasurementsPerBatch {
			flushBatch()
			resetBatch()
		}
	}
	resetBatch()
	for key, value := range report.Counts {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if value != 0 {
			m.Value = float64(value)
		}
		addMeasurement(m)
	}
	// TODO: refactor to use Aggregator methods
	for key, agg := range report.Aggregators {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if agg.Sum != 0 {
			m.Sum = agg.Sum
		}
		if agg.Count != 0 {
			m.Count = agg.Count
		}
		if agg.Min != 0 {
			m.Min = agg.Min
		}
		if agg.Max != 0 {
			m.Max = agg.Max
		}
		if agg.Last != 0 {
			m.Last = agg.Last
		}
		addMeasurement(m)
	}
	if len(batch.Measurements) > 0 {
		flushBatch()
	}
}

func (r *Reporter) flushReportsForever() {
	defer close(r.batchChan)
	shutdown := false
	for !shutdown {
		// Sleep until the beginning of the next reporting period.
		now := time.Now()
		nextInterval := now.Truncate(r.period).Add(r.period)
		select {
		case <-time.After(nextInterval.Sub(now)):
		case <-r.stopChan:
			shutdown = true
		}
		report := r.measurementSet.Reset()
		if len(report.Aggregators) > 0 || len(report.Counts) > 0 {
			r.flushReport(report, nextInterval)
		}
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
