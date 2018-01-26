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

// MeasurementSet represents a map of SynchronizedCounters and SynchronizedSummarys. All functions
// of MeasurementSet are safe for concurrent use.
type MeasurementSet struct {
	counters       map[string]*SynchronizedCounter
	summaries      map[string]*SynchronizedSummary
	countersMutex  sync.RWMutex
	summariesMutex sync.RWMutex
}

// NewMeasurementSet returns a new empty MeasurementSet
func NewMeasurementSet() *MeasurementSet {
	return &MeasurementSet{
		counters:  map[string]*SynchronizedCounter{},
		summaries: map[string]*SynchronizedSummary{},
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

// GetSummary returns a SynchronizedSummary assigned to the specified key, creating a new one
// if necessary.
func (s *MeasurementSet) GetSummary(key string) *SynchronizedSummary {
	s.summariesMutex.RLock()
	summary, ok := s.summaries[key]
	s.summariesMutex.RUnlock()
	if !ok {
		s.summariesMutex.Lock()
		summary, ok = s.summaries[key]
		if !ok {
			summary = &SynchronizedSummary{}
			s.summaries[key] = summary
		}
		s.summariesMutex.Unlock()
	}
	return summary
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

// UpdateSummaryValue is a convenience to get the specified Summary and call UpdateValue on it.
// See Summary.UpdateValue.
func (s *MeasurementSet) UpdateSummaryValue(key string, val int64) {
	s.GetSummary(key).UpdateValue(val)
}

// UpdateSummary is a convenience to get the specified Summary and call Update on it. See Summary.Update.
func (s *MeasurementSet) UpdateSummary(key string, other Summary) {
	s.GetSummary(key).Update(other)
}

// Merge takes a MeasurementSetReport and merges all of it Counters and Summarys into this MeasurementSet.
// This in turn calls Counter.Add for each Counter in the report, and Summary.Update for each Summary in
// the report. Any keys that do not exist in this MeasurementSet will be created.
func (s *MeasurementSet) Merge(report *MeasurementSetReport) {
	for key, value := range report.Counts {
		s.GetCounter(key).Add(value)
	}
	for key, summary := range report.Summaries {
		s.GetSummary(key).Update(summary)
	}
}

// Reset generates a MeasurementSetReport with a copy of the state of each of the non-zero Counters and
// Summarys in this MeasurementSet. Counters with a value of 0 and Summarys with a count of 0 are omitted.
// All Counters and Summarys are reset to the zero/nil state but are never removed from this
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
	s.summariesMutex.Lock()
	for key, syncSummary := range s.summaries {
		summary := syncSummary.Reset()
		if summary.Count != 0 {
			report.Summaries[key] = summary
		}
	}
	s.summariesMutex.Unlock()
	return report
}

// ContextWithMeasurementSet wraps the specified context with a MeasurementSet.
// XXX TODO: add convenience methods to read that MeasurementSet and manipulate Counters/Summarys on it.
func ContextWithMeasurementSet(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, NewMeasurementSet())
}
