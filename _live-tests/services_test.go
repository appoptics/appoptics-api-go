package live_tests

import (
	"testing"

	appoptics "github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
//
// THESE TESTS PRODUCE LIVE DATA IN APPOPTICS - ORDER MATTERS HERE!
//
//

func TestServices(t *testing.T) {
	var (
		service *appoptics.Service
	)

	t.Run("Create", func(t *testing.T) {
		createdService, err := client.ServicesService().Create(testService("test"))
		require.Nil(t, err)
		assert.Equal(t, testService("test").Title, createdService.Title)
		service = createdService
	})

	t.Run("List", func(t *testing.T) {
		_, err := client.ServicesService().List()
		require.Nil(t, err)
	})

	t.Run("Retrieve", func(t *testing.T) {
		retrievedService, err := client.ServicesService().Retrieve(service.ID)
		require.Nil(t, err)
		assert.Equal(t, retrievedService.Title, service.Title)
		assert.Equal(t, retrievedService.Type, service.Type)
		assert.Equal(t, retrievedService.Settings["token"], service.Settings["token"])
	})

	t.Run("Update", func(t *testing.T) {
		title := "new-title"
		service.Title = title
		err := client.ServicesService().Update(service)
		require.Nil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.ServicesService().Delete(service.ID)
		require.Nil(t, err)
	})
}

func testService(title string) *appoptics.Service {
	return &appoptics.Service{
		Title: title,
		Type:  "slack",
		Settings: map[string]string{
			"room":      "deployments",
			"token":     "deadbeef",
			"subdomain": "acme",
			"url":       "https://example.com",
		},
	}
}
