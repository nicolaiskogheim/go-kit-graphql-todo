package user

import (
	"errors"
	"strings"

	"github.com/pborman/uuid"
)

type UserID string
type UserName string
type UserEmail string
type UserPassword string

type User struct {
	ID       UserID       `json:"id"`
	Name     UserName     `json:"name"`
	Email    UserEmail    `json:"email"`
	Password UserPassword `json:"password"`
}

func New(id UserID, name UserName, email UserEmail, password UserPassword) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// ErrUnknown is used when a user could not be found.
var ErrUnknown = errors.New("unknown user")

// NextUserID generates a new UserID.
func NextUserID() UserID {
	return UserID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}
