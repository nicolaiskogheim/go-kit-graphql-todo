package session

import "time"

var (
	Session1 = &Session{UID: "2C2E7C8D",
		Token:   "a282e4ca-b74a-4f51-a27d-28bbf6287729",
		Expires: time.Now().Add(time.Minute * 60),
	}
)
