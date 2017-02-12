package todo

import "github.com/nicolaiskogheim/go-kit-graphql-todo/user"

type Service interface {
	Add(t *Todo) error
	Toggle(id TodoID) (*Todo, error)
	Remove(id TodoID) (*Todo, error)
	FindAll() []*Todo
	Find(id TodoID) (*Todo, error)
	FindByUserID(id user.UserID) []*Todo
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

func (s *service) Toggle(id TodoID) (*Todo, error) {

	t, err := s.repository.Find(id)
	if err != nil {
		return nil, err
	}

	t.ToggleDone()
	s.repository.Update(t)

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

func (s *service) FindByUserID(id user.UserID) []*Todo {
	return s.repository.FindByUserID(id)
}

type TodoRepository interface {
	Store(t *Todo) error
	Update(t *Todo) error
	Delete(id TodoID) (*Todo, error)
	Find(id TodoID) (*Todo, error)
	FindAll() []*Todo
	FindByUserID(id user.UserID) []*Todo
}
