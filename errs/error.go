package errs

import (
	"fmt"

	"util/logs"
	"util/run"
)

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%v, msg=%v\n", e.Code, e.Msg)
}

func New(code int, format string, v ...interface{}) *Error {
	caller := run.Caller(1)
	msg := fmt.Sprintf(format, v...)
	e := &Error{Code: code, Msg: caller + msg}

	logs.Error(e.Error())

	return e
}

var ErrNil = &Error{-1, "<nil>"}
