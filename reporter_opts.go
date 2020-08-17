package appoptics

import (
	"os"
	"time"
)

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

func WithMeasurementSet(measurementSet *MeasurementSet) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.MeasurementSet = measurementSet
	}
}

func WithMeasurementsCommunicator(communicator MeasurementsCommunicator) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.measurementsComms = append(opts.measurementsComms, communicator)
	}
}

func WithMetricNamespace(metricNamespace string) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.metricPrefix = metricNamespace
		if metricNamespace != "" {
			opts.metricPrefix += "."
		}
	}
}

func WithReportingPeriod(period time.Duration) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.period = period
	}
}

func WithGlobalTags(globalTags map[string]string) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.globalTags = globalTags
	}
}

func WithMaxPostRetries(maxPostRetries int) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.maxPostRetries = maxPostRetries
	}
}
