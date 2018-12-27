package appoptics

import "time"

type MultiReporter struct {
	measurementSet *MeasurementSet
	reporters      []*Reporter
	period         time.Duration
	stopChan       chan struct{}
	stoppedChan    chan struct{}
}

func NewMultiReporter(m *MeasurementSet, reporters []*Reporter, period time.Duration) *MultiReporter {
	return &MultiReporter{
		measurementSet: m,
		reporters:      reporters,
		period:         period,
		stopChan:       make(chan struct{}, 1),
		stoppedChan:    make(chan struct{}),
	}
}

func (m *MultiReporter) Start() {
	for _, r := range m.reporters {
		r.Start()
	}
	go m.flushReportsForever()
}

// Close forces an immediate flush of the metrics and stops further reporting.
func (m *MultiReporter) Close() error {
	// Notify the flushReportsForever worker that it should exit.
	select {
	case m.stopChan <- struct{}{}:
	default:
	}
	// Wait until the flushReportsForever worker returns
	<-m.stoppedChan
	for _, r := range m.reporters {
		r.Close()
	}
	return nil
}

func (m *MultiReporter) flushReport(report *MeasurementSetReport, reportTime time.Time) {
	for _, r := range m.reporters {
		r.flushReport(report, reportTime)
	}
}

func (m *MultiReporter) flushReportsForever() {
	defer close(m.stoppedChan)
	shutdown := false
	for !shutdown {
		// Sleep until the beginning of the next reporting period.
		now := time.Now()
		nextInterval := now.Truncate(m.period).Add(m.period)
		select {
		case <-time.After(nextInterval.Sub(now)):
		case <-m.stopChan:
			shutdown = true
		}
		report := m.measurementSet.Reset()
		m.flushReport(report, nextInterval)
	}
}
