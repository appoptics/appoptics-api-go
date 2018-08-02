package appoptics_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/gorilla/mux"
)

var (
	client *appoptics.Client
	server *httptest.Server
)

func setup() {
	router := NewServerTestMux()
	server = httptest.NewServer(router)
	serverURLWithVersion := fmt.Sprintf("%s/v1/", server.URL)
	client = appoptics.NewClient("deadbeef", appoptics.BaseURLClientOption(serverURLWithVersion))
}

func teardown() {
	server.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func NewServerTestMux() *mux.Router {
	router := mux.NewRouter()

	// Measurements

	// Metrics

	// Spaces
	router.Handle("/v1/spaces", ListSpacesHandler()).Methods("GET")
	router.Handle("/v1/spaces", CreateSpaceHandler()).Methods("POST")
	router.Handle("/v1/spaces/{id}", RetrieveSpaceHandler()).Methods("GET")
	router.Handle("/v1/spaces/{id}", UpdateSpaceHandler()).Methods("PUT")
	router.Handle("/v1/spaces/{id}", DeleteSpaceHandler()).Methods("DELETE")

	// Charts
	router.Handle("/v1/spaces/{spaceId}/charts", ListChartsHandler()).Methods("GET")
	router.Handle("/v1/spaces/{spaceId}/charts", CreateChartHandler()).Methods("POST")
	router.Handle("/v1/spaces/{spaceId}/charts/{chartId}", RetrieveChartHandler()).Methods("GET")
	router.Handle("/v1/spaces/{spaceId}/charts/{chartId}", UpdateChartHandler()).Methods("PUT")
	router.Handle("/v1/spaces/{spaceId}/charts/{chartId}", DeleteChartHandler()).Methods("DELETE")

	// Annotations

	// Alerts

	// API Tokens

	// Jobs

	// Snapshots

	return router
}
