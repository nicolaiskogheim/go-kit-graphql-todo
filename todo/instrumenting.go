package todo

import (
	"time"

	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) Add(t *Todo) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add").Add(1)
		s.requestLatency.With("method", "add").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Add(t)
}

func (s *instrumentingService) Toggle(id TodoID) (*Todo, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "toggle").Add(1)
		s.requestLatency.With("method", "toggle").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Toggle(id)
}

func (s *instrumentingService) Remove(id TodoID) (*Todo, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "remove").Add(1)
		s.requestLatency.With("method", "remove").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Remove(id)
}

func (s *instrumentingService) FindAll() []*Todo {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find_all").Add(1)
		s.requestLatency.With("method", "find_all").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindAll()
}

func (s *instrumentingService) Find(id TodoID) (*Todo, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find").Add(1)
		s.requestLatency.With("method", "find").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Find(id)
}
