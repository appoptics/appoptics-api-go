package live_tests

import (
	"testing"

	"fmt"

	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
//
// THESE TESTS PRODUCE LIVE DATA IN APPOPTICS - ORDER MATTERS HERE!
//
//

func TestCharts(t *testing.T) {
	var (
		chartID int
		spaceID int
		chart   *appoptics.Chart
	)

	space, err := client.SpacesService().Create(spaceName("chart"))
	require.Nil(t, err)
	spaceID = space.ID

	time.Sleep(4 * time.Second) // account for potential replica lag

	defer client.SpacesService().Delete(spaceID)

	t.Run("Create", func(t *testing.T) {
		createdChart, err := client.ChartsService().Create(chartFixture("test"), spaceID)
		require.Nil(t, err)
		require.Equal(t, chartFixture("test").Name, createdChart.Name)
		chart = createdChart
		chartID = chart.ID
	})

	t.Run("List", func(t *testing.T) {
		charts, err := client.ChartsService().List(spaceID)
		require.Nil(t, err)
		firstChart := charts[0]
		assert.Equal(t, chartID, firstChart.ID)
	})

	t.Run("Retrieve", func(t *testing.T) {
		retrievedChart, err := client.ChartsService().Retrieve(chartID, spaceID)
		require.Nil(t, err)
		assert.Equal(t, chart.Name, retrievedChart.Name)
	})

	t.Run("Update", func(t *testing.T) {
		fmt.Printf("UPDATE CHART: %+v\n", chart)
		otherName := "new-name"
		chart.Name = otherName
		updatedChart, err := client.ChartsService().Update(chart, spaceID)
		require.Nil(t, err)
		assert.Equal(t, otherName, updatedChart.Name)
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
