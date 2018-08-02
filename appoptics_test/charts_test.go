package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
)

func TestChartsService_List(t *testing.T) {
	charts, err := client.ChartsService().List(123)

	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	assert.Equal(t, 1, len(charts))

	firstChart := charts[0]
	assert.Equal(t, 2, len(firstChart.Streams))
	firstChartStream := firstChart.Streams[0]
	firstChartStreamTags := firstChartStream.Tags

	assert.Equal(t, 27035309, firstChartStream.ID)
	assert.Equal(t, "cpu.percent.idle", firstChartStream.Metric)
	assert.Equal(t, "gauge", firstChartStream.Type)
	assert.Equal(t, "environment", firstChartStreamTags[0].Name)
	assert.Equal(t, "*", firstChartStreamTags[0].Values[0])
}

func TestChartsService_Retrieve(t *testing.T) {
	chart, err := client.ChartsService().Retrieve(1, 2)

	if err != nil {
		t.Errorf("error running Retrieve: %v", err)
	}

	assert.Equal(t, 3700969, chart.ID)
	assert.Equal(t, "CPU Usage", chart.Name)
	assert.Equal(t, "line", chart.Type)

	stream := chart.Streams[0]

	assert.Equal(t, 27003258, stream.ID)
	assert.Equal(t, "cpu.percent.idle", stream.Metric)
	assert.Equal(t, "gauge", stream.Type)
	assert.Equal(t, "region", stream.Tags[0].Name)
	assert.Equal(t, "us-east-1", stream.Tags[0].Values[0])
}

func TestChartsService_Create(t *testing.T) {
	chart, err := client.ChartsService().Create(&appoptics.Chart{}, 123)
	if err != nil {
		t.Errorf("error running Create: %v", err)
	}

	assert.Equal(t, 1234567, chart.ID)
	assert.Equal(t, "CPU Usage", chart.Name)
	assert.Equal(t, "line", chart.Type)

	assert.Equal(t, 27032885, chart.Streams[0].ID)
	assert.Equal(t, "cpu.percent.idle", chart.Streams[0].Metric)
	assert.Equal(t, "gauge", chart.Streams[0].Type)

	assert.Equal(t, "environment", chart.Streams[1].Tags[0].Name)
	assert.Equal(t, "prod", chart.Streams[1].Tags[0].Values[0])
}

func TestChartsService_Update(t *testing.T) {
	chart, err := client.ChartsService().Update(&appoptics.Chart{}, 123)

	if err != nil {
		t.Errorf("error running Update: %v", err)
	}

	assert.Equal(t, 3700969, chart.ID)
	assert.Equal(t, "Temperature", chart.Name)
	assert.Equal(t, "line", chart.Type)

	stream := chart.Streams[0]

	assert.Equal(t, 27003258, stream.ID)
	assert.Equal(t, "collectd.cpu.0.cpu.user", stream.Metric)
	assert.Equal(t, "gauge", stream.Type)
	assert.Equal(t, "region", stream.Tags[0].Name)
	assert.Equal(t, "us-east-1", stream.Tags[0].Values[0])
}

// TODO: do we even care to have this since there's no structure to process/test?
func TestChartsService_Delete(t *testing.T) {
	err := client.ChartsService().Delete(1, 1)
	if err != nil {
		t.Errorf("error running Delete: %v", err)
	}
}
