package session

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

type SessionUID string

func (uid SessionUID) ToString() string {
	return string(uid)
}

type SessionToken string

func (token SessionToken) ToString() string {
	return string(token)
}

type Session struct {
	UID     SessionUID
	Token   SessionToken
	Expires time.Time
}

func New(uid SessionUID, expires time.Time) *Session {
	return &Session{
		UID:     uid,
		Token:   newToken(),
		Expires: expires,
	}
}

func newToken() SessionToken {
	return SessionToken(uuid.New())
}

// ErrUnknown is used when a session could not be found.
var ErrUnknown = errors.New("unknown session")
