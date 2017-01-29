package graphql

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/graphql-go/graphql"
	"golang.org/x/net/context"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) Do(ctx context.Context, request interface{}) *graphql.Result {
	defer func(begin time.Time) {
		s.requestCount.With("method", "do").Add(1)
		s.requestLatency.With("method", "do").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Do(ctx, request)
}
