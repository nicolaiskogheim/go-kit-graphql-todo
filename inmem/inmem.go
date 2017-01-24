package inmem

import (
	"sync"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
)

type todoRepository struct {
	mtx   sync.RWMutex
	todos map[todo.TodoID]*todo.Todo
}

func (r *todoRepository) Store(t *todo.Todo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.todos[t.ID] = t

	return nil
}

func (r *todoRepository) Delete(id todo.TodoID) (*todo.Todo, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if val, ok := r.todos[id]; ok {
		delete(r.todos, id)
		return val, nil
	}

	return nil, todo.ErrUnknown
}

func (r *todoRepository) Find(id todo.TodoID) (*todo.Todo, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.todos[id]; ok {
		return val, nil
	}

	return nil, todo.ErrUnknown
}

func (r *todoRepository) FindAll() []*todo.Todo {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	todos := make([]*todo.Todo, 0, len(r.todos))
	for _, val := range r.todos {
		todos = append(todos, val)
	}

	return todos
}

func NewTodoRepository() todo.TodoRepository {
	r := &todoRepository{
		todos: make(map[todo.TodoID]*todo.Todo),
	}

	r.todos[todo.Todo1.ID] = todo.Todo1
	r.todos[todo.Todo2.ID] = todo.Todo2
	r.todos[todo.Todo3.ID] = todo.Todo3
	r.todos[todo.Todo4.ID] = todo.Todo4

	return r
}
