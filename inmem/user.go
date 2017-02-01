package inmem

import (
	"sync"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

type userRepository struct {
	mtx   sync.RWMutex
	users map[user.UserID]*user.User
}

func (r *userRepository) Store(u *user.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.users[u.ID] = u

	return nil
}

func (r *userRepository) Update(u *user.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.users[u.ID]; ok {
		r.users[u.ID] = u
		return nil
	}

	return user.ErrUnknown
}

func (r *userRepository) Delete(id user.UserID) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.users[id]; ok {
		delete(r.users, id)
		return nil
	}

	return user.ErrUnknown
}

func (r *userRepository) Find(id user.UserID) (*user.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.users[id]; ok {
		return val, nil
	}

	return nil, user.ErrUnknown
}

func (r *userRepository) FindAll() []*user.User {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	users := make([]*user.User, 0, len(r.users))
	for _, val := range r.users {
		users = append(users, val)
	}

	return users
}

func NewUserRepository() user.UserRepository {
	r := &userRepository{
		users: make(map[user.UserID]*user.User),
	}

	r.users[user.User1.ID] = user.User1
	r.users[user.User2.ID] = user.User2
	r.users[user.User3.ID] = user.User3
	r.users[user.User4.ID] = user.User4

	return r
}
