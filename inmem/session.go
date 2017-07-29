package inmem

import (
	"sync"
	"time"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
)

type sessionRepository struct {
	mtx      sync.RWMutex
	sessions map[session.SessionUID]session.SessionToken
}

func (r *sessionRepository) Store(s *session.Session) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.sessions[s.UID] = s.Token

	return nil
}

func (r *sessionRepository) Find(uid session.SessionUID) (*session.Session, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.sessions[uid]; ok {
		s := &session.Session{
			UID:     uid,
			Token:   val,
			Expires: time.Now(),
		}

		return s, nil
	}

	return nil, session.ErrUnknown
}

func NewSessionRepository() session.SessionRepository {
	r := &sessionRepository{
		sessions: make(map[session.SessionUID]session.SessionToken),
	}

	r.sessions[session.Session1.UID] = session.Session1.Token

	return r
}
