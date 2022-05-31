package identifier

import (
	"log"

	"github.com/rs/xid"
)

// ID is a unique identifier
type ID xid.ID

// New generates a new ID
func New() ID {
	return ID(xid.New())
}

// String returns the string representation of the ID
func (id ID) String() string {
	return xid.ID(id).String()
}

// FromString parses a string into an ID
func (id ID) FromString(idString string) ID {
	value, err := xid.FromString(idString)
	if err != nil {
		log.Println("failed to parse id: ", err.Error())
	}

	return ID(value)
}
