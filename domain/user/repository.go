package user

import (
	"context"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

var ErrUserNotFound = shared.NewDomainError(shared.CodeNotFound, "user not found")

// Repository abstracts persistence for the User aggregate.
type Repository interface {
	Get(ctx context.Context, id ID) (*User, error)
	Save(ctx context.Context, user *User) error
}
