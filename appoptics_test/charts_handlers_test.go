package appoptics_test

import "net/http"

func ListChartsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `[
  {
    "id": 1234567,
    "name": "CPU Usage",
    "type": "line",
    "streams": [
      {
        "id": 27035309,
        "metric": "cpu.percent.idle",
        "type": "gauge",
        "tags": [
          {
            "name": "environment",
            "values": [
              "*"
            ]
          }
        ]
      },
      {
        "id": 27035310,
        "metric": "cpu.percent.user",
        "type": "gauge",
        "tags": [
          {
            "name": "environment",
            "values": [
              "prod"
            ]
          }
        ]
      }
    ],
    "thresholds": null
  }
]`
		w.Write([]byte(responseBody))
	}
}

func RetrieveChartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 3700969,
  "name": "CPU Usage",
  "type": "line",
  "streams": [
    {
      "id": 27003258,
      "metric": "cpu.percent.idle",
      "type": "gauge",
      "tags": [
        {
          "name": "region",
          "values": [
            "us-east-1"
          ]
        }
      ]
    }
  ],
  "thresholds": null
}`
		w.Write([]byte(responseBody))
	}
}

func CreateChartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 1234567,
  "name": "CPU Usage",
  "type": "line",
  "streams": [
    {
      "id": 27032885,
      "metric": "cpu.percent.idle",
      "type": "gauge",
      "tags": [
        {
          "name": "environment",
          "values": [
            "*"
          ]
        }
      ]
    },
    {
      "id": 27032886,
      "metric": "cpu.percent.user",
      "type": "gauge",
      "tags": [
        {
          "name": "environment",
          "values": [
            "prod"
          ]
        }
      ]
    }
  ],
  "thresholds": null
}`
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseBody))
	}
}

func UpdateChartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{
  "id": 3700969,
  "name": "Temperature",
  "type": "line",
  "streams": [
    {
      "id": 27003258,
      "metric": "collectd.cpu.0.cpu.user",
      "type": "gauge",
      "tags": [
        {
          "name": "region",
          "values": [
            "us-east-1"
          ]
        }
      ]
    }
  ],
  "thresholds": null
}`
		w.Write([]byte(responseBody))
	}
}

func DeleteChartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
