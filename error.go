package gocache

import (
	"fmt"
	"time"
)

type Error struct {
	StatusCode int
	Cause      string
	Instant    string
}

func (e *Error) Error() string {
	return fmt.Sprintf(e.Cause)
}

func NewError(statusCode int, cause string) *Error {
	return &Error{StatusCode: statusCode, Cause: cause, Instant: time.Now().Format("2006-01-02T15:04:05.999999Z")}
}
