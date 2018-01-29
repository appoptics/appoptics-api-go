package appoptics

import "testing"

func TestUpdateValue(t *testing.T) {
	s := &Aggregator{
		Count: 3,
		Sum:   8,
		Min:   3,
		Max:   5,
		Last:  3,
	}

	t.Run("UpdateValue", func(t *testing.T) {
		newValue := float64(4)
		preUpdate := *s
		s.UpdateValue(newValue)

		newCount := preUpdate.Count + 1
		if s.Count != newCount {
			t.Errorf("expected Count to be %v but was %v", newCount, s.Count)
		}

		newSum := preUpdate.Sum + newValue
		if s.Sum != newSum {
			t.Errorf("expected Sum to be %v but was %v", newSum, s.Sum)
		}

		if s.Min != 3 {
			t.Errorf("expected Min to be 3 but was %v", s.Min)
		}

		if s.Max != 5 {
			t.Errorf("expected Max to be 5 but was %v", s.Max)
		}

		if s.Last != 4 {
			t.Errorf("expected Last to be 4 but was %v", s.Last)
		}
	})

	t.Run("UpdateValue with new Min", func(t *testing.T) {
		newMin := float64(1)
		s.UpdateValue(newMin)

		if s.Min != newMin {
			t.Errorf("expected Min to be %v but was %v", newMin, s.Min)
		}
	})

	t.Run("UpdateValue with new Max", func(t *testing.T) {
		newMax := float64(7)
		s.UpdateValue(newMax)
		if s.Max != newMax {
			t.Errorf("expected Max to be %v but was %v", newMax, s.Max)
		}
	})
}

func TestUpdateWithZeroValues(t *testing.T) {
	newAgg := Aggregator{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	emptyAgg := &Aggregator{}

	emptyAgg.Update(newAgg)

	if emptyAgg.Count != newAgg.Count {
		t.Errorf("expected Count to match but %v != %v", emptyAgg.Count, newAgg.Count)
	}

	if emptyAgg.Sum != newAgg.Sum {
		t.Errorf("expected Sum to match but %v != %v", emptyAgg.Sum, newAgg.Sum)
	}

	if emptyAgg.Min != newAgg.Min {
		t.Errorf("expected Min to match but %v != %v", emptyAgg.Min, newAgg.Min)
	}

	if emptyAgg.Max != newAgg.Max {
		t.Errorf("expected Max to match but %v != %v", emptyAgg.Max, newAgg.Max)
	}

	if emptyAgg.Last != newAgg.Last {
		t.Errorf("expected Last to match but %v != %v", emptyAgg.Last, newAgg.Last)
	}

}

func TestUpdateAggregation(t *testing.T) {
	oldAgg := Aggregator{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	newAgg := Aggregator{
		Count: 2,
		Sum:   5,
		Min:   2,
		Max:   3,
		Last:  3,
	}

	oldAgg.Update(newAgg)

	if oldAgg.Count != 4 {
		t.Errorf("expected Count to be aggregate but was %v", oldAgg.Count)
	}

	if oldAgg.Sum != 8 {
		t.Errorf("expected Sum to be aggregate but was %v", oldAgg.Sum)
	}
}

func TestUpdateWithNewMin(t *testing.T) {
	oldAgg := Aggregator{
		Count: 2,
		Sum:   6,
		Min:   2,
		Max:   4,
		Last:  2,
	}

	newAgg := Aggregator{
		Count: 2,
		Sum:   4,
		Min:   1,
		Max:   3,
		Last:  3,
	}

	oldAgg.Update(newAgg)

	if oldAgg.Min != newAgg.Min {
		t.Errorf("expected Min to be reset to %v but was %v", newAgg.Min, oldAgg.Min)
	}

}

func TestUpdateWithNewMax(t *testing.T) {
	oldAgg := Aggregator{
		Count: 2,
		Sum:   3,
		Min:   1,
		Max:   2,
		Last:  2,
	}

	newAgg := Aggregator{
		Count: 2,
		Sum:   4,
		Min:   1,
		Max:   3,
		Last:  3,
	}

	oldAgg.Update(newAgg)

	if oldAgg.Max != newAgg.Max {
		t.Errorf("expected Max to be reset to %v but was %v", newAgg.Max, oldAgg.Max)
	}
}
