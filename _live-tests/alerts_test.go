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

func TestAlerts(t *testing.T) {
	var (
		alert  *appoptics.Alert
		metric *appoptics.Metric
		err    error
	)

	metric, err = client.MetricsService().Create(testMetric("test.alert"))
	require.Nil(t, err)

	defer client.MetricsService().Delete(metric.Name)

	t.Run("Create", func(t *testing.T) {
		newAlert, err := client.AlertsService().Create(testAlert("test", metric.Name))
		require.Nil(t, err)
		assert.Equal(t, testAlert("test", metric.Name).Name, newAlert.Name)
		alert = newAlert
	})

	t.Run("List", func(t *testing.T) {
		_, err := client.AlertsService().List()
		assert.Nil(t, err)
	})

	t.Run("Retrieve", func(t *testing.T) {
		fetchedAlert, err := client.AlertsService().Retrieve(alert.ID)
		assert.Nil(t, err)
		assert.Equal(t, alert.Name, fetchedAlert.Name)
	})

	t.Run("Update", func(t *testing.T) {
		alert.Name = "other.name"
		err := client.AlertsService().Update(alert)
		assert.Nil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.AlertsService().Delete(alert.ID)
		require.Nil(t, err)
	})
}

func testAlert(n, metricName string) *appoptics.Alert {
	alertName := fmt.Sprintf("%s-%s", TestPrefix, n)
	return &appoptics.Alert{
		Name:        alertName,
		Description: "A Test Alert",
		Conditions: []*appoptics.AlertCondition{
			{
				Type:       "above",
				MetricName: metricName,
				Threshold:  200,
			},
		},
	}
}
