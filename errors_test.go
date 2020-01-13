package appoptics

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	ErrorRequestBodyUnauthorized = `{
	"errors":{
		"request":["This token is not permitted to perform this action on this resource"]
	}
}
`
)

func TestErrorResponse_Error(t *testing.T) {
	errResp := &ErrorResponse{}
	body := strings.NewReader(ErrorRequestBodyUnauthorized)
	decoder := json.NewDecoder(body)
	decodeErr := decoder.Decode(errResp)

	t.Run("it decodes without error", func(t *testing.T) {
		require.NoError(t, decodeErr)
	})

	t.Run("it holds detailed error information", func(t *testing.T) {
		v := errResp.Errors.(map[string]interface{})
		mapV := v["request"].([]interface{})
		sVal := mapV[0].(string)
		assert.Equal(t, "This token is not permitted to perform this action on this resource", sVal)
	})
	
	t.Run("it places error information in Error() output", func(t *testing.T) {
		errResp.Status = "403 Forbidden"
		actual := `403 Forbidden - {"request":["This token is not permitted to perform this action on this resource"]}`
		assert.Equal(t, errResp.Error(), actual)
	} )
}


