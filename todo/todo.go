package todo

import (
	"errors"
	"strings"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
	"github.com/pborman/uuid"
)

// TodoID uniquely identifies a todo
type TodoID string

// Contents of a Todo
type TodoText string

// True for done, false otherwise
type TodoDone bool

type Todo struct {
	ID    TodoID      `json:"id"`
	Text  TodoText    `json:"text"`
	Done  TodoDone    `json:"done"`
	Owner user.UserID `json:"owner"`
}

func (t *Todo) UpdateText(text TodoText) {
	t.Text = text
}

func (t *Todo) ToggleDone() {
	t.Done = !t.Done
}

func New(id TodoID, text TodoText, done TodoDone) *Todo {
	return &Todo{
		ID:   id,
		Text: text,
		Done: done,
	}
}

// ErrUnknown is used when a todo could not be found.
var ErrUnknown = errors.New("unknown todo")

// NextTodoID generates a new todo ID.
// TODO: Move to infrastructure(?)
func NextTodoID() TodoID {
	return TodoID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}
