package todo

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{
		logger:  logger,
		Service: s,
	}
}

func (s *loggingService) Add(t *Todo) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add",
			"id", t.ID,
			"text", t.Text,
			"done", t.Done,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())

	return s.Service.Add(t)
}

func (s *loggingService) Toggle(user user.User, id TodoID) (t *Todo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "toggle",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())

	return s.Service.Toggle(user, id)
}

func (s *loggingService) Remove(id TodoID) (t *Todo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "remove",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())

	return s.Service.Remove(id)
}

func (s *loggingService) Find(id TodoID) (t *Todo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())

	return s.Service.Find(id)
}

func (s *loggingService) FindAll() []*Todo {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find_all",
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.FindAll()
}
