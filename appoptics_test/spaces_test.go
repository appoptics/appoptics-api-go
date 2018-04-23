package appoptics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpacesService_Create(t *testing.T) {
	space, err := client.SpacesService().Create("CPUs")
	if err != nil {
		t.Errorf("error running Create: %v", err)
	}

	assert.Equal(t, space.Name, "CPUs")
	assert.NotEmpty(t, space.ID)
}

func TestSpacesService_Retrieve(t *testing.T) {
	resp, err := client.SpacesService().Retrieve(129)

	if err != nil {
		t.Errorf("error running Retrieve: %v", err)
	}

	assert.Equal(t, resp.Name, "CPUs")
	assert.Equal(t, resp.ID, 129)
	assert.Equal(t, len(resp.Charts), 4)
}

func TestSpacesService_List(t *testing.T) {
	spaces, err := client.SpacesService().List(nil)

	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	assert.Equal(t, len(spaces), 1)
	assert.Equal(t, spaces[0].ID, 4)
	assert.Equal(t, spaces[0].Name, "staging_ops")
}

func TestSpacesService_Update(t *testing.T) {
	space, err := client.SpacesService().Update(1, "MEMORY")

	if err != nil {
		t.Errorf("error running List: %v", err)
	}

	assert.Equal(t, space.Name, "MEMORY")
}

func TestSpacesService_Delete(t *testing.T) {
	err := client.SpacesService().Delete(1)
	if err != nil {
		t.Errorf("error running Delete: %v", err)
	}
}
