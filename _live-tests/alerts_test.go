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
		alert   *appoptics.Alert
		metric  *appoptics.Metric
		service *appoptics.Service
		err     error
	)

	metric, err = client.MetricsService().Create(testMetric("test.alert"))
	require.Nil(t, err)

	service, err = client.ServicesService().Create(testService("test"))
	require.Nil(t, err)

	defer client.MetricsService().Delete(metric.Name)
	defer client.ServicesService().Delete(service.ID)

	t.Run("Create", func(t *testing.T) {
		newAlert, err := client.AlertsService().Create(testAlertRequest("test", metric.Name))
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

	t.Run("AssociateToService", func(t *testing.T) {
		err := client.AlertsService().AssociateToService(alert.ID, service.ID)
		assert.Nil(t, err)
	})

	t.Run("DisassociateFromService", func(t *testing.T) {
		err := client.AlertsService().DisassociateFromService(alert.ID, service.ID)
		assert.Nil(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		alert.Name = "other.name"
		newRequest := testAlertRequest(alert.Name, metric.Name)
		newRequest.ID = alert.ID
		err := client.AlertsService().Update(newRequest)
		assert.Nil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.AlertsService().Delete(alert.ID)
		require.Nil(t, err)
	})
}

func testAlertRequest(n, metricName string) *appoptics.AlertRequest {
	alertName := fmt.Sprintf("%s-%s", TestPrefix, n)
	return &appoptics.AlertRequest{
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
