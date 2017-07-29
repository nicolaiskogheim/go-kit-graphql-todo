package models

import (
	"errors"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/auth"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/schema"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
	"golang.org/x/net/context"
)

var _ schema.QueryInterface = (*Root)(nil)

var _ schema.MutationInterface = (*Root)(nil)

type Root struct {
	TodoService todo.Service
	UserService user.Service
}

// TodosQuery ...
func (root Root) TodosQuery(
	ctx context.Context,
) ([]schema.TodoInterface, error) {
	var todos []schema.TodoInterface
	for _, val := range root.TodoService.FindAll() {
		todos = append(todos, Todo{
			source:      *val,
			UserService: root.UserService,
		})
	}
	return todos, nil
}

// TodoQuery resolves todo( id: ID! )
func (root Root) TodoQuery(
	ctx context.Context,
	id string,
) (schema.TodoInterface, error) {
	todo, err := root.TodoService.Find(todo.TodoID(id))

	if err != nil {
		return nil, err
	}

	t := Todo{
		source:      *todo,
		UserService: root.UserService,
	}

	return t, nil
}

// UsersQuery resolves users()
func (root Root) UsersQuery(
	ctx context.Context,
) ([]schema.UserInterface, error) {
	var users []schema.UserInterface
	for _, val := range root.UserService.FindAll() {
		users = append(users, User{
			source: *val,
		})
	}

	return users, nil
}

// UserQuery resolves user( id: ID )
func (root Root) UserQuery(
	ctx context.Context,
	id string,
) (schema.UserInterface, error) {
	user, err := root.UserService.Find(user.UserID(id))

	if err != nil {
		return nil, err
	}

	u := User{
		source: *user,
	}

	return u, nil
}

// AddTodoMutation resolves addTodo( text: String!, done: Boolean = false )
func (root Root) AddTodoMutation(
	ctx context.Context,
	text string,
	done bool,
	owner string,
) (schema.TodoInterface, error) {
	todo := todo.New(
		todo.NextTodoID(),
		todo.TodoText(text),
		todo.TodoDone(done),
		user.UserID(owner),
	)
	err := root.TodoService.Add(todo)

	if err != nil {
		return nil, err
	}

	return Todo{source: *todo, UserService: root.UserService}, nil
}

// ToggleTodoMutation resolves toggleTodo( id: ID! )
func (root Root) ToggleTodoMutation(
	ctx context.Context,
	id string,
) (schema.TodoInterface, error) {
	uid := auth.Viewer(ctx)
	user, err := root.UserService.Find(user.UserID(*uid))

	if err != nil {
		return nil, errors.New("user not found")
	}

	todo, err := root.TodoService.Toggle(*user, todo.TodoID(id))

	if err != nil {
		return nil, err
	}

	return Todo{source: *todo, UserService: root.UserService}, nil
}

// DeleteTodoMutation resolves deleteTodo( id: ID! )
func (root Root) DeleteTodoMutation(
	ctx context.Context,
	id string,
) (schema.TodoInterface, error) {
	todo, err := root.TodoService.Remove(todo.TodoID(id))

	if err != nil {
		return nil, err
	}

	return Todo{source: *todo, UserService: root.UserService}, nil
}

// AddUserMutation resolves addUser( name: String!, email: String!, password: String! )
func (root Root) AddUserMutation(
	ctx context.Context,
	name string,
	email string,
	password string,
) (schema.UserInterface, error) {
	user := user.New(
		user.NextUserID(),
		user.UserName(name),
		user.UserEmail(email),
		user.UserPassword(password),
	)

	err := root.UserService.Add(user)

	if err != nil {
		return nil, err
	}

	return User{source: *user}, nil
}

func (root Root) ViewerQuery(
	ctx context.Context,
) (schema.UserInterface, error) {
	id := auth.Viewer(ctx)

	if id == nil {
		return nil, nil
	}

	user, err := root.UserService.Find(user.UserID(id.ToString()))
	if err != nil {
		return nil, err
	}

	return User{source: *user}, nil
}
