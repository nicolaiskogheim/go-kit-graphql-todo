package models

import (
	"github.com/nicolaiskogheim/go-kit-graphql-todo/schema"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

var _ schema.UserInterface = (*User)(nil)

// User wraps user.User and stuff needed to resolve it
type User struct {
	source user.User
}

// IdField resolves the ID field on user.User
func (user User) IdField() (*string, error) {
	id := string(user.source.ID)
	return &id, nil
}

// NameField resolves the Name field on user.User
func (user User) NameField() (*string, error) {
	name := string(user.source.Name)
	return &name, nil
}
