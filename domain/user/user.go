package user

import (
	"errors"
	"strings"
	"time"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidUserID   = errors.New("invalid user id")
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

// User carries minimal identity info with an aggregate-scoped ID.
type User struct {
	id        ID
	name      string
	createdAt time.Time
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
	}, nil
}

func Reconstitute(id ID, name string, createdAt time.Time) (*User, error) {
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
	return &User{id: id, name: name, createdAt: createdAt}, nil
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
