package appoptics_test

import "net/http"

func ListAlertsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
	"query": {
		"offset": 0,
		"length": 1,
		"found": 1,
		"total": 7
	},
	"alerts": [{
		"id": 1400310,
		"name": "CPU.utilization",
		"description": "How much power the box got",
		"conditions": [{
			"id": 1016,
			"type": "above",
			"metric_name": "AWS.EC2.CPUUtilization",
			"source": "*prod*",
			"threshold": 90,
			"duration": 300,
			"summary_function": "max"
		}],
		"services": [{
			"id": 1153,
			"type": "mail",
			"settings": {
				"addresses": "foo@domain.com,bar@domain.com"
			},
			"title": "Ops Team"
		}],
		"attributes": {},
		"active": true,
		"created_at": 1394745670,
		"updated_at": 1394745670,
		"version": 2,
		"rearm_seconds": 600,
		"rearm_per_signal": false
	}]
}`
		w.Write([]byte(responseBody))
	}
}

func RetrieveAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 123,
  "name": "production.web.frontend.response_time",
  "description":"Web Response Time",
  "conditions":[
      {
         "id":19375969,
         "type":"above",
         "metric_name":"web.nginx.response_time",
         "source":null,
         "threshold":200.0,
         "summary_function":"average",
         "tags":[
            {
               "name":"environment",
               "grouped":false,
               "values":[
                  "production"
               ]
            }
         ]
      }
   ],
  "services":[
      {
         "id":17584,
         "type":"slack",
         "settings":{
            "url":"https://hooks.slack.com/services/XYZABC/a1b2c3/asdf"
         },
         "title":"appoptics-services"
      }
   ],
  "attributes": {
    "runbook_url": "http://myco.com/runbooks/response_time"
  },
  "active":true,
  "created_at":1484588756,
  "updated_at":1484588756,
  "version":2,
  "rearm_seconds":600,
  "rearm_per_signal":false,
  "md":true
}`
		w.Write([]byte(responseBody))
	}
}

func CreateAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
   "id":1234567,
   "name":"production.web.frontend.response_time",
   "description":"Web Response Time",
   "conditions":[
      {
         "id":19376030,
         "type":"above",
         "metric_name":"web.nginx.response_time",
         "threshold":200.0,
         "summary_function":"max",
         "tags":[
            {
               "name":"tag_name",
               "grouped":false,
               "values":[
                  "tag_value"
               ]
            }
         ]
      }
   ],
   "services":[
      {
         "id":17584,
         "type":"slack",
         "settings":{
            "url":"https://hooks.slack.com/services/ABCDEFG/A1B2C3/asdfg1234"
         },
         "title":"librato-services"
      }
   ],
   "attributes":{
      "runbook_url":"http://myco.com/runbooks/response_time"
   },
   "active":true,
   "created_at":1484594787,
   "updated_at":1484594787,
   "version":2,
   "rearm_seconds":600,
   "rearm_per_signal":false,
   "md":true
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func StatusAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
    "alert": {
        "id": 120
    },
    "status": "triggered"
}`
		w.Write([]byte(responseBody))
	}
}

func UpdateAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
func DisassociateAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
}
func AssociateAlertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
}
