package appoptics_test

import (
	"testing"

	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnotationsService_List(t *testing.T) {
	t.Run("without arg", func(t *testing.T) {
		annResponse, err := client.AnnotationsService().List(nil)

		require.Nil(t, err)

		streams := annResponse.AnnotationStreams

		assert.Equal(t, "api-deploys", streams[0].Name)
		assert.Equal(t, "Deploys to API", streams[0].DisplayName)

		assert.Equal(t, "app-deploys", streams[1].Name)
		assert.Equal(t, "Deploys to UX app", streams[1].DisplayName)
	})

	t.Run("with arg", func(t *testing.T) {
		filter := "foobar"
		annResponse, err := client.AnnotationsService().List(&filter)

		require.Nil(t, err)

		streams := annResponse.AnnotationStreams

		assert.Equal(t, "api-deploys", streams[0].Name)
		assert.Equal(t, "Deploys to API", streams[0].DisplayName)

		assert.Equal(t, "app-deploys", streams[1].Name)
		assert.Equal(t, "Deploys to UX app", streams[1].DisplayName)
	})
}

func TestAnnotationsService_Create(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(time.Minute * time.Duration(5))
	newAnnotation := &appoptics.AnnotationEvent{
		Title:       "Deafbeef Deploys",
		Source:      "non-prod-eu-west-01",
		Description: "Shipped Beef Vortex 2.7",
		StartTime:   startTime.Unix(),
		EndTime:     endTime.Unix(),
	}

	annResponse, err := client.AnnotationsService().Create(newAnnotation, "dead-beef-app-deploys")
	require.Nil(t, err)
	require.NotNil(t, annResponse)

	assert.Equal(t, 123, annResponse.ID)
	assert.Equal(t, "My Annotation", annResponse.Title)
	assert.Equal(t, int64(1234567890), annResponse.StartTime)
}

func TestAnnotationsService_Retrieve(t *testing.T) {
	rar := &appoptics.RetrieveAnnotationsRequest{Name: "foobar"} // name doesn't matter in this test
	annResponse, err := client.AnnotationsService().Retrieve(rar)
	require.Nil(t, err)

	assert.Equal(t, "api-deploys", annResponse.Name)
	assert.Equal(t, 45, annResponse.Events[0]["unassigned"][0].ID)
	assert.Equal(t, "Deployed v91", annResponse.Events[0]["unassigned"][0].Title)
}

func TestAnnotationsService_RetrieveEvent(t *testing.T) {
	annResponse, err := client.AnnotationsService().RetrieveEvent("foobar", 42)
	require.Nil(t, err)

	assert.Equal(t, 189, annResponse.ID)
	assert.Equal(t, int64(1234567890), annResponse.StartTime)

	assert.Equal(t, "Github commit", annResponse.Links[0].Label)
	assert.Equal(t, "github", annResponse.Links[0].Rel)
}

func TestAnnotationsService_UpdateStream(t *testing.T) {
	err := client.AnnotationsService().UpdateStream("foo", "bar")
	require.Nil(t, err)
}

func TestAnnotationsService_UpdateEvent(t *testing.T) {
	link := &appoptics.AnnotationLink{} // blank b/c we're testing client processing fixture data
	annResponse, err := client.AnnotationsService().UpdateEvent("foobar", 123, link)
	require.Nil(t, err)

	assert.Equal(t, "github", annResponse.Rel)
	assert.Equal(t, "https://github.com/acme/app/commits/01beaf", annResponse.Href)
}
