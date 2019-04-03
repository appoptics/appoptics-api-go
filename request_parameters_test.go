package appoptics

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRequestParameters_AddToRequest_WithParams(t *testing.T) {
	baseUrl := "http://example.com"
	orderby := "name"
	sort := "desc"
	req, _ := http.NewRequest("GET", baseUrl, nil)

	rp := PaginationParameters{
		Offset:  10,
		Length:  20,
		Orderby: orderby,
		Sort:    sort,
	}

	rp.AddToRequest(req)

	fullUrl := fmt.Sprintf("%s?length=%d&offset=%d&orderby=%s&sort=%s",
		baseUrl, rp.Length, rp.Offset, rp.Orderby, rp.Sort)

	assert.Equal(t, req.URL.String(), fullUrl)
}
