package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
)

type sessionRepository struct {
	client redis.Client
}

func (r *sessionRepository) Find(uid session.SessionUID) (*session.Session, error) {
	val, err := r.client.Get(uid.ToString()).Result()

	if err != nil {
		return nil, err
	}

	s := session.Session{UID: uid,
		Token:   session.SessionToken(val),
		Expires: time.Now(),
	}

	return &s, nil
}

func (r *sessionRepository) Store(s *session.Session) error {
	expires := s.Expires.Sub(time.Now())
	err := r.client.Set(s.Token.ToString(), s.UID.ToString(), expires).Err()

	return err
}

func NewSessionRepository(client redis.Client) session.SessionRepository {

	// XXX(nicolai): This is how we check that the client is alive
	// We may want to do this in regular intervals and alert something
	// when/if the client dies. Or maybe there are recovery strategies
	// for this. If this app is orchestrated by Kubernetes, then
	// Kubernetes could be responsible for keeping redis alive.
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	r := &sessionRepository{
		client: client,
	}

	return r
}
