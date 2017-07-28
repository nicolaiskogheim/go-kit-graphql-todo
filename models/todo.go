package models

import (
	"github.com/nicolaiskogheim/go-kit-graphql-todo/schema"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

var _ schema.TodoInterface = (*Todo)(nil)

// Todo wraps todo.Todo and stuff needed to resolve it
type Todo struct {
	source      todo.Todo
	UserService user.Service
}

// IdField resolves the ID field on todo.Todo
func (todo Todo) IdField() (*string, error) {
	id := string(todo.source.ID)
	return &id, nil
}

// TextField resolves the Text field on todo.Todo
func (todo Todo) TextField() (*string, error) {
	text := string(todo.source.Text)
	return &text, nil
}

// DoneField resolves the Done field on todo.Todo
func (todo Todo) DoneField() (*bool, error) {
	done := bool(todo.source.Done)
	return &done, nil
}

// OwnerField resolves the Owner field on todo.Todo
func (todo Todo) OwnerField() (schema.UserInterface, error) {
	user, err := todo.UserService.Find(todo.source.OwnerID)

	if err != nil {
		return nil, err
	}

	return User{source: *user}, nil
}
