package user

type Service interface {
	Add(u *User) error
	Remove(id UserID) error
	Find(id UserID) (*User, error)
	FindAll() []*User
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

type UserRepository interface {
	Store(u *User) error
	Update(t *User) error
	Delete(id UserID) error
	Find(id UserID) (*User, error)
	FindAll() []*User
}
