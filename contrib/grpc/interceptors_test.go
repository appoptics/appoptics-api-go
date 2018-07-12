package interceptors

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"github.com/appoptics/appoptics-api-go"
	"github.com/stretchr/testify/assert"
	"fmt"
	"strings"
)

var (
	srv = "server_full_of_testing"
	uInfo = &grpc.UnaryServerInfo{
		FullMethod: "/something/blah",
		Server:     srv,
	}
)

func TestUnaryRequest(t *testing.T) {
	measures := appoptics.NewMeasurementSet()
	intercept := UnaryServerInterceptor(measures)
	ctx := context.Background()
	var handler grpc.UnaryHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
		return "something", nil
	}
	intercept(ctx, "some data", uInfo, handler)

	stupidTestShenanigans := fmt.Sprintf("%v", measures)
	assert.NotNil(t, strings.Contains(stupidTestShenanigans, "something.blah.received"))
	assert.NotNil(t, strings.Contains(stupidTestShenanigans, "something.blah.result::status::OK"))
	assert.NotNil(t, strings.Contains(stupidTestShenanigans, "something.blah.time_ms"))
}
