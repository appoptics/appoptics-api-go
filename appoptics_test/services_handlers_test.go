package appoptics_test

import "net/http"

func ListServicesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "query": {
    "found": 2,
    "length": 2,
    "offset": 0,
    "total": 2
  },
  "services": [
    {
      "id": 145,
      "type": "slack",
      "settings": {
        "room": "Ops",
        "token": "1234567890ABCDEF",
        "subdomain": "acme"
      },
      "title": "Notify Ops Room"
    },
    {
      "id": 156,
      "type": "mail",
      "settings": {
        "addresses": "george@example.com,fred@example.com"
      },
      "title": "Email ops team"
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func RetrieveServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 156,
  "type": "mail",
  "settings": {
    "addresses": "george@example.com,fred@example.com"
  },
  "title": "Email ops team"
}`
		w.Write([]byte(responseBody))
	}
}

func CreateServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 145,
  "type": "campfire",
  "settings": {
    "room": "Ops",
    "token": "1234567890ABCDEF",
    "subdomain": "acme"
  },
  "title": "Notify Ops Room"
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func UpdateServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
