package inmem

import (
	"testing"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
)

func TestInmem(t *testing.T) {

	tr := &todoRepository{
		todos: make(map[todo.TodoID]*todo.Todo),
	}

	todo1 := &todo.Todo{ID: "123", Text: "Test Inmem", Done: false}

	// Repo should not contain anything
	if want, have := 0, len(tr.FindAll()); want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	tr.Store(todo1)

	t1, err := tr.Find(todo1.ID)
	if err != nil {
		t.Fatalf("should have found todo, got %s", err.Error())
	}

	if want, have := todo1.ID, t1.ID; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := todo1.Text, t1.Text; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := todo1.Done, t1.Done; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	ts := tr.FindAll()
	if want, have := 1, len(ts); want != have {
		t.Fatalf("want %d, have %d", want, have)
	}

	t1, err = tr.Delete(todo1.ID)
	if err != nil {
		t.Fatalf("should delete todo, got %s", err.Error())
	}

	if want, have := 0, len(tr.FindAll()); want != have {
		t.Fatalf("want %d, have %d", want, have)
	}

}
