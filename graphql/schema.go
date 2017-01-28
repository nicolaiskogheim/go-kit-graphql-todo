package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
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
		"done": {
			Type: TodoDoneType,
		},
	},
})

// The other types on todo is handled correctly, but we need
// to handle TodoBoolean, or else it will always evaluate to
// false. See graphql.Boolean and their coerceBool for why
// this is.
var TodoDoneType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "TodoBoolean",
	Description: "The `TodoBoolean` scalar type represents `true` or `false`.",
	Serialize:   coerceTodoBool,
	ParseValue:  coerceTodoBool,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.BooleanValue:
			return valueAST.Value
		}
		return nil
	},
})

func coerceTodoBool(value interface{}) interface{} {
	if val, ok := value.(todo.TodoDone); ok {
		return bool(val)
	}
	return false
}

// TODO(nicolai): Does the schemas belong in the services they administer?
func NewSchema(s todo.Service) (graphql.Schema, error) {

	// TODO(nicolai): Run `gofmt -s` on this sometime
	return graphql.NewSchema(graphql.SchemaConfig{

		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name:        "TodosQuery",
				Description: "Returns all todos",
				Fields: graphql.Fields{
					"todos": &graphql.Field{
						Type:        graphql.NewList(todoType),
						Description: "List of todos",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return s.FindAll(), nil
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
							return s.Find(todoID)
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
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						todoText := todo.TodoText(p.Args["text"].(string))
						todoDone := todo.TodoDone(p.Args["done"].(bool))

						id := todo.NextTodoID()
						t := todo.New(id, todoText, todoDone)
						err := s.Add(t)
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
						todo, err := s.Toggle(todoID)

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

						return s.Remove(todoID)
					},
				},
			},
		}),
	})

}
