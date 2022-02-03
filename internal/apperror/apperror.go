package apperror

import (
	"fmt"
)

type Error struct {
	Msg    string
	DevMsg string
	Code   string
	Err    error
}

func NewError(msg, devMsg, code string, err error) *Error {
	if err == nil {
		err = fmt.Errorf(msg)
	}
	return &Error{
		Msg:    msg,
		DevMsg: devMsg,
		Code:   code,
		Err:    err,
	}
}

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) Unwrap() error {
	return e.Err
}
