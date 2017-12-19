package appoptics

import (
	"runtime"
	"time"
)

const (
	runtimeRecordInterval = 10 * time.Second
)

func RecordRuntimeMetrics(m *MeasurementSet) {
	go recordRuntimeMetrics(m)
}

func recordRuntimeMetrics(m *MeasurementSet) {
	var (
		memStats       = &runtime.MemStats{}
		lastSampleTime = time.Now()
		lastPauseNs    uint64
		lastNumGC      uint64
	)

	for {
		runtime.ReadMemStats(memStats)

		now := time.Now()

		m.UpdateGaugeValue("go.goroutines", int64(runtime.NumGoroutine()))
		m.UpdateGaugeValue("go.memory.allocated", int64(memStats.Alloc))
		m.UpdateGaugeValue("go.memory.mallocs", int64(memStats.Mallocs))
		m.UpdateGaugeValue("go.memory.frees", int64(memStats.Frees))
		m.UpdateGaugeValue("go.memory.gc.total_pause", int64(memStats.PauseTotalNs))
		m.UpdateGaugeValue("go.memory.gc.heap", int64(memStats.HeapAlloc))
		m.UpdateGaugeValue("go..memory.gc.stack", int64(memStats.StackInuse))

		if lastPauseNs > 0 {
			pauseSinceLastSample := memStats.PauseTotalNs - lastPauseNs
			m.UpdateGaugeValue("go.memory.gc.pause_per_second", int64(float64(pauseSinceLastSample)/runtimeRecordInterval.Seconds()))
		}
		lastPauseNs = memStats.PauseTotalNs

		countGC := int(uint64(memStats.NumGC) - lastNumGC)
		if lastNumGC > 0 {
			diff := float64(countGC)
			diffTime := now.Sub(lastSampleTime).Seconds()
			m.UpdateGaugeValue("go.memory.gc.gc_per_second", int64(diff/diffTime))
		}

		if countGC > 0 {
			if countGC > 256 {
				countGC = 256
			}

			for i := 0; i < countGC; i++ {
				idx := int((memStats.NumGC-uint32(i))+255) % 256
				pause := time.Duration(memStats.PauseNs[idx])
				m.UpdateGaugeValue("go.memory.gc.pause", int64(pause))
			}
		}

		lastNumGC = uint64(memStats.NumGC)
		lastSampleTime = now

		time.Sleep(runtimeRecordInterval)
	}
}
