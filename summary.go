package appoptics

import "sync"

// A Summary corresponds roughly to the "summary fields" measurement feature in AppOptics, storing a count/sum/min/max/last.
// It can be either updated by passing sequential values to UpdateValue or by passing a Summary
// to Update, e.g. s.Update(Summary{Sum:100,Count:10,Min:5,Max:15})
type Summary struct {
	Count int64
	Sum   int64
	Min   int64
	Max   int64
	Last  int64
}

// UpdateValue sets the most recently observed value for this Summary, updating sum/count/min/max/last
// accordingly.
func (g *Summary) UpdateValue(val int64) {
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

// Update merges another Summary into this Summary, merging sum/count/min/max/last accordingly. It can
// be used to facilitate efficient input of many data points into a Summary in one call, and it can
// also be used to merge two different Summarys (for example, workers can each maintain their own
// and periodically merge them).
func (g *Summary) Update(other Summary) {
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

// SynchronizedSummary augments a Summary with a mutex to allow concurrent access from multiple
// goroutines.
type SynchronizedSummary struct {
	Summary
	m sync.Mutex
}

// UpdateValue is a concurrent-safe wrapper around Summary.UpdateValue
func (g *SynchronizedSummary) UpdateValue(val int64) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Summary.UpdateValue(val)
}

// Update is a concurrent-safe wrapper around Summary.Update
func (g *SynchronizedSummary) Update(other Summary) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Summary.Update(other)
}

// Reset returns a copy the current Summary state and resets it back to its zero state.
func (g *SynchronizedSummary) Reset() Summary {
	g.m.Lock()
	defer g.m.Unlock()
	current := g.Summary
	g.Summary = Summary{}
	return current
}
