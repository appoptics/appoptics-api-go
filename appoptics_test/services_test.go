package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
)

func TestServicesService_List(t *testing.T) {
	serviceResponse, err := client.ServicesService().List()

	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	query := serviceResponse.Query

	assert.Equal(t, 2, query.Found)
	assert.Equal(t, 2, query.Length)
	assert.Equal(t, 0, query.Offset)
	assert.Equal(t, 2, query.Total)

	firstService := serviceResponse.Services[0]

	assert.Equal(t, 145, *firstService.ID)
	assert.Equal(t, "slack", *firstService.Type)
	assert.Equal(t, "Notify Ops Room", *firstService.Title)

	settings := *firstService.Settings

	assert.Equal(t, "Ops", settings["room"])
	assert.Equal(t, "1234567890ABCDEF", settings["token"])
	assert.Equal(t, "acme", settings["subdomain"])

}

func TestServicesService_Retrieve(t *testing.T) {
	service, err := client.ServicesService().Retrieve(123)

	if err != nil {
		t.Errorf("error running Retrieve: %v", err)
	}

	assert.Equal(t, 156, *service.ID)
	assert.Equal(t, "mail", *service.Type)
	assert.Equal(t, "Email ops team", *service.Title)
}

func TestServicesService_Create(t *testing.T) {
	service, err := client.ServicesService().Create(&appoptics.Service{})

	if err != nil {
		t.Errorf("error running Create: %v", err)
	}

	assert.Equal(t, 145, *service.ID)
	assert.Equal(t, "campfire", *service.Type)
	assert.Equal(t, "Notify Ops Room", *service.Title)
}
