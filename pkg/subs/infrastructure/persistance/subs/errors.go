package subs

import (
	"errors"
	"fmt"
)

var (
	ErrConcurrentModification = errors.New("concurrent modification")
	ErrInvalidSortingField    = errors.New("invalid sorting field")
)

type ErrorRetriesExceeded struct {
	err error
}

func NewErrorRetriesExceeded(err error) *ErrorRetriesExceeded {
	return &ErrorRetriesExceeded{err: err}
}
func (e ErrorRetriesExceeded) Error() string {
	return fmt.Sprintf("превышено максимальное кол-во попыток, последняя ошибка: %v", e.err)
}
