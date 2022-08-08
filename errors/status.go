package errors

import (
	"fmt"
)

type status struct {
	Status  int
	Code    int
	Message string
}

type WarpError struct {
	e *status
}

func new(s int, c int, msg string) *status {
	return &status{Status: s, Code: c, Message: msg}
}

func (e *WarpError) Error() string {
	return e.e.Message
}

func (e *WarpError) Status() int {
	return e.e.Status
}

func (e *WarpError) Code() int {
	return e.e.Code
}

func (s *status) Err() error {
	return &WarpError{e: s}
}

func Error(s int, c int, msg string) error {
	return new(s, c, msg).Err()
}

func Errorf(s int, c int, format string, a ...interface{}) error {
	return Error(s, c, fmt.Sprintf(format, a...))
}

func Code(s int, c int) error {
	return Error(s, c, "")
}
