package appoptics

type MeasurementSetReport struct {
	Counts map[string]int64
	Gauges map[string]Gauge
}

func NewMeasurementSetReport() *MeasurementSetReport {
	return &MeasurementSetReport{
		Counts: map[string]int64{},
		Gauges: map[string]Gauge{},
	}
}
