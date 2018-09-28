package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
)

func TestAlertsService_List(t *testing.T) {
	alertResponse, err := client.AlertsService().List()
	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	query := alertResponse.Query
	alert := alertResponse.Alerts[0]
	condition := alert.Conditions[0]
	service := alert.Services[0]

	// query
	assert.Equal(t, 1, query.Found)
	assert.Equal(t, 7, query.Total)

	// Alert
	assert.Equal(t, 1400310, *alert.ID)
	assert.Equal(t, "CPU.utilization", *alert.Name)
	assert.Equal(t, "How much power the box got", *alert.Description)

	// Conditions
	assert.Equal(t, 1016, *condition.ID)
	assert.Equal(t, "above", *condition.Type)
	assert.Equal(t, "AWS.EC2.CPUUtilization", *condition.MetricName)
	assert.Equal(t, float64(90), *condition.Threshold)
	assert.Equal(t, "max", *condition.SummaryFunction)

	// Service
	assert.Equal(t, 1153, *service.ID)
	assert.Equal(t, "mail", *service.Type)
	assert.Equal(t, "Ops Team", *service.Title)

	firstSetting := *service.Settings
	assert.Equal(t, "foo@domain.com,bar@domain.com", firstSetting["addresses"])

}

func TestAlertsService_Create(t *testing.T) {
	alert, err := client.AlertsService().Create(&appoptics.Alert{})
	if err != nil {
		t.Errorf("error running Create: %v", err)
	}
	condition := alert.Conditions[0]

	// Alert
	assert.Equal(t, 1234567, *alert.ID)
	assert.Equal(t, "production.web.frontend.response_time", *alert.Name)
	assert.Equal(t, "Web Response Time", *alert.Description)

	// Condition
	assert.Equal(t, 19376030, *condition.ID)
	assert.Equal(t, "above", *condition.Type)
	assert.Equal(t, "web.nginx.response_time", *condition.MetricName)
	assert.Equal(t, float64(200.0), *condition.Threshold)
	assert.Equal(t, "max", *condition.SummaryFunction)
}

func TestAlertsService_Retrieve(t *testing.T) {
	alert, err := client.AlertsService().Retrieve(123)

	if err != nil {
		t.Errorf("error running Retrieve: %v", err)
	}

	service := alert.Services[0]
	serviceSetting := *service.Settings

	// Alert
	assert.Equal(t, 123, *alert.ID)
	assert.Equal(t, "production.web.frontend.response_time", *alert.Name)
	assert.Equal(t, "Web Response Time", *alert.Description)

	// Services
	assert.Equal(t, 17584, *service.ID)
	assert.Equal(t, "slack", *service.Type)
	assert.Equal(t, "https://hooks.slack.com/services/XYZABC/a1b2c3/asdf", serviceSetting["url"])
}

func TestAlertsService_Status(t *testing.T) {
	alertStatus, err := client.AlertsService().Status(123)

	if err != nil {
		t.Errorf("error running Status: %v", err)
	}

	alert := *alertStatus.Alert

	assert.Equal(t, 120, *alert.ID)
	assert.Equal(t, "triggered", *alertStatus.Status)

}

func TestAlertsService_Associate(t *testing.T) {
	err := client.AlertsService().AssociateToService(123, 17584)

	if err != nil {
		t.Errorf("error running Status: %v", err)
	}
	assert.Equal(t, nil, err)
}

func TestAlertsService_Disassociate(t *testing.T) {
	err := client.AlertsService().DisassociateFromService(123, 17584)
	if err != nil {
		t.Errorf("error running Retrieve: %v", err)
	}

	assert.Equal(t, nil, err)
}
