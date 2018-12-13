package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSnapshot(t *testing.T) {
	sResponse, err := client.SnapshotsService().Create(&appoptics.Snapshot{})
	require.Nil(t, err)

	assert.Equal(t, "https://api.appoptics.com/v1/jobs/123456", sResponse.JobHref)
	assert.Equal(t, 1, sResponse.Subject["chart"].ID)
	assert.Equal(t, "*", sResponse.Subject["chart"].Sources[0])
	assert.False(t, sResponse.EndTime.IsZero())
	assert.False(t, sResponse.CreatedAt.IsZero())
	assert.False(t, sResponse.UpdatedAt.IsZero())
}

func TestRetrieveSnapshot(t *testing.T) {
	sResponse, err := client.SnapshotsService().Retrieve(123)
	require.Nil(t, err)

	assert.Equal(t, "https://api.appoptics.com/v1/jobs/123456", sResponse.JobHref)
	assert.Equal(t, 1, sResponse.Subject["chart"].ID)
	assert.Equal(t, "*", sResponse.Subject["chart"].Sources[0])
	assert.False(t, sResponse.EndTime.IsZero())
	assert.False(t, sResponse.CreatedAt.IsZero())
	assert.False(t, sResponse.UpdatedAt.IsZero())
}
