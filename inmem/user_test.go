package inmem

import (
	"testing"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
)

func TestInmem(t *testing.T) {
	r := &userRepository{
		users: make(map[user.UserID]*user.User),
	}

	user := &users.User{ID: "123", Name: "jon", Email: "jon@jon.com", Password: "pass"}

	// Repo should not contain anything
	if want, have := 0, len(r.FindAll()); want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	r.Store(user)

	u, err := r.Find(user.ID)
	if err != nil {
		t.Fatalf("should have found user, got %s", err)
	}

	if want, have := u.ID, user.ID; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := u.Name, user.Name; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := u.Email, user.Email; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := u.Password, user.Password; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	us := r.FindAll()
	if want, have := 1, len(r); want != have {
		t.Fatalf("repo should contain %d, but contains %d users", want, have)
	}

	err = tr.Delete(todo1.ID)
	if err != nil {
		t.Fatalf("should delete todo, got %s", err.Error())
	}

	if want, have := 0, len(tr.FindAll()); want != have {
		t.Fatalf("repo should contain %d, but contains %d users", want, have)
	}
}
