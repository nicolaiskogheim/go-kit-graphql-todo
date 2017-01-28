package todo

import "testing"

var (
	id      = TodoID("todo_id")
	text1   = TodoText("test todo")
	text2   = TodoText("updated text")
	notDone = TodoDone(false)
	done    = TodoDone(true)
)

func TestTodo(t *testing.T) {
	t1 := New(id, text1, notDone)

	if want, have := id, t1.ID; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := text1, t1.Text; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := notDone, t1.Done; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	t1.UpdateText(text2)
	if want, have := text2, t1.Text; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	t1.ToggleDone(done)
	if want, have := done, t1.Done; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}
}
