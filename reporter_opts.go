package appoptics

import "time"

type ReporterOpts struct {
	measurementSet   *MeasurementSet
	measurementsComm MeasurementsCommunicator
	prefix           string
	period           time.Duration
	globalTags       map[string]string
	maxPostRetries   int
}

var defaultReporterOpts = ReporterOpts{
	maxPostRetries: 3,
}

type ReporterOptsFn func(*ReporterOpts)

func WithMeasurementSet(measurementSet *MeasurementSet) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.measurementSet = measurementSet
	}
}

func WithMeasurementsCommunicator(communicator MeasurementsCommunicator) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.measurementsComm = communicator
	}
}

func WithMetricNamePrefix(prefix string) ReporterOptsFn {
	return func(opts *ReporterOpts) {
		opts.prefix = prefix
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
