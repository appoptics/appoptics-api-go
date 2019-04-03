package live_tests

import (
	"log"
	"testing"

	"os"

	"github.com/appoptics/appoptics-api-go"
)

const (
	TestPrefix        = "iNtEgRaTiOnTEST"
	CreatedNameString = "created"
	UpdatedNameString = "updated"
)

var (
	client *appoptics.Client
)

func TestMain(m *testing.M) {
	token := os.Getenv("APPOPTICS_TOKEN")
	if token == "" {
		log.Fatal("set APPOPTICS_TOKEN in the environment")
	}

	if debugMode := os.Getenv("AO_CLIENT_DEBUG"); debugMode != "" {
		client = appoptics.NewClient(token, appoptics.SetDebugMode())
	} else {
		client = appoptics.NewClient(token)
	}

	os.Exit(m.Run())
}
