package integration_tests

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/require"
)

func TestCharts(t *testing.T) {
	var (
		chartID int
		spaceID int
	)

	space, err := client.SpacesService().Create(spaceName("chart"))
	require.Nil(t, err)
	spaceID = space.ID

	defer client.SpacesService().Delete(spaceID)

	t.Run("Create", func(t *testing.T) {
		chart, err := client.ChartsService().Create(chartFixture("test"), spaceID)
		require.Nil(t, err)
		require.Equal(t, chartFixture("test").Name, chart.Name)
		chartID = chart.ID
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.ChartsService().Delete(chartID, spaceID)
		require.Nil(t, err)
	})
}

func chartFixture(name string) *appoptics.Chart {
	return &appoptics.Chart{
		Name: name,
		Type: "line",
	}
}
