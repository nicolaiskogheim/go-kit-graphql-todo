package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

// The schema for our Todo type
// This lets GraphQL know what fields are available, and of which types.
var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": {
			Type: graphql.String,
		},
		"text": {
			Type: graphql.String,
		},
		// The other types on todo is handled correctly, but we need
		// to handle TodoBoolean, or else it will always evaluate to
		// false. See graphql.Boolean and their coerceBool for why
		// this is.
		"done": {
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*todo.Todo); ok {
					return bool(t.Done), nil
				}
				return nil, nil
			},
		},
		"owner": {
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// TODO(nicolai): how the fuck do we get user.Service
				// in here?
				// u, err := us.Find(user.UserID(p.Source.OwnerID.(string)))
				// return u, err
				return nil, nil
			},
		},
	},
})

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": {
			Type: graphql.String,
		},
		"name": {
			Type: graphql.String,
		},
		"email": {
			Type: graphql.String,
		},
	},
})

// TODO(nicolai): Does the schemas belong in the services they administer?
func NewSchema(ts todo.Service, us user.Service) (graphql.Schema, error) {

	// TODO(nicolai): Run `gofmt -s` on this sometime
	return graphql.NewSchema(graphql.SchemaConfig{

		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name:        "Query",
				Description: "Query todos and users",
				Fields: graphql.Fields{
					"todos": &graphql.Field{
						Type:        graphql.NewList(todoType),
						Description: "List of todos",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {

							return ts.FindAll(), nil
						},
					},
					"todo": &graphql.Field{
						Type:        todoType,
						Description: "Find a specific todo by id. Returns \"unknown todo\" if not found.",
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(graphql.String),
							},
						},
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							todoID := todo.TodoID(p.Args["id"].(string))
							t, err := ts.Find(todoID)
							if err != nil {
								return nil, err
							}

							u, err := us.Find(t.OwnerID)
							if err != nil {
								return t, user.ErrUnknown
							}

							todoWithUser := struct {
								ID    todo.TodoID   `json:"id"`
								Text  todo.TodoText `json:"text"`
								Done  todo.TodoDone `json:"done"`
								Owner user.User     `json:"owner"`
							}{t.ID, t.Text, t.Done, *u}

							return todoWithUser, nil
						},
					},
					"users": &graphql.Field{
						Type:        graphql.NewList(userType),
						Description: "List of users",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return us.FindAll(), nil
						},
					},
					"user": &graphql.Field{
						Type:        userType,
						Description: "Find a user",
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(graphql.String),
							},
						},
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							id := user.UserID(p.Args["id"].(string))
							return us.Find(id)
						},
					},
				},
			},
		),

		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "TodosMutations",
			Fields: graphql.Fields{
				"addTodo": &graphql.Field{
					Type:        todoType,
					Description: "Creates a new todo and stores it.",
					Args: graphql.FieldConfigArgument{
						"text": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"done": &graphql.ArgumentConfig{
							Type:         graphql.Boolean,
							DefaultValue: false,
						},
						"owner_id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						todoText := todo.TodoText(p.Args["text"].(string))
						todoDone := todo.TodoDone(p.Args["done"].(bool))
						todoOwnerID := user.UserID(p.Args["owner_id"].(string))

						id := todo.NextTodoID()
						t := todo.New(id, todoText, todoDone, todoOwnerID)
						err := ts.Add(t)
						return t, err
					},
				},
				"toggleTodo": &graphql.Field{
					Type:        todoType,
					Description: "Toggles the 'done' field of a todo",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						todoID := todo.TodoID(p.Args["id"].(string))
						todo, err := ts.Toggle(todoID)

						return todo, err
					},
				},
				"deleteTodo": &graphql.Field{
					Type:        todoType,
					Description: "Deletes the speciefied todo, or returns \"unknown todo\" if not found",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						todoID := todo.TodoID(p.Args["id"].(string))

						return ts.Remove(todoID)
					},
				},
				"addUser": &graphql.Field{
					Type:        userType,
					Description: "Add a user",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						name := user.UserName(p.Args["name"].(string))
						email := user.UserEmail(p.Args["email"].(string))
						password := user.UserPassword(p.Args["password"].(string))

						// TODO(nicolai): do this through the service?
						u := user.New(user.NextUserID(), name, email, password)
						err := us.Add(u)
						return u, err
					},
				},
			},
		}),
	})

}
