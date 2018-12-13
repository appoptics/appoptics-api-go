package appoptics_test

import "net/http"

func CreateSnapshotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "href": "https://api.appoptics.com/v1/snapshots/1",
  "job_href": "https://api.appoptics.com/v1/jobs/123456",
  "image_href": "http://snapshots.appoptics.com/chart/tuqlgn1i-71569.png",
  "duration": 3600,
  "end_time": "2016-02-20T01:18:46Z",
  "created_at": "2016-02-20T01:18:46Z",
  "updated_at": "2016-02-20T01:18:46Z",
  "subject": {
    "chart": {
      "id": 1,
      "sources": [
        "*"
      ],
      "type": "stacked"
    }
  }
}`
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responseBody))
	}
}

func RetrieveSnapshotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "href": "https://api.appoptics.com/v1/snapshots/1",
  "job_href": "https://api.appoptics.com/v1/jobs/123456",
  "image_href": "http://snapshots.appoptics.com/chart/tuqlgn1i-71569.png",
  "duration": 3600,
  "end_time": "2016-02-20T01:18:46Z",
  "created_at": "2016-02-20T01:18:46Z",
  "updated_at": "2016-02-20T01:18:46Z",
  "subject": {
    "chart": {
      "id": 1,
      "sources": [
        "*"
      ],
      "type": "stacked"
    }
  }
}}`
		w.Write([]byte(responseBody))
	}
}
