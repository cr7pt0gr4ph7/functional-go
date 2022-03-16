package cursor

// Represents an immutable pointer into an immutable data structure
// that can be advanced (which generates a new, separate Cursor)
// to iterate over the data structure in a generic way.
type Cursor[T any] CursorWithSelfType[Cursor[T], T]

type CursorWithSelfType[Self any, T any] interface {
	Advance() (T, Self, bool)
}

// Returns an empty cursor.
func Empty[T any]() Cursor[T] {
	return emptyCursor[T]{}
}

// Returns a cursor that first traverses `left` and then `right`.
func Concat[T any](left Cursor[T], right Cursor[T]) Cursor[T] {
	if isKnownEmpty(right) {
		return left
	} else if isKnownEmpty(left) {
		return right
	} else {
		return chainedCursor[T]{left, right}
	}
}

type maybeEmpty interface {
	isKnownEmpty() bool
}

func isKnownEmpty[T any](cursor Cursor[T]) bool {
	if maybeEmpty, ok := cursor.(maybeEmpty); ok && maybeEmpty.isKnownEmpty() {
		return true
	}
	return false
}

// Returns a cursor that iterates over the specified slice.
func FromSlice[T any, S ~[]T](slice S) Cursor[T] {
	if len(slice) == 0 {
		// Optimization: Return an Empty[T]() cursor if the slice is empty.
		//               This allows the returned cursor to be treated as the identity for Concat().
		return emptyCursor[T]{}
	}
	return sliceCursor[T]{slice}
}

type emptyCursor[T any] struct{}

func (c emptyCursor[T]) isKnownEmpty() bool {
	return true
}

func (c emptyCursor[T]) Advance() (T, Cursor[T], bool) {
	var t T
	return t, c, false
}

type chainedCursor[T any] struct {
	left  Cursor[T]
	right Cursor[T]
}

func (c chainedCursor[T]) Advance() (T, Cursor[T], bool) {
	if t, next, ok := c.left.Advance(); ok {
		return t, chainedCursor[T]{next, c.right}, true
	}
	return c.right.Advance()
}

type sliceCursor[T any] struct {
	slice []T
}

func (c sliceCursor[T]) Advance() (T, Cursor[T], bool) {
	if len(c.slice) == 0 {
		var t T
		return t, emptyCursor[T]{}, false
	}
	return c.slice[0], sliceCursor[T]{c.slice[1:]}, true
}
