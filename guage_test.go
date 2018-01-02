package appoptics

import "testing"

func TestUpdateValue(t *testing.T) {
	guage := &Gauge{
		Count: 3,
		Sum:   8,
		Min:   3,
		Max:   5,
		Last:  3,
	}

	t.Run("UpdateValue", func(t *testing.T) {
		newValue := int64(4)
		preUpdate := *guage
		guage.UpdateValue(newValue)

		newCount := preUpdate.Count + 1
		if guage.Count != newCount {
			t.Errorf("expected Count to be %d but was %d", newCount, guage.Count)
		}

		newSum := preUpdate.Sum + newValue
		if guage.Sum != newSum {
			t.Errorf("expected Sum to be %d but was %d", newSum, guage.Sum)
		}

		if guage.Min != 3 {
			t.Errorf("expected Min to be 3 but was %d", guage.Min)
		}

		if guage.Max != 5 {
			t.Errorf("expected Max to be 5 but was %d", guage.Max)
		}

		if guage.Last != 4 {
			t.Errorf("expected Last to be 4 but was %d", guage.Last)
		}
	})

	t.Run("UpdateValue with new Min", func(t *testing.T) {
		newMin := int64(1)
		guage.UpdateValue(newMin)

		if guage.Min != newMin {
			t.Errorf("expected Min to be %d but was %d", newMin, guage.Min)
		}
	})

	t.Run("UpdateValue with new Max", func(t *testing.T) {
		newMax := int64(7)
		guage.UpdateValue(newMax)
		if guage.Max != newMax {
			t.Errorf("expected Max to be %d but was %d", newMax, guage.Max)
		}
	})
}

func TestUpdateWithZeroValues(t *testing.T) {
	newGauge := Gauge{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	emptyGauge := &Gauge{}

	emptyGauge.Update(newGauge)

	if emptyGauge.Count != newGauge.Count {
		t.Errorf("expected Count to match but %d != %d", emptyGauge.Count, newGauge.Count)
	}

	if emptyGauge.Sum != newGauge.Sum {
		t.Errorf("expected Sum to match but %d != %d", emptyGauge.Sum, newGauge.Sum)
	}

	if emptyGauge.Min != newGauge.Min {
		t.Errorf("expected Min to match but %d != %d", emptyGauge.Min, newGauge.Min)
	}

	if emptyGauge.Max != newGauge.Max {
		t.Errorf("expected Max to match but %d != %d", emptyGauge.Max, newGauge.Max)
	}

	if emptyGauge.Last != newGauge.Last {
		t.Errorf("expected Last to match but %d != %d", emptyGauge.Last, newGauge.Last)
	}

}

func TestUpdateAggregation(t *testing.T) {
	oldGauge := Gauge{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	newGauge := Gauge{
		Count: 2,
		Sum:   5,
		Min:   2,
		Max:   3,
		Last:  3,
	}

	oldGauge.Update(newGauge)

	if oldGauge.Count != 4 {
		t.Errorf("expected Count to be aggregate but was %d", oldGauge.Count)
	}

	if oldGauge.Sum != 8 {
		t.Errorf("expected Sum to be aggregate but was %d", oldGauge.Sum)
	}
}

func TestUpdateWithNewMin(t *testing.T) {
	oldGauge := Gauge{
		Count: 2,
		Sum:   6,
		Min:   2,
		Max:   4,
		Last:  2,
	}

	newGauge := Gauge{
		Count: 2,
		Sum:   4,
		Min:   1,
		Max:   3,
		Last:  3,
	}

	oldGauge.Update(newGauge)

	if oldGauge.Min != newGauge.Min {
		t.Errorf("expected Min to be reset to %d but was %d", newGauge.Min, oldGauge.Min)
	}

}

func TestUpdateWithNewMax(t *testing.T) {
	oldGauge := Gauge{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	newGauge := Gauge{
		Count: 2,
		Sum:   4,
		Min:   1,
		Max:   3,
		Last:  3,
	}

	oldGauge.Update(newGauge)

	if oldGauge.Max != newGauge.Max {
		t.Errorf("expected Max to be reset to %d but was %d", newGauge.Max, oldGauge.Max)
	}
}
