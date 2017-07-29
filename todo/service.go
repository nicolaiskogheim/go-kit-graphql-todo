package todo

import (
	"errors"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

type Service interface {
	Add(t *Todo) error
	Toggle(user user.User, id TodoID) (*Todo, error)
	Remove(id TodoID) (*Todo, error)
	FindAll() []*Todo
	Find(id TodoID) (*Todo, error)
}

type service struct {
	repository TodoRepository
}

func NewService(repository TodoRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Add(t *Todo) error {
	return s.repository.Store(t)
}

func (s *service) Toggle(user user.User, id TodoID) (*Todo, error) {

	t, err := s.repository.Find(id)
	if err != nil {
		return nil, err
	}

	if user.ID == t.OwnerID {
		t.ToggleDone()
		s.repository.Update(t)
	} else {
		return nil, errors.New("user is not owner of todo")
	}

	return t, nil
}

func (s *service) Remove(id TodoID) (*Todo, error) {
	return s.repository.Delete(id)
}

func (s *service) Find(id TodoID) (*Todo, error) {
	return s.repository.Find(id)
}

func (s *service) FindAll() []*Todo {
	return s.repository.FindAll()
}

type TodoRepository interface {
	Store(t *Todo) error
	Update(t *Todo) error
	Delete(id TodoID) (*Todo, error)
	Find(id TodoID) (*Todo, error)
	FindAll() []*Todo
}
