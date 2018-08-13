package appoptics_test

import (
	"net/http"
)

func ListApiTokensHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "query": {
    "found": 3,
    "length": 3,
    "offset": 0,
    "total": 3
  },
  "api_tokens": [
    {
      "name": "Boss Token",
      "token": "9b593b3dff7cef442268c3056625981b92d5feb622cfe23fef8d716d5eecd2c3",
      "active": true,
      "role": "admin",
      "href": "http://api.appoptics.com/v1/api_tokens/2",
      "created_at": "2013-01-22 18:08:15 UTC",
      "updated_at": "2013-01-22 18:08:15 UTC"
    },
    {
      "name": "Token for collectors",
      "token": "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf",
      "active": true,
      "role": "recorder",
      "href": "http://api.appoptics.com/v1/api_tokens/28",
      "created_at": "2013-02-01 18:53:38 UTC",
      "updated_at": "2013-02-01 18:53:38 UTC"
    },
    {
      "name": "Token that has been disabled",
      "active": false,
      "role": "viewer",
      "href": "http://api.appoptics.com/v1/api_tokens/29",
      "token": "d1ffbbbe327c6839a71023c2a8c9ba921207e32e427a49fb221843d74d63f7b8",
      "created_at": "2013-02-01 18:54:28 UTC",
      "updated_at": "2013-02-01 18:54:28 UTC"
    }
  ]
}`
		w.Write([]byte(responseBody))
	}
}

func RetrieveApiTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "api_tokens": [
    {
      "name": "Token for collectors",
      "token": "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf",
      "active": true,
      "role": "recorder",
      "href": "http://api.appoptics.com/v1/api_tokens/28",
      "created_at": "2013-02-01 18:53:38 UTC",
      "updated_at": "2013-02-01 18:53:38 UTC"
    },
    {
      "name": "Token that has been disabled",
      "active": false,
      "role": "viewer",
      "href": "http://api.appoptics.com/v1/api_tokens/29",
      "token": "d1ffbbbe327c6839a71023c2a8c9ba921207e32e427a49fb221843d74d63f7b8",
      "created_at": "2013-02-01 18:54:28 UTC",
      "updated_at": "2013-02-01 18:54:28 UTC"
    }
  ],
    "query": {
        "found": 2,
        "length": 2,
        "offset": 0,
        "total": 2
    }
}`
		w.Write([]byte(responseBody))
	}
}

func CreateApiTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "My New Token",
  "token": "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf",
  "active": true,
  "role": "admin",
  "href": "http://api.appoptics.dev/v1/api_tokens/28",
  "created_at": "2013-02-01 18:53:38 UTC",
  "updated_at": "2013-02-01 18:53:38 UTC"
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func UpdateApiTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "New Token Name",
  "token": "24f9fb2134399595b91da1dcac39cb6eafc68a07fa08ad3d70892b7aad10e1cf",
  "active": false,
  "role": "admin",
  "href": "http://api.appoptics.dev/v1/api_tokens/28",
  "created_at": "2013-02-01 18:53:38 UTC",
  "updated_at": "2013-02-01 19:51:22 UTC"
}`
		w.Write([]byte(responseBody))
	}
}

func DeleteApiTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
