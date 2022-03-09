package eval

import (
	"github.com/cr7pt0gr4ph7/functional-go/option"
)

type evalKind byte

const (
	kNow evalKind = iota
	kLater
	kAlways
	kFlatMap
	kMemoize
	kDefer
)

type evalImpl[A any] interface {
	kind() evalKind
	Value() A
	Memoize() Eval[A]
}

func fromImpl[A any, I evalImpl[A]](impl I) Eval[A] {
	return Eval[A]{impl: impl}
}

func (d *deferImpl[A]) kind() evalKind      { return kDefer }
func (f *flatMapImpl[S, A]) kind() evalKind { return kFlatMap }
func (m *memoizeImpl[A]) kind() evalKind    { return kMemoize }

// deferImpl provides the implementation for `Defer()`.
type deferImpl[A any] struct {
	run func() Eval[A]
}

func (d *deferImpl[A]) Value() A {
	return d.run().Value()
}

func (d *deferImpl[A]) Memoize() Eval[A] {
	return wrapWithMemoize(fromImpl[A](d))
}

// flatMapImpl provides the implementation for `FlatMap()` and `Map()`.
type flatMapImpl[Start any, A any] struct {
	start func() Eval[Start]
	run   func(start Start) Eval[A]
}

func (f *flatMapImpl[S, A]) Value() A {
	return f.run(f.start().Value()).Value()
}

func (f *flatMapImpl[S, A]) Memoize() Eval[A] {
	return wrapWithMemoize(fromImpl[A](f))
}

// memoizeImpl provides the implementation for `wrapWithMemoize(eval)`.
type memoizeImpl[A any] struct {
	eval   Eval[A]
	result option.Optional[A]
}

func wrapWithMemoize[A any](eval Eval[A]) Eval[A] {
	return fromImpl[A](&memoizeImpl[A]{eval: eval})
}

func (m *memoizeImpl[A]) Value() A {
	if v, ok := m.result.Value(); ok {
		return v
	} else {
		r := m.eval.Value()
		m.result = option.Some(r)
		return r
	}

}

func (m *memoizeImpl[A]) Memoize() Eval[A] {
	return Eval[A]{impl: m}
}
