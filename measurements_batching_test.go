package appoptics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
				bp.stopBatchingChan <- struct{}{}
			}
		}

		assert.Equal(t, batchCount, len(batchCollection))
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
				bp.stopBatchingChan <- struct{}{}
			}
		}

		measurementsInBatch := len(batchCollection[0].Measurements)

		assert.Equal(t, measurementsCount, measurementsInBatch)
	})

	t.Run("errors accumulate and stop batching at limit", func(t *testing.T) {
		bp := NewBatchPersister(mms, false)
		var pSig struct{}
		var eSig struct{}
		go bp.batchMeasurements()
		go bp.managePersistenceErrors()

		for i := 0; i < bp.errorLimit; i++ {
			bp.errorChan <- fmt.Errorf("some error")
		}

		pSig = <-bp.stopPersistingChan
		eSig = <-bp.stopErrorChan

		assert.NotNil(t, pSig)
		assert.NotNil(t, eSig)
		assert.Equal(t, bp.errorLimit, len(bp.errors))
	})
}
