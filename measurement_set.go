package appoptics

import (
	"sync"

	"golang.org/x/net/context"
)

type ctxMarker struct{}

// DefaultSink is a convenience instance of MeasurementSet that can be used to centrally aggregate
// measurements for an entire process.
var (
	DefaultSink  = NewMeasurementSet()
	ctxMarkerKey = &ctxMarker{}
)

// MeasurementSet represents a map of SynchronizedCounters and SynchronizedGauges. All functions
// of MeasurementSet are safe for concurrent use.
type MeasurementSet struct {
	counters      map[string]*SynchronizedCounter
	gauges        map[string]*SynchronizedGauge
	countersMutex sync.RWMutex
	gaugesMutex   sync.RWMutex
}

// NewMeasurementSet returns a new empty MeasurementSet
func NewMeasurementSet() *MeasurementSet {
	return &MeasurementSet{
		counters: map[string]*SynchronizedCounter{},
		gauges:   map[string]*SynchronizedGauge{},
	}
}

// GetCounter returns a SynchronizedCounter assigned to the specified key, creating a new one
// if necessary.
func (s *MeasurementSet) GetCounter(key string) *SynchronizedCounter {
	s.countersMutex.RLock()
	counter, ok := s.counters[key]
	s.countersMutex.RUnlock()
	if !ok {
		s.countersMutex.Lock()
		counter, ok = s.counters[key]
		if !ok {
			counter = NewCounter()
			s.counters[key] = counter
		}
		s.countersMutex.Unlock()
	}
	return counter
}

// GetGauge returns a SynchronizedGauge assigned to the specified key, creating a new one
// if necessary.
func (s *MeasurementSet) GetGauge(key string) *SynchronizedGauge {
	s.gaugesMutex.RLock()
	gauge, ok := s.gauges[key]
	s.gaugesMutex.RUnlock()
	if !ok {
		s.gaugesMutex.Lock()
		gauge, ok = s.gauges[key]
		if !ok {
			gauge = &SynchronizedGauge{}
			s.gauges[key] = gauge
		}
		s.gaugesMutex.Unlock()
	}
	return gauge
}

// Incr is a convenience function to get the specified Counter and call Incr on it. See Counter.Incr.
func (s *MeasurementSet) Incr(key string) {
	s.GetCounter(key).Incr()
}

// Add is a convenience function to get the specified Counter and call Add on it. See Counter.Add.
func (s *MeasurementSet) Add(key string, delta int64) {
	s.GetCounter(key).Add(delta)
}

// AddInt is a convenience function to get the specified Counter and call AddInt on it. See
// Counter.AddInt.
func (s *MeasurementSet) AddInt(key string, delta int) {
	s.GetCounter(key).AddInt(delta)
}

// UpdateGaugeValue is a convenience to get the specified Gauge and call UpdateValue on it.
// See Gauge.UpdateValue.
func (s *MeasurementSet) UpdateGaugeValue(key string, val int64) {
	s.GetGauge(key).UpdateValue(val)
}

// UpdateGauge is a convenience to get the specified Gauge and call Update on it. See Gauge.Update.
func (s *MeasurementSet) UpdateGauge(key string, other Gauge) {
	s.GetGauge(key).Update(other)
}

// Merge takes a MeasurementSetReport and merges all of it Counters and Gauges into this MeasurementSet.
// This in turn calls Counter.Add for each Counter in the report, and Gauge.Update for each Gauge in
// the report. Any keys that do not exist in this MeasurementSet will be created.
func (s *MeasurementSet) Merge(report *MeasurementSetReport) {
	for key, value := range report.Counts {
		s.GetCounter(key).Add(value)
	}
	for key, gauge := range report.Gauges {
		s.GetGauge(key).Update(gauge)
	}
}

// Reset generates a MeasurementSetReport with a copy of the state of each of the non-zero Counters and
// Gauges in this MeasurementSet. Counters with a value of 0 and Gauges with a count of 0 are omitted.
// All Counters and Gauges are reset to the zero/nil state but are never removed from this
// MeasurementSet, so they can continue be used indefinitely.
func (s *MeasurementSet) Reset() *MeasurementSetReport {
	report := NewMeasurementSetReport()
	s.countersMutex.Lock()
	for key, counter := range s.counters {
		val := counter.Reset()
		if val != 0 {
			report.Counts[key] = val
		}
	}
	s.countersMutex.Unlock()
	s.gaugesMutex.Lock()
	for key, syncGauge := range s.gauges {
		gauge := syncGauge.Reset()
		if gauge.Count != 0 {
			report.Gauges[key] = gauge
		}
	}
	s.gaugesMutex.Unlock()
	return report
}

// ContextWithMeasurementSet wraps the specified context with a MeasurementSet.
// XXX TODO: add convenience methods to read that MeasurementSet and manipulate Counters/Gauges on it.
func ContextWithMeasurementSet(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, NewMeasurementSet())
}
