package immutable

// Represents an immutable list.
type List[L List_[L, T], T any] interface {
	List_[L, T]
}

// Defines the methods for `List[L, T]`.
//
// This is a separate type due to compiler restrictions.
// This type is public so that it is visible in the documentation.
type List_[L any, T any] interface {
	Prepend(item T) L
	Append(item T) L
	Concat(other L) L
	Empty() bool
	Len() int
	Cursor() Cursor[T]
}

type Cursor[T any] interface {
	Advance() (T, Cursor[T], bool)
}

// CursorIter implements a mutable iterator over an immutable cursor.
type CursorIter[C interface{ Advance() (T, C, bool) }, T any] struct {
	Cursor C
}

func ToCursorIter[T any, C interface{ Advance() (T, C, bool) }](cursor C) CursorIter[C, T] {
	return CursorIter[C, T]{cursor}
}

func NewCursorIter[T any, C interface{ Advance() (T, C, bool) }](cursor C) *CursorIter[C, T] {
	return &CursorIter[C, T]{cursor}
}

func (it *CursorIter[C, T]) Next() (item T, ok bool) {
	item, it.Cursor, ok = it.Cursor.Advance()
	return
}
