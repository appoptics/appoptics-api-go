package integration_tests

import (
	"log"
	"testing"

	"os"

	"github.com/appoptics/appoptics-api-go"
)

const (
	TestPrefix = "iNtEgRaTiOnTEST"
)

var (
	client *appoptics.Client
)

func TestMain(m *testing.M) {
	token := os.Getenv("APPOPTICS_TOKEN")
	if token == "" {
		log.Fatal("set APPOPTICS_TOKEN in the environment")
	}
	//client = appoptics.NewClient(token)
	client = appoptics.NewClient(token, appoptics.SetDebugMode())

	os.Exit(m.Run())
}
