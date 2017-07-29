package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
)

type ContextKey int

const (
	AuthContextID ContextKey = iota
)

func (id ContextKey) ToInt() int {
	return int(id)
}

type Service interface {
	Authenticate(ctx context.Context, request *http.Request) context.Context
	Login(req *http.Request) (string, error)
}

type Authenticatable interface {
	Authenticate(req http.Request) (string, error)
}

type service struct {
	session    session.Service
	authable   Authenticatable
	cookieName string
}

func (s *service) Authenticate(ctx context.Context, request *http.Request) context.Context {
	cookie, err := request.Cookie("session")

	if err != nil {
		return ctx
	}

	sess, err := s.session.Get(cookie.Value)

	if err != nil {
		return ctx
	}

	spew.Dump(sess)

	authctx := context.WithValue(ctx, AuthContextID, sess.UID)

	return authctx
}

func (s *service) Login(req *http.Request) (string, error) {
	id, err := s.authable.Authenticate(*req)

	if err != nil {
		return "", err
	}

	token, err := s.session.Make(session.SessionUID(id), time.Date(
		time.Now().Year()+1,
		time.January,
		0, 0, 0, 0, 0,
		time.FixedZone("Europe/Oslo", 0)),
	)

	if err != nil {
		return "", err
	}

	return token.ToString(), nil
}

func NewService(session session.Service, authable Authenticatable) Service {
	return &service{
		session:  session,
		authable: authable,
	}
}
