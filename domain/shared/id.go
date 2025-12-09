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

// ParseID parses a string into an ID.
func ParseID(s string) (ID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID(u), nil
}

// String returns the string representation of the ID.
func (id ID) String() string {
	return ulid.ULID(id).String()
}
