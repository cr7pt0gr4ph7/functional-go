package result

import (
	"errors"
)

var nilError = errors.New("no error given")

func ensureError(err error) error {
	if err == nil {
		return nilError
	}
	return err
}

type Result[T any] struct {
	value T
	err   error
}

func From[T any](value T, err error) Result[T] {
	return Result[T]{value: value, err: err}
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Error[T any](err error) Result[T] {
	return Result[T]{err: ensureError(err)}
}

func (r Result[_]) IsOk() bool {
	return r.err == nil
}

func (r Result[_]) IsError() bool {
	return r.err != nil
}

func (r Result[T]) Extract() (T, error) {
	return r.value, r.err
}

func (r Result[T]) ValueOrPanic() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

func (r Result[_]) Error() error {
	return r.err
}
