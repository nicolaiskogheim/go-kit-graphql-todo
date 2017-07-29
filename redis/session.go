package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
)

type sessionRepository struct {
	client redis.Client
}

func (r *sessionRepository) Find(uid session.SessionToken) (*session.Session, error) {
	return nil, nil
}

func (r *sessionRepository) Store(s *session.Session) error {
	return nil
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

	return &sessionRepository{
		client: client,
	}
}
