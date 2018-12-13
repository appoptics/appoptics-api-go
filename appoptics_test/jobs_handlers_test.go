package appoptics_test

import "net/http"

func RetrieveJobsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 123456,
  "state": "failed",
  "progress": 76.5,
  "errors": {
    "name": [
      "is invalid"
    ]
  }
}}`
		w.Write([]byte(responseBody))
	}
}
