package appoptics_test

import "net/http"

func ListAnnotationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "query": {
    "found": 2,
    "length": 2,
    "offset": 0,
    "total": 2
  },
  "annotations": [
    {
      "name": "api-deploys",
      "display_name": "Deploys to API"
    },
    {
      "name": "app-deploys",
      "display_name": "Deploys to UX app"
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func RetrieveAnnotationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "api-deploys",
  "display_name": "Deploys to API",
  "events": [
    {
      "unassigned": [
        {
          "id": 45,
          "title": "Deployed v91",
          "description": null,
          "start_time": 1234567890,
          "end_time": null,
          "source": null,
          "links": [
            {
              "label": "Github commit",
              "rel": "github",
              "href": "http://github.com/org/app/commit/01beaf"
            },
            {
              "label": "Jenkins CI build job",
              "rel": "jenkins",
              "href": "http://ci.acme.com/job/api/34"
            }
          ]
        }
      ],
      "foo3.bar.com": [
        {
          "id": 123,
          "title": "My Annotation",
          "description": null,
          "start_time": 1234567890,
          "end_time": null,
          "source": "foo3.bar.com",
          "links": [

          ]
        },
        {
          "id": 1098,
          "title": "Chef run finished",
          "description": "Chef migrated up to SHA 44a87f9",
          "start_time": 1234567908,
          "end_time": 1234567958,
          "source": "foo3.bar.com",
          "links": [

          ]
        }
      ]
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func RetrieveAnnotationEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 189,
  "title": "Deployed v91",
  "description": null,
  "start_time": 1234567890,
  "end_time": null,
  "source": null,
  "links": [
    {
      "label": "Github commit",
      "rel": "github",
      "href": "http://github.com/org/app/commit/01beaf"
    },
    {
      "label": "Jenkins CI build job",
      "rel": "jenkins",
      "href": "http://ci.acme.com/job/api/34"
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func CreateAnnotationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 123,
  "title": "My Annotation",
  "description": "Joe deployed v29 to metrics",
  "source": null,
  "start_time": 1234567890,
  "end_time": null,
  "links": [

  ]
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func UpdateAnnotationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "rel": "github",
  "label": "Github Commit",
  "href": "https://github.com/acme/app/commits/01beaf"
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func DeleteAnnotationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
