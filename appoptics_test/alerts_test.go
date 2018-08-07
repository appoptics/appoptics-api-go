package appoptics_test

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestAlertsService_List(t *testing.T) {
	alertResponse, err := client.AlertsService().List()
	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	query := alertResponse.Query
	alert := alertResponse.Alerts[0]
	condition := alert.Conditions[0]
	tag := condition.Tags[0]

	// query
	assert.Equal(t, 1, query.Found)
	assert.Equal(t, 1, query.Total)

	// Alert
	assert.Equal(t, 1234567, alert.ID)
	assert.Equal(t, "production.web.frontend.response_time", alert.Name)
	assert.Equal(t, "Web Response Time", alert.Description)

	// Conditions
	assert.Equal(t, 19376030, condition.ID)
	assert.Equal(t, "above", condition.Type)
	assert.Equal(t, "metric_name", condition.MetricName)
	assert.Equal(t, "threshold", condition.Threshold)
	assert.Equal(t, "summary_function", condition.SummaryFunction)

	// Tags
	assert.Equal(t, "tag_name", tag.Name)
	assert.Equal(t, false, tag.Grouped)
	assert.Equal(t, "tag_value", tag.Values[0])

}

//func TestAlertsService_Create(t *testing.T) {
//	alert, err := client.AlertsService().Create(&appoptics.Alert{})
//
//}
//
//func TestAlertsService_Retrieve(t *testing.T) {
//	alert, err := client.AlertsService().Retrieve(123)
//
//}
//
//func TestAlertsService_Status(t *testing.T) {
//
//}
