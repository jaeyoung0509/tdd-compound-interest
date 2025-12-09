package user

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

// Repository abstracts persistence for the User aggregate.
type Repository interface {
	Get(ctx context.Context, id ID) (*User, error)
	Save(ctx context.Context, user *User) error
}
