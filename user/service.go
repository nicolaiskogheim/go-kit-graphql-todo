package user

import (
	"net/http"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/auth"
)

type Service interface {
	Add(u *User) error
	Remove(id UserID) error
	Find(id UserID) (*User, error)
	FindAll() []*User
	Authenticate(req http.Request) (*auth.Identifier, error)
}

type service struct {
	repository UserRepository
}

func NewService(repository UserRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Add(u *User) error {
	return s.repository.Store(u)
}

func (s *service) Remove(id UserID) error {
	return s.repository.Delete(id)
}

func (s *service) Find(id UserID) (*User, error) {
	return s.repository.Find(id)
}

func (s *service) FindAll() []*User {
	return s.repository.FindAll()
}

func (s *service) Authenticate(req http.Request) (*auth.Identifier, error) {

	req.ParseForm()
	password := req.Form.Get("password")
	username := req.Form.Get("username")

	user, err := s.repository.FindByCredentials(
		username,
		password,
	)

	if err != nil {
		return nil, err
	}

	identifier := auth.Identifier(user.ID.ToString())
	return &identifier, nil
}

type UserRepository interface {
	Store(u *User) error
	Update(t *User) error
	Delete(id UserID) error
	Find(id UserID) (*User, error)
	FindAll() []*User
	FindByCredentials(username, password string) (*User, error)
}
