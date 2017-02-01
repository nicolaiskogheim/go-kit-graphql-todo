package user

import (
	"time"

	"github.com/go-kit/kit/log"
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

func (s *loggingService) Add(u *User) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add",
			"id", u.ID,
			"name", u.Name,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())

	return s.Service.Add(u)
}

func (s *loggingService) Remove(id UserID) (err error) {
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

func (s *loggingService) Find(id UserID) (u *User, err error) {
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

func (s *loggingService) FindAll() []*User {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find_all",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.FindAll()
}
