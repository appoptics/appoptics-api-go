package main

import (
	"testing"

	"github.com/appoptics/appoptics-api-go"
)

func TestPrintVersion(t *testing.T) {
	if printVersion() != appoptics.Version() {
		t.Errorf("tag version doesn't match declared version")
	}
}
