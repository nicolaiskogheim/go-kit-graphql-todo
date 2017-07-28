package inmem

import "sync"
import "github.com/nicolaiskogheim/go-kit-graphql-todo/session"

type sessionRepository struct {
	mtx      sync.RWMutex
	sessions map[session.SessionToken]*session.Session
}

func (r *sessionRepository) Store(s *session.Session) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.sessions[s.Token] = s

	return nil
}

func (r *sessionRepository) Find(token session.SessionToken) (*session.Session, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.sessions[token]; ok {
		return val, nil
	}

	return nil, session.ErrUnknown
}

func NewSessionRepository() session.SessionRepository {
	r := &sessionRepository{
		sessions: make(map[session.SessionToken]*session.Session),
	}

	return r
}
