package user

import "testing"

var (
	id       = UserID("user_id")
	name     = UserName("jimmy")
	email    = UserEmail("jimmy@jimmy.com")
	password = UserPassword("password")
)

func TestUser(t *testing.T) {
	u := New(id, name, email, password)

	if want, have := id, u.ID; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := name, u.Name; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := email, u.Email; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}

	if want, have := password, u.Password; want != have {
		t.Fatalf("want %+v, have %+v", want, have)
	}
}
