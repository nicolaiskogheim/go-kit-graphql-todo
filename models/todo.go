package models

import (
	"github.com/nicolaiskogheim/go-kit-graphql-todo/schema"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
)

var _ schema.TodoInterface = (*Todo)(nil)

// Todo wraps todo.Todo and stuff needed to resolve it
type Todo struct {
	source todo.Todo
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
