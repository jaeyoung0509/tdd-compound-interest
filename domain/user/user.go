package user

import (
	"strings"
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

var (
	ErrInvalidUserName = shared.NewDomainError(shared.CodeValidation, "invalid user name")
	ErrInvalidUserID   = shared.NewDomainError(shared.CodeValidation, "invalid user id")
)

type ID struct {
	value shared.ID
}

func NewID() ID {
	return ID{value: shared.NewID()}
}

func (id ID) Value() shared.ID {
	return id.value
}

func (id ID) IsZero() bool {
	return shared.IsZero(id.value)
}

// IDFromShared creates a user.ID from shared.ID.
func IDFromShared(id shared.ID) ID {
	return ID{value: id}
}

// User carries minimal identity info with an aggregate-scoped ID.
type User struct {
	id        ID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func New(name string, now time.Time) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidUserName
	}
	if now.IsZero() {
		now = time.Now()
	}
	return &User{
		id:        NewID(),
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Reconstitute(id ID, name string, createdAt, updatedAt time.Time) (*User, error) {
	name = strings.TrimSpace(name)
	if id.IsZero() {
		return nil, ErrInvalidUserID
	}
	if name == "" {
		return nil, ErrInvalidUserName
	}
	if createdAt.IsZero() {
		return nil, ErrInvalidUserID
	}
	return &User{id: id, name: name, createdAt: createdAt, updatedAt: updatedAt}, nil
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
