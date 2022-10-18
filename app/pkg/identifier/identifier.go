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

func (id ID) Bytes() []byte {
	return xid.ID(id).Bytes()
}

// FromString parses a string into an ID
func (id ID) FromString(idString string) (ID, error) {
	value, err := xid.FromString(idString)
	if err != nil {
		log.Println("failed to parse id: ", err.Error())
		return ID{}, err
	}

	return ID(value), nil
}

func FromBytes(idString []byte) (ID, error) {
	value, err := xid.FromBytes(idString)
	if err != nil {
		return ID{}, err
	}

	return ID(value), nil
}

func (id ID) FromBytes(idString []byte) (ID, error) {
	return FromBytes(idString)
}
