package appoptics

import "sync"

// A Gauge corresponds roughly to the Gauge feature in Librato, storing a count/sum/min/max/last.
// It can be either updated by passing sequential values to UpdateValue or by passing a Gauge
// to Update, e.g. gauge.Update(Gauge{Sum:100,Count:10,Min:5,Max:15})
type Gauge struct {
	Count int64
	Sum   int64
	Min   int64
	Max   int64
	Last  int64
}

// UpdateValue sets the most recently observed value for this gauge, updating sum/count/min/max/last
// accordingly.
func (g *Gauge) UpdateValue(val int64) {
	if g.Count == 0 {
		g.Min = val
		g.Max = val
	} else {
		if val < g.Min {
			g.Min = val
		}
		if val > g.Max {
			g.Max = val
		}
	}
	g.Count++
	g.Sum += val
	g.Last = val
}

// Update merges another Gauge into this Gauge, merging sum/count/min/max/last accordingly. It can
// be used to facilitate efficient input of many data points into a Gauge in one call, and it can
// also be used to merge two different Gauges (for example, workers can each maintain their own
// and periodically merge them).
func (g *Gauge) Update(other Gauge) {
	if g.Count == 0 {
		g.Count = other.Count
		g.Sum = other.Sum
		g.Min = other.Min
		g.Max = other.Max
		g.Last = other.Last
	} else {
		g.Count += other.Count
		g.Sum += other.Sum
		if other.Min < g.Min {
			g.Min = other.Min
		}
		if other.Max > g.Max {
			g.Max = other.Max
		}
		g.Last = other.Last
	}
}

// SynchronizedGauge augments a Gauge with a mutex to allow concurrent access from multiple
// goroutines.
type SynchronizedGauge struct {
	Gauge
	m sync.Mutex
}

// UpdateValue is a concurrent-safe wrapper around Gauge.UpdateValue
func (g *SynchronizedGauge) UpdateValue(val int64) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Gauge.UpdateValue(val)
}

// Update is a concurrent-safe wrapper around Gauge.Update
func (g *SynchronizedGauge) Update(other Gauge) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Gauge.Update(other)
}

// Reset returns a copy the current Gauge state and resets it back to its zero state.
func (g *SynchronizedGauge) Reset() Gauge {
	g.m.Lock()
	defer g.m.Unlock()
	current := g.Gauge
	g.Gauge = Gauge{}
	return current
}
