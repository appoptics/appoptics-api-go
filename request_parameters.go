package appoptics

import (
	"net/http"
	"strconv"
)

// RequestParameters holds pagination values
// https://docs.appoptics.com/api/?shell#request-parameters
type RequestParameters struct {
	Offset int
	Length int
	Sort   string
}

// AddToRequest mutates the provided http.Request with the RequestParameters values
// Note that only valid values for Sort are "asc" and "desc" but the client does not enforce this.
func (rp *RequestParameters) AddToRequest(req *http.Request) {
	if rp == nil {
		return
	}
	values := req.URL.Query()
	values.Add("offset", strconv.Itoa(rp.Offset))
	values.Add("length", strconv.Itoa(rp.Length))
	values.Add("sort", rp.Sort)

	req.URL.RawQuery = values.Encode()
}
