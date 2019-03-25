package appoptics_test

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
)

func TestApiTokensService_List(t *testing.T) {
	apiTokenResponse, err := client.ApiTokensService().List()
	assert.Nil(t, err)
	query := apiTokenResponse.Query
	firstToken := apiTokenResponse.ApiTokens[0]

	// Query
	assert.Equal(t, 3, query.Found)
	assert.Equal(t, 3, query.Length)
	assert.Equal(t, 3, query.Total)

	// ApiToken
	assert.Equal(t, "Boss Token", *firstToken.Name)
	assert.Equal(t, "9b593b3dff7cef442268c3056625981b92d5feb622cfe23fef8d716d5eecd2c3", *firstToken.Token)
	assert.Equal(t, true, *firstToken.Active)
	assert.Equal(t, "admin", *firstToken.Role)
}

func TestApiTokensService_Create(t *testing.T) {
	apiToken, err := client.ApiTokensService().Create(&appoptics.ApiToken{})
	assert.Nil(t, err)
	assert.Equal(t, "My New Token", *apiToken.Name)
	assert.Equal(t, "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf", *apiToken.Token)
	assert.Equal(t, true, *apiToken.Active)
	assert.Equal(t, "admin", *apiToken.Role)
}

func TestApiTokensService_Retrieve(t *testing.T) {
	apiTokenResponse, err := client.ApiTokensService().Retrieve("foobar")
	assert.Nil(t, err)

	query := apiTokenResponse.Query
	firstToken := apiTokenResponse.ApiTokens[0]

	// Query
	assert.Equal(t, 2, query.Found)
	assert.Equal(t, 2, query.Length)
	assert.Equal(t, 2, query.Total)

	// ApiToken
	assert.Equal(t, "Token for collectors", *firstToken.Name)
	assert.Equal(t, "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf", *firstToken.Token)
	assert.Equal(t, true, *firstToken.Active)
	assert.Equal(t, "recorder", *firstToken.Role)

	if err != nil {
		t.Errorf("error running List: %v", err)
	}
}

func TestApiTokensService_Update(t *testing.T) {
	apiToken, err := client.ApiTokensService().Update(&appoptics.ApiToken{})
	assert.Nil(t, err)
	assert.Equal(t, "New Token Name", *apiToken.Name)
	assert.Equal(t, "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf", *apiToken.Token)
	assert.Equal(t, false, *apiToken.Active)
	assert.Equal(t, "admin", *apiToken.Role)
}
