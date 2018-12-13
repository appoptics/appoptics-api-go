package appoptics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveJobs(t *testing.T) {
	jobsResponse, err := client.JobsService().Retrieve(123456)
	require.Nil(t, err)

	assert.Equal(t, 123456, jobsResponse.ID)
	assert.Equal(t, float64(76.5), jobsResponse.Progress)
	assert.Equal(t, "failed", jobsResponse.State)
	assert.Equal(t, "is invalid", jobsResponse.Errors["name"][0])
}
