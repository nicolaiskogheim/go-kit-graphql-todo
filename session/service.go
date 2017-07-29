package session

import (
	"time"
)

type Service interface {
	Make(uid SessionUID, expires time.Time) (*SessionToken, error)
	Get(token SessionToken) (*Session, error)
}

type service struct {
	repository SessionRepository
}

func NewService(repository SessionRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Make(uid SessionUID, expires time.Time) (*SessionToken, error) {
	session := New(uid, expires)
	err := s.repository.Store(session)

	if err != nil {
		return nil, err
	}

	return &session.Token, nil
}

func (s *service) Get(token SessionToken) (*Session, error) {
	session, err := s.repository.Find(SessionToken(token))

	if err != nil {
		return nil, err
	}

	return session, nil
}

type SessionRepository interface {
	Find(token SessionToken) (*Session, error)
	Store(s *Session) error
}
