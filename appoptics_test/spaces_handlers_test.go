package appoptics_test

import (
	"net/http"
)

func ListSpacesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "query": {
    "found": 1,
    "length": 1,
    "offset": 0,
    "total": 15
  },
  "spaces": [
    {
      "id": 4,
      "name": "staging_ops"
    }
  ]
}`

		w.Write([]byte(responseBody))
	}
}

func RetrieveSpaceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "CPUs",
  "id": 129,
  "charts": [
    {
      "id": 915
    },
    {
      "id": 1321
    },
    {
      "id": 47842
    },
    {
      "id": 922
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func CreateSpaceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 129,
  "name": "CPUs"
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func UpdateSpaceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "MEMORY"
}`

		w.Write([]byte(responseBody))
	}
}

func DeleteSpaceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
