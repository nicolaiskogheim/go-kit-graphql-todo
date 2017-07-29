package auth

import (
	"context"
	"net/http"
	"time"

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
	Login(req *http.Request) (*session.SessionToken, error)
}

type Authenticatable interface {
	Authenticate(req http.Request) (*Identifier, error)
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

	sess, err := s.session.Get(session.SessionToken(cookie.Value))

	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, AuthContextID, sess.UID)
}

func (s *service) Login(req *http.Request) (*session.SessionToken, error) {
	id, err := s.authable.Authenticate(*req)

	if err != nil || id == nil {
		return nil, err
	}

	token, err := s.session.Make(session.SessionUID(id.ToString()), time.Date(
		time.Now().Year()+1,
		time.January,
		0, 0, 0, 0, 0,
		time.FixedZone("Europe/Oslo", 0)),
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// Viewer returns the identifier of an authorized viewer
// it returns nil if no viewer was found in the session
func Viewer(ctx context.Context) *Identifier {
	val := ctx.Value(AuthContextID)
	if val == nil {
		return nil
	}

	uid, ok := val.(session.SessionUID)
	if ok == false {
		return nil
	}

	id := Identifier(uid)

	return &id
}

func NewService(session session.Service, authable Authenticatable) Service {
	return &service{
		session:  session,
		authable: authable,
	}
}
