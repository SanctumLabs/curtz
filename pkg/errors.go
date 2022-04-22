package pkg

import "errors"

var (
	ErrServerError = errors.New("unexpected error encountered in server side")
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func NewError(msg string) Error {
	return Error{msg}
}
