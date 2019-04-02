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

	client = appoptics.NewClient(token)
	// Uncomment the below to see response status/body stdout while tests run
	//client = appoptics.NewClient(token, appoptics.SetDebugMode())
	os.Exit(m.Run())
}
