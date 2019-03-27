package live_tests

import (
	"testing"

	"fmt"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
//
// THESE TESTS PRODUCE LIVE DATA IN APPOPTICS - ORDER MATTERS HERE!
//
//

func TestMetrics(t *testing.T) {
	var (
		metricName string
		metric     *appoptics.Metric
	)

	t.Run("Create", func(t *testing.T) {
		newMetric, err := client.MetricsService().Create(testMetric("test"))
		require.Nil(t, err)
		require.Equal(t, testMetric("test").Name, newMetric.Name)
		metric = newMetric
		metricName = metric.Name
	})

	t.Run("List", func(t *testing.T) {
		_, err := client.MetricsService().List()
		require.Nil(t, err)
	})

	t.Run("Retrieve", func(t *testing.T) {
		metric, err := client.MetricsService().Retrieve(metricName)
		require.Nil(t, err)
		assert.Equal(t, testMetric("test").Name, metric.Name)
	})

	t.Run("Update", func(t *testing.T) {
		metric.Attributes.Color = "#beefee"
		err := client.MetricsService().Update(metricName, metric)
		assert.Nil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.MetricsService().Delete(metricName)
		assert.Nil(t, err)
	})
}

func testMetric(n string) *appoptics.Metric {
	name := fmt.Sprintf("%s-%s", TestPrefix, n)
	return &appoptics.Metric{
		Name: name,
		Type: "gauge",
		Attributes: appoptics.MetricAttributes{
			Color:     "#deadbe",
			Aggregate: false,
		},
	}
}
