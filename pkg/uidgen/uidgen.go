package uidgen

import "github.com/google/uuid"

type UIDGen interface {
	Create() uuid.UUID
	Generate(chars []byte) (string, error)
}

type uidgen struct{}

func New() UIDGen {
	return &uidgen{}
}

// Create creates a new UUID
func (u uidgen) Create() uuid.UUID {
	return uuid.New()
}

// Generate creates a unique id using the given chars
func (u uidgen) Generate(chars []byte) (string, error) {
	id, err := uuid.FromBytes(chars)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
