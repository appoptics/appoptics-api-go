package appoptics_test

import "net/http"

func ListMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
}`
		w.Write([]byte(responseBody))
	}

}

func CreateMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "cpu.percent.used",
  "display_name": null,
  "type": "composite",
  "description": null,
  "period": null,
  "source_lag": null,
  "composite": "s(\"cpu.percent.user\", {\"environment\" : \"prod\", \"service\": \"api\"})"
}`
		w.Write([]byte(responseBody))
	}

}

func UpdateMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}
}

func RetrieveMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "name": "cpu_temp",
  "display_name": "cpu_temp",
  "description": "Current CPU temperature in Fahrenheit",
  "period": 60,
  "type": "gauge",
  "attributes": {
    "created_by_ua": "appoptics-metrics/0.7.4 (ruby; 1.9.3p194; x86_64-linux) direct-faraday/0.8.4",
    "display_max": null,
    "display_min": 0,
    "display_stacked": true,
    "display_units_long": "Fahrenheit",
    "display_units_short": "°F"
  }
}`
		w.Write([]byte(responseBody))
	}
}

func DeleteMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}
}
