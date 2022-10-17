package appoptics

import (
	"os"
	"strings"
	"time"
)

// ReporterOpts is a container for Reporter configuration. This struct should not
// be configured directly; instead, users should pass a slice of OptsFn to NewReporter.
type ReporterOpts struct {
	MeasurementSet    *MeasurementSet
	measurementsComms []MeasurementsCommunicator
	period            time.Duration
	metricPrefix      string
	globalTags        map[string]string
	maxPostRetries    int
}

var defaultReporterOpts = ReporterOpts{
	period:         time.Minute,
	maxPostRetries: 3,
}

// ReporterOptsFn is a mutator to configure an aspect of the Opts struct. This interface
// is implemented by the return values of the various ReporterWithXX functions in this package.
type ReporterOptsFn func(*ReporterOpts)

func withDefaultGlobalTags(opts *ReporterOpts) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "na"
	}
	opts.globalTags = map[string]string{
		"hostname": hostname + os.Getenv("HOST_SUFFIX"),
	}
}

// ReporterWithMeasurementSet sets the MeasurementSet for a Reporter.
//
// If no MeasurementSet is set, a new one will be created and can be retrieved via
// `reporter.MeasurementSet`.
func ReporterWithMeasurementSet(measurementSet *MeasurementSet) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.MeasurementSet = measurementSet
	}
}

// ReporterWithMeasurementsCommunicator adds a MeasurementsCommunicator to the Reporter.
// Multiple communicators can be attached to a single Reporter, and they each receive
// identical batches of events.
func ReporterWithMeasurementsCommunicator(communicator MeasurementsCommunicator) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.measurementsComms = append(opts.measurementsComms, communicator)
	}
}

// ReporterWithMetricNamespace sets the metric namespace for metrics emitted by the Reporter.
func ReporterWithMetricNamespace(metricNamespace string) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.metricPrefix = metricNamespace
		if opts.metricPrefix != "" && !strings.HasSuffix(opts.metricPrefix, ".") {
			opts.metricPrefix += "."
		}
	}
}

// ReporterWithReportingPeriod sets the interval at which the Reporter emits metrics.
func ReporterWithReportingPeriod(period time.Duration) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.period = period
	}
}

// ReporterWithGlobalTags sets a global set of tags to apply to all metrics emitted
// by the reporter. These tags may be overridden by individual metrics that contain
// a tag with a matching key.
func ReporterWithGlobalTags(globalTags map[string]string) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.globalTags = globalTags
	}
}

// ReporterWithMaxPostRetries sets the maximum number of retries to attempt for
// each batch of metrics.
func ReporterWithMaxPostRetries(maxPostRetries int) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.maxPostRetries = maxPostRetries
	}
}
