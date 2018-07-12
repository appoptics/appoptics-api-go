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

type InstrumentedServerStream struct {
	grpc.ServerStream
	InstrumentedServer
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

func UnaryServerInterceptor(m *appoptics.MeasurementSet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		instrument := InstrumentedServer{
			m,
			path.Dir(info.FullMethod)[1:],
			path.Base(info.FullMethod),
			ctxtags.Extract(ctx).Values(),
		}

		instrument.received()

		start := time.Now()
		resp, err := handler(ctx, req)
		instrument.timed(time.Now().Sub(start))
		instrument.handled(err)
		return resp, err
	}
}

func StreamServerInterceptor(m *appoptics.MeasurementSet) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		instrument := &InstrumentedServerStream{
			ss,
			InstrumentedServer{
				m,
				path.Dir(info.FullMethod)[1:],
				path.Base(info.FullMethod),
				nil,
			},
		}

		start := time.Now()
		err := handler(srv, instrument)
		instrument.timed(time.Now().Sub(start))
		instrument.handled(err)

		return err
	}
}
