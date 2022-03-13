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

func (r Result[T]) ErrorOrPanic() error {
	if r.err == nil {
		panic("not an error")
	}
	return r.err
}

func Map[A any, B any](r Result[A], mapValue func(value A) B) Result[B] {
	if v, e := r.Extract(); e == nil {
		return Ok(mapValue(v))
	} else {
		return Error[B](e)
	}
}

func MapError[A any](r Result[A], mapError func(err error) error) Result[A] {
	if _, e := r.Extract(); e == nil {
		return r
	} else {
		return Error[A](mapError(e))
	}
}

func FlatMap[A any, B any](r Result[A], mapValue func(value A) Result[B]) Result[B] {
	if v, e := r.Extract(); e == nil {
		return mapValue(v)
	} else {
		return Error[B](e)
	}
}
