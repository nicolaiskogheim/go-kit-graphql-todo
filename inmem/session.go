package inmem

import (
	"sync"
	"time"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
)

type sessionRepository struct {
	mtx      sync.RWMutex
	sessions map[session.SessionToken]session.SessionUID
}

func (r *sessionRepository) Store(s *session.Session) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.sessions[s.Token] = s.UID

	return nil
}

func (r *sessionRepository) Find(token session.SessionToken) (*session.Session, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.sessions[token]; ok {
		s := &session.Session{
			UID:     val,
			Token:   token,
			Expires: time.Now(),
		}

		return s, nil
	}

	return nil, session.ErrUnknown
}

func NewSessionRepository() session.SessionRepository {
	r := &sessionRepository{
		sessions: make(map[session.SessionToken]session.SessionUID),
	}

	r.sessions[session.Session1.Token] = session.Session1.UID

	return r
}
