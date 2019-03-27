package live_tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
//
// THESE TESTS PRODUCE LIVE DATA IN APPOPTICS - ORDER MATTERS HERE!
//
//

func TestSpaces(t *testing.T) {
	var spaceID int

	t.Run("Create", func(t *testing.T) {
		space, err := client.SpacesService().Create(spaceName(CreatedNameString))
		require.Nil(t, err)
		assert.Equal(t, spaceName(CreatedNameString), space.Name)
		spaceID = space.ID
	})

	t.Run("Retrieve", func(t *testing.T) {
		space, err := client.SpacesService().Retrieve(spaceID)
		require.Nil(t, err)
		assert.Equal(t, spaceName(CreatedNameString), space.Name)
	})

	t.Run("List", func(t *testing.T) {
		_, err := client.SpacesService().List(client.DefaultRequestParameters(5))
		require.Nil(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		err := client.SpacesService().Update(spaceID, spaceName(UpdatedNameString))
		require.Nil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.SpacesService().Delete(spaceID)
		assert.Nil(t, err)
	})
}

func spaceName(s string) string {
	return fmt.Sprintf("%s-Space-%s", TestPrefix, s)
}
