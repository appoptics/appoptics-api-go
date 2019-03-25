package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsService_List(t *testing.T) {
	apiMetricsResponse, err := client.MetricsService().List()

	require.Nil(t, err)

	query := apiMetricsResponse.Query

	firstMetric := apiMetricsResponse.Metrics[0]
	firstAttributes := firstMetric.Attributes

	// Query
	assert.Equal(t, 2, query.Found)
	assert.Equal(t, 2, query.Length)
	assert.Equal(t, 2, query.Total)

	// Metric metadata
	assert.Equal(t, "app_requests", firstMetric.Name)
	assert.Equal(t, "app_requests", firstMetric.DisplayName)
	assert.Equal(t, "HTTP requests serviced by the app per-minute", firstMetric.Description)
	assert.Equal(t, 60, firstMetric.Period)
	assert.Equal(t, "gauge", firstMetric.Type)

	// Attributes
	assert.Equal(t, "appoptics-metrics/0.7.4 (ruby; 1.9.3p194; x86_64-linux) direct-faraday/0.8.4", firstAttributes.CreatedByUA)
	assert.Equal(t, float64(0), firstAttributes.DisplayMin)
	assert.Equal(t, true, firstAttributes.DisplayStacked)
	assert.Equal(t, "Requests", firstAttributes.DisplayUnitsLong)
	assert.Equal(t, "reqs", firstAttributes.DisplayUnitsShort)

}

func TestMetricsService_Create(t *testing.T) {
	apiMetricsResponse, err := client.MetricsService().Create(&appoptics.Metric{Name: "cpu.percent.used"})
	require.Nil(t, err)

	assert.Equal(t, "cpu.percent.used", apiMetricsResponse.Name)
	assert.Equal(t, "CPU Used", apiMetricsResponse.DisplayName)
	assert.Equal(t, "composite", apiMetricsResponse.Type)
	assert.Equal(t, "all the cpu used on the machine", apiMetricsResponse.Description)
}

func TestMetricsService_Retrieve(t *testing.T) {
	apiMetricsResponse, err := client.MetricsService().Retrieve("cpu.temp")

	require.Nil(t, err)

	assert.Equal(t, "cpu_temp", apiMetricsResponse.Name)
	assert.Equal(t, "cpu_temp", apiMetricsResponse.DisplayName)
	assert.Equal(t, "gauge", apiMetricsResponse.Type)
}

// Yes, this test doesn't really do anything
func TestMetricsService_Update(t *testing.T) {
	err := client.MetricsService().Update(&appoptics.MetricAttributes{})
	require.Nil(t, err)
}

// Yes, this test doesn't really do anything
func TestMetricsService_Delete(t *testing.T) {
	err := client.MetricsService().Delete("cpu.temp")
	require.Nil(t, err)
}
