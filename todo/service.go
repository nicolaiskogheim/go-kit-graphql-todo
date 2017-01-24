package todo

type Service interface {
	Add(t *Todo) error
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
	Delete(id TodoID) (*Todo, error)
	Find(id TodoID) (*Todo, error)
	FindAll() []*Todo
}
