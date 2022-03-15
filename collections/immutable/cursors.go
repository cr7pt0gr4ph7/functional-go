package immutable

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

type emptyCursor[T any] struct{}

func (c emptyCursor[T]) Advance() (T, Cursor[T], bool) {
	var t T
	return t, c, false
}
