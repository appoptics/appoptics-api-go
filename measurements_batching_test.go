package appoptics

import (
	"fmt"
	"testing"
	"time"
)

func TestBatchPersister(t *testing.T) {
	mms := &MockMeasurementsService{}

	t.Run("creates maximal batches", func(t *testing.T) {
		bp := NewBatchPersister(mms, false)
		batchCount := 2
		batchCollection := []*MeasurementsBatch{}

		go bp.batchMeasurements()

		go func() {
			for i := 0; i < (MeasurementPostMaxBatchSize * batchCount); i++ {
				// MeasurementsSink() is write-only bp.prepChan
				bp.MeasurementsSink() <- []Measurement{{Value: 3.14, Time: time.Now().UTC().Unix()}}
			}
			close(bp.prepChan)
		}()

		for batch := range bp.batchChan {
			batchCollection = append(batchCollection, batch)
			if len(batchCollection) == batchCount {
				bp.stopBatchingChan <- true
			}
		}

		if len(batchCollection) != batchCount {
			t.Errorf("expected batch count to be %d but was %d", batchCount, len(batchCollection))
		}
	})

	t.Run("respects push interval", func(t *testing.T) {
		bp := NewBatchPersister(mms, false)
		batchCollection := []*MeasurementsBatch{}
		measurementsCount := 5

		if measurementsCount > MeasurementPostMaxBatchSize {
			t.Errorf("measurements count larger than batch size makes no sense in this test")
		}

		bp.SetMaximumPushInterval(100)

		go bp.batchMeasurements()

		go func() {
			for i := 0; i < measurementsCount; i++ {
				// MeasurementsSink() is write-only bp.prepChan
				bp.MeasurementsSink() <- []Measurement{{Value: 3.14, Time: time.Now().UTC().Unix()}}
			}
			close(bp.prepChan)
		}()

		// grab the single batch that will be on the channel
		for batch := range bp.batchChan {
			batchCollection = append(batchCollection, batch)
			if len(batchCollection[0].Measurements) == measurementsCount {
				bp.stopBatchingChan <- true
			}
		}

		measurementsInBatch := len(batchCollection[0].Measurements)

		if measurementsInBatch != measurementsCount {
			t.Errorf("expected Measurements in single batch to be %d but was %d", measurementsCount, measurementsInBatch)
		}
	})

	t.Run("errors accumulate and stop batching at limit", func(t *testing.T) {
		bp := NewBatchPersister(mms, false)
		var pSig bool
		var eSig bool
		go bp.batchMeasurements()
		go bp.managePersistenceErrors()

		for i := 0; i < bp.errorLimit; i++ {
			bp.errorChan <- fmt.Errorf("some error")
		}

		pSig = <-bp.stopPersistingChan
		eSig = <-bp.stopErrorChan

		if !pSig {
			t.Errorf("expected to find a value on the stopPersistingChan")
		}

		if !eSig {
			t.Errorf("expected to find a value on the stopPersistingChan")
		}

		if len(bp.errors) != bp.errorLimit {
			t.Errorf("expected errors count (%d) to equal error limit (%d) but it did not", len(bp.errors), bp.errorLimit)
		}

	})
}
