package hlist

import (
	"fmt"
)

// HList can be used as a constraint for heterogeneous lists.
type HList interface {
	hList() // marker method
	hExtract() (head any, tail HList, ok bool)
	Empty() bool
	Len() int
	countTailRec(count int) int
	Slice() []any
	String() string
	Each(f func(item any))
	EachWithIndex(f func(index int, item any))
}

// NonEmpty can be used as a constraint for non-empty heterogeneous lists.
type NonEmpty interface {
	HList
	nonEmpty() // marker method
}

func (_ Nil) hList()           {}
func (_ Cons[_, _]) hList()    {}
func (_ Cons[_, _]) nonEmpty() {}

func (_ Nil) hExtract() (any, HList, bool)        { return nil, nil, false }
func (c Cons[H, T]) hExtract() (any, HList, bool) { return c.Head, c.Tail, true }

// =========
// :: Nil ::
// =========

// Nil represents the empty list.
type Nil struct{}

func (_ Nil) Empty() bool                { return true }
func (_ Nil) Len() int                   { return 0 }
func (_ Nil) countTailRec(count int) int { return count }
func (_ Nil) Slice() (_ []any)           { return }
func (_ Nil) String() string             { return "Nil" }

func (_ Nil) Each(f func(item any))                     {}
func (_ Nil) EachWithIndex(f func(index int, item any)) {}

// ==========
// :: Cons ::
// ==========

// Cons represents a list with a head element and a tail list.
type Cons[Head any, Tail HList] struct {
	Head Head
	Tail Tail
}

func (_ Cons[_, _]) Empty() bool                { return false }
func (c Cons[_, _]) Len() int                   { return c.countTailRec(0) }
func (c Cons[_, _]) countTailRec(count int) int { return c.Tail.countTailRec(1 + count) }

func (c Cons[_, _]) Slice() []any {
	return append([]any{c.Head}, c.Tail.Slice()...)
}

func (c Cons[H, T]) String() string {
	return fmt.Sprintf("%v ::: %v", c.Head, c.Tail)
}

func (c Cons[H, T]) Each(f func(item any)) {
	for h, t, ok := c.hExtract(); ok; h, t, ok = t.hExtract() {
		f(h)
	}
}

func (c Cons[H, T]) EachWithIndex(f func(index int, item any)) {
	i := 0
	for h, t, ok := c.hExtract(); ok; h, t, ok = t.hExtract() {
		f(i, h)
		i++
	}
}

// =============
// :: Methods ::
// =============

func Prepend[H any, T HList](head H, tail T) Cons[H, T] {
	return Cons[H, T]{Head: head, Tail: tail}
}

func FoldLeft[L HList, R any](l L, initial R, foldFn func(item any, acc R) R) R {
	acc := initial
	for h, t, ok := l.hExtract(); ok; h, t, ok = t.hExtract() {
		acc = foldFn(h, acc)
	}
	return acc
}

func MapFoldLeft[L HList, T any, R any](l L, initial R, mapFn func(item any) T, foldFn func(item T, acc R) R) R {
	acc := initial
	for h, t, ok := l.hExtract(); ok; h, t, ok = t.hExtract() {
		acc = foldFn(mapFn(h), acc)
	}
	return acc
}

func MapReduceLeft[L NonEmpty, R any](l L, mapFn func(item any) R, reduceFn func(item R, acc R) R) R {
	h, t, ok := l.hExtract()
	acc := mapFn(h)
	for h, t, ok = l.hExtract(); ok; h, t, ok = t.hExtract() {
		acc = reduceFn(mapFn(h), acc)
	}
	return acc
}

func UnsafeSliceOf[A any, H any, T HList](l Cons[H, T]) []A {
	xs := l.Slice()
	r := make([]A, len(xs))
	for i, x := range xs {
		r[i] = x.(any).(A)
	}
	return r
}
