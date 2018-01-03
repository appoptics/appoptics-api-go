package appoptics

type TaggedMeasurementSet struct {
	*MeasurementSet
	tags map[string]interface{}
}

// Tags returns the tags map
func (s *TaggedMeasurementSet ) Tags() map[string] interface{}  {
	return s.tags
}

// SetTags sets the unexported value of the unexported tags struct member and returns the struct
func (s *TaggedMeasurementSet ) SetTags(tags map[string]interface{}) {
	s.tags = tags
}

// GetCounter returns a SynchronizedCounter assigned to the specified key with tags, creating a new one
// if necessary.
func (s *TaggedMeasurementSet) GetCounter(key string) *SynchronizedCounter {
	return s.MeasurementSet.GetCounter(MetricWithTags(key, s.tags))
}

// GetGauge returns a SynchronizedGauge assigned to the specified key with tags, creating a new one
// if necessary.
func (s *TaggedMeasurementSet) GetGauge(key string) *SynchronizedGauge {
	return s.MeasurementSet.GetGauge(MetricWithTags(key, s.tags))
}

// Incr is a convenience function to get the specified Counter and call Incr on it. See Counter.Incr.
func (s *TaggedMeasurementSet) Incr(key string) {
	s.GetCounter(key).Incr()
}

// Add is a convenience function to get the specified Counter and call Add on it. See Counter.Add.
func (s *TaggedMeasurementSet) Add(key string, delta int64) {
	s.GetCounter(key).Add(delta)
}

// AddInt is a convenience function to get the specified Counter and call AddInt on it. See
// Counter.AddInt.
func (s *TaggedMeasurementSet) AddInt(key string, delta int) {
	s.GetCounter(key).AddInt(delta)
}

// UpdateGaugeValue is a convenience to get the specified Gauge and call UpdateValue on it.
// See Gauge.UpdateValue.
func (s *TaggedMeasurementSet) UpdateGaugeValue(key string, val int64) {
	s.GetGauge(key).UpdateValue(val)
}

// UpdateGauge is a convenience to get the specified Gauge and call Update on it. See Gauge.Update.
func (s *TaggedMeasurementSet) UpdateGauge(key string, other Gauge) {
	s.GetGauge(key).Update(other)
}

// Merge takes a MeasurementSetReport and merges all of it Counters and Gauges into this MeasurementSet.
// This in turn calls Counter.Add for each Counter in the report, and Gauge.Update for each Gauge in
// the report. Any keys that do not exist in this MeasurementSet will be created.
func (s *TaggedMeasurementSet) Merge(report *MeasurementSetReport) {
	for key, value := range report.Counts {
		s.GetCounter(key).Add(value)
	}
	for key, gauge := range report.Gauges {
		s.GetGauge(key).Update(gauge)
	}
}
