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

	metric := firstChartStream.Metric
	gauge := firstChartStream.Type
	firstTagName := firstChartStreamTags[0].Name
	firstTagValue := firstChartStreamTags[0].Values[0]

	assert.Equal(t, 27035309, firstChartStream.ID)
	assert.Equal(t, "cpu.percent.idle", metric)
	assert.Equal(t, "gauge", gauge)
	assert.Equal(t, "environment", firstTagName)
	assert.Equal(t, "*", firstTagValue)
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

	tagName := stream.Tags[0].Name
	tagValue := stream.Tags[0].Values[0]

	assert.Equal(t, 27003258, stream.ID)
	assert.Equal(t, "cpu.percent.idle", stream.Metric)
	assert.Equal(t, "gauge", stream.Type)
	assert.Equal(t, "region", tagName)
	assert.Equal(t, "us-east-1", tagValue)
}

func TestChartsService_Create(t *testing.T) {
	chart, err := client.ChartsService().Create(&appoptics.Chart{}, 123)
	if err != nil {
		t.Errorf("error running Create: %v", err)
	}

	assert.Equal(t, 1234567, chart.ID)
	assert.Equal(t, "CPU Usage", chart.Name)
	assert.Equal(t, "line", chart.Type)

	firstStream := chart.Streams[0]

	assert.Equal(t, 27032885, firstStream.ID)
	assert.Equal(t, "cpu.percent.idle", firstStream.Metric)
	assert.Equal(t, "gauge", firstStream.Type)

	secondStreamTagName := chart.Streams[1].Tags[0].Name
	secondStreamTagValue := chart.Streams[1].Tags[0].Values[0]

	assert.Equal(t, "environment", secondStreamTagName)
	assert.Equal(t, "prod", secondStreamTagValue)
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

	streamFirstTag := stream.Tags[0]
	streamFirstTagValues := streamFirstTag.Values[0]

	assert.Equal(t, "region", streamFirstTag.Name)
	assert.Equal(t, "us-east-1", streamFirstTagValues)
}
