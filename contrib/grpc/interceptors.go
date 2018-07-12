package interceptors

import (
	"fmt"
	"path"
	"time"

	"github.com/appoptics/appoptics-api-go"
	ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type InstrumentedServer struct {
	m       *appoptics.MeasurementSet
	service string
	method  string
	tags    map[string]interface{}
}

func NewInstrumentedServer(m *appoptics.MeasurementSet, service string, method string, tags map[string]interface{}) *InstrumentedServer {
	if tags == nil {
		tags = make(map[string]interface{})
	}
	return &InstrumentedServer{
		m:       m,
		service: service,
		method:  method,
		tags:    tags,
	}
}

func (s *InstrumentedServer) key(key string) string {
	return appoptics.MetricWithTags(fmt.Sprintf("%s.%s.%s", s.service, s.method, key), s.tags)
}

func (s *InstrumentedServer) received() {
	s.m.Incr(s.key("received"))
}

func (s *InstrumentedServer) handled(err error) {
	s.tags["status"] = status.Code(err).String()
	s.m.Incr(s.key("result"))
}

func (s *InstrumentedServer) timed(t time.Duration) {
	s.m.UpdateAggregatorValue(s.key("time_ms"), float64(t/time.Millisecond)+float64(t%time.Millisecond)/1e9)
}

// InstrumentedServerStream implements gRPC's `Stream` interface, providing metrics for
// the number of times Send- and RecvMsg() are called. Timing metrics are not emitted
type InstrumentedServerStream struct {
	grpc.ServerStream
	*InstrumentedServer
}

func (s *InstrumentedServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	s.handled(err)
	return err
}

func (s *InstrumentedServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	s.received()
	return err
}

// Creates a UnaryServerInterceptor that submits AO metrics using the given MeasurementSet. Emits
// counts of requests received, counts of requests handled (tagged by status code) and timings
func UnaryServerInterceptor(m *appoptics.MeasurementSet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		instrument := NewInstrumentedServer(
			m,
			path.Dir(info.FullMethod)[1:],
			path.Base(info.FullMethod),
			ctxtags.Extract(ctx).Values(),
		)

		instrument.received()

		start := time.Now()
		resp, err := handler(ctx, req)
		instrument.timed(time.Now().Sub(start))
		instrument.handled(err)
		return resp, err
	}
}

// Creates a StreamServerInterceptor that submits AO metrics using the given MeasurementSet. See
// InstrumentedServerStream for a list of metrics emitted
func StreamServerInterceptor(m *appoptics.MeasurementSet) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		instrument := &InstrumentedServerStream{
			ss,
			NewInstrumentedServer(
				m,
				path.Dir(info.FullMethod)[1:],
				path.Base(info.FullMethod),
				nil,
			),
		}

		start := time.Now()
		err := handler(srv, instrument)
		instrument.timed(time.Now().Sub(start))
		instrument.handled(err)

		return err
	}
}
