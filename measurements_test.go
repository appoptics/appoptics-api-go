package appoptics

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	t    *testing.T
	reqs []*http.Request
}

var errMockHTTPClientExpected = errors.New("mock http.Response not implemented")

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, req)
	return nil, errMockHTTPClientExpected
}

// Test Create works by using a mock HTTP client and checking the request resulting from
// posting a MeasurementsBatch.
func TestMeasurementsServiceCreate(t *testing.T) {
	// prep batch
	ts1 := time.Now().Unix()
	ts2 := time.Now().Unix() - 100

	testBatch := &MeasurementsBatch{
		Measurements: []Measurement{
			{Name: "metric1", Tags: map[string]string{"a": "b"}, Value: 5},
			{Name: "metric2", Tags: map[string]string{"c": "d"}, Sum: 10, Count: 2},
			{Name: "metric3", Tags: map[string]string{"x": "y"}, Value: 4, Time: ts1},
		},
		Period: 60,
		Time:   ts2,
		Tags:   &map[string]string{"j": "k", "l": "m"},
	}

	// post batch
	mockClient := &mockHTTPClient{t: t}
	baseURL := "http://unused:8383/v1/" // XXX remove /v1/ from baseURL
	c := NewClient("abcdef", BaseURLClientOption(baseURL))
	c.httpClient = mockClient
	ms := &MeasurementsService{client: c}
	resp, err := ms.Create(testBatch)
	assert.Nil(t, resp)
	assert.Equal(t, errMockHTTPClientExpected, err)

	// assert expected http.Request
	require.Len(t, mockClient.reqs, 1)
	req := mockClient.reqs[0]
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "unused:8383", req.URL.Host)
	assert.Equal(t, "/v1/measurements", req.URL.Path)

	// unmarshal JSON and assert
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(req.Body)
	assert.NoError(t, err)
	data := make(map[string]interface{})
	err = json.Unmarshal(buf.Bytes(), &data)
	assert.NoError(t, err)

	assert.EqualValues(t, 60, data["period"])
	assert.EqualValues(t, ts2, data["time"])
	assert.Equal(t, map[string]interface{}{"j": "k", "l": "m"}, data["tags"])

	mts, ok := data["measurements"].([]interface{})
	require.True(t, ok)
	require.Len(t, mts, len(testBatch.Measurements))

	m0, _ := mts[0].(map[string]interface{})
	assert.Equal(t, "metric1", m0["name"])
	assert.Equal(t, map[string]interface{}{"a": "b"}, m0["tags"])
	assert.EqualValues(t, 5, m0["value"])
	assert.NotContains(t, m0, "time")
	assert.NotContains(t, m0, "sum")
	assert.NotContains(t, m0, "count")
	assert.NotContains(t, m0, "min")
	assert.NotContains(t, m0, "max")

	m1, _ := mts[1].(map[string]interface{})
	assert.Equal(t, "metric2", m1["name"])
	assert.Equal(t, map[string]interface{}{"c": "d"}, m1["tags"])
	assert.EqualValues(t, 10, m1["sum"])
	assert.EqualValues(t, 2, m1["count"])
	assert.NotContains(t, m1, "time")
	assert.NotContains(t, m1, "value")

	m2, _ := mts[2].(map[string]interface{})
	assert.Equal(t, "metric3", m2["name"])
	assert.Equal(t, map[string]interface{}{"x": "y"}, m2["tags"])
	assert.EqualValues(t, 4, m2["value"])
	assert.EqualValues(t, ts1, m2["time"])
	assert.NotContains(t, m2, "sum")
	assert.NotContains(t, m2, "count")
	assert.NotContains(t, m2, "min")
	assert.NotContains(t, m2, "max")
}
