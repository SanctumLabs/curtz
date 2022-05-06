package identifier

import (
	"log"

	"github.com/rs/xid"
)

type ID xid.ID

func New() ID {
	return ID(xid.New())
}

func (id ID) String() string {
	return xid.ID(id).String()
}

func (id ID) FromString(idString string) ID {
	value, err := xid.FromString(idString)
	if err != nil {
		log.Println("failed to parse id: ", err.Error())
	}

	return ID(value)
}
