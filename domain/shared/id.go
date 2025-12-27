package shared

import "github.com/oklog/ulid/v2"

type ID ulid.ULID

var zero ulid.ULID

func NewID() ID {
	return ID(ulid.Make())
}

func IsZero(id ID) bool {
	return ulid.ULID(id) == zero
}
