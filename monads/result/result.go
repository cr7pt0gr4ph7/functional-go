package result

import (
	"errors"
	"fmt"
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

func (r Result[T]) MustBeValue() T {
	if r.err != nil {
		panic(fmt.Errorf("expected value but got error: %w", r.err))
	}
	return r.value
}

func (r Result[_]) Error() error {
	return r.err
}

func (r Result[T]) MustBeError() error {
	if r.err == nil {
		panic("expected error value, got nil")
	}
	return r.err
}

func (r Result[A]) Catch(whenError func(err error) Result[A]) Result[A] {
	return Catch(r, whenError)
}

func Catch[A any](r Result[A], whenError func(err error) Result[A]) Result[A] {
	if _, e := r.Extract(); e == nil {
		return r
	} else {
		return whenError(e)
	}
}

func (r Result[A]) Tap(tapAction func(r Result[A])) Result[A] {
	tapAction(r)
	return r
}

func (r Result[A]) PassThru(hijackFunc func(r Result[A]) Result[A]) Result[A] {
	return hijackFunc(r)
}

func (r Result[A]) Do(whenOk func(value A)) Result[A] {
	return Do(r, whenOk)
}

func Do[A any](r Result[A], whenOk func(value A)) Result[A] {
	if v, e := r.Extract(); e == nil {
		whenOk(v)
	}
	return r
}

func Fold[A any, B any](r Result[A], whenOk func(value A) B, whenError func(err error) B) B {
	if v, e := r.Extract(); e == nil {
		return whenOk(v)
	} else {
		return whenError(e)
	}
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

func Wrap0[R any, ErrT error](f func() (R, ErrT)) func() Result[R] {
	return func() Result[R] {
		return From(f())
	}
}

func Wrap1[A any, R any, ErrT error](f func(a A) (R, ErrT)) func(a A) Result[R] {
	return func(a A) Result[R] {
		return From(f(a))
	}
}

func Wrap2[A any, B any, R any, ErrT error](f func(a A, b B) (R, ErrT)) func(a A, b B) Result[R] {
	return func(a A, b B) Result[R] {
		return From(f(a, b))
	}
}

func Wrap3[A any, B any, C any, R any, ErrT error](f func(a A, b B, c C) (R, ErrT)) func(a A, b B, c C) Result[R] {
	return func(a A, b B, c C) Result[R] {
		return From(f(a, b, c))
	}
}

func FlatMap2[A any, B any, ErrT error](r Result[A], mapValue func(value A) (B, ErrT)) Result[B] {
	if v, e := r.Extract(); e == nil {
		return From(mapValue(v))
	} else {
		return Error[B](e)
	}
}
