package appoptics

type MeasurementSetReport struct {
	Counts    map[string]int64
	Summaries map[string]Summary
}

func NewMeasurementSetReport() *MeasurementSetReport {
	return &MeasurementSetReport{
		Counts:    map[string]int64{},
		Summaries: map[string]Summary{},
	}
}
