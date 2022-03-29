package chain

import (
	"github.com/cr7pt0gr4ph7/functional-go/collections/immutable/cursor"
)

// Immutable list with O(1) prepend, append and concat.
type Chain[T any] struct {
	impl chainImpl[T]
}

func (c Chain[T]) Prepend(item T) Chain[T] {
	if c.impl == nil {
		return Chain[T]{one[T]{item}}
	}
	return Chain[T]{concat[T]{one[T]{item}, c.impl}}
}

func (c Chain[T]) Append(item T) Chain[T] {
	if c.impl == nil {
		return Chain[T]{one[T]{item}}
	}
	return Chain[T]{concat[T]{c.impl, one[T]{item}}}
}

func (c Chain[T]) Concat(other Chain[T]) Chain[T] {
	if c.impl == nil {
		return other
	}
	if other.impl == nil {
		return c
	}
	return Chain[T]{concat[T]{c.impl, other.impl}}
}

func (c Chain[T]) Empty() bool {
	return c.impl == nil
}

func (c Chain[T]) Len() int {
	i := 0
	cur := c.Cursor()
	for {
		var ok bool
		_, cur, ok = cur.Advance()
		if !ok {
			break
		}
		i++
	}
	return i
}

func (c Chain[T]) Cursor() cursor.Cursor[T] {
	return chainCursor[T]{c.impl}
}

type chainImpl[T any] interface{ chainImpl(_ T) }

func (_ one[T]) chainImpl(_ T)        {}
func (_ concat[T]) chainImpl(_ T)     {}
func (_ fromSlice[T]) chainImpl(_ T)  {}
func (_ fromCursor[T]) chainImpl(_ T) {}

type one[T any] struct {
	item T
}

type concat[T any] struct {
	left  chainImpl[T]
	right chainImpl[T]
}

type chainCursor[T any] struct {
	impl chainImpl[T]
}

func (c chainCursor[T]) Advance() (T, cursor.Cursor[T], bool) {
	switch i := c.impl.(type) {
	case nil:
		var t T
		return t, c, false
	case one[T]:
		return i.item, cursor.Empty[T](), true
	case concat[T]:
		return cursor.Concat[T](i.left.Cursor(), i.right.Cursor()).Advance()
	case fromSlice[T]:
		return cursor.FromSlice(i.slice).Advance()
	case fromCursor[T]:
		return i.cursor.Advance()
	default:
		panic("unreachable")
	}
}
