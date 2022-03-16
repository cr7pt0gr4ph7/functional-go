package cursor

// CursorIter implements a mutable iterator over an immutable cursor.
type CursorIter[C CursorWithSelfType[C, T], T any] struct {
	Cursor C
}

func ToIter[T any, C CursorWithSelfType[C, T]](cursor C) *CursorIter[C, T] {
	return &CursorIter[C, T]{cursor}
}

func (it *CursorIter[C, T]) Next() (item T, ok bool) {
	item, it.Cursor, ok = it.Cursor.Advance()
	return
}

// CursorIterWithCurrent implements a mutable iterator over an immutable cursor.
//
// We try to attempt to implement all stateful iteration patterns
// using a single type, which leads to some code acrobatics internally.
type CursorIterWithCurrent[C CursorWithSelfType[C, T], T any] struct {
	current   iterState[C, T]
	next      iterState[C, T]
	nextValid bool
}

func ToIterWithCurrent[T any, C CursorWithSelfType[C, T]](cursor C) *CursorIterWithCurrent[C, T] {
	it := new(CursorIterWithCurrent[C, T])
	it.current.cursor = cursor
	return it
}

func (it *CursorIterWithCurrent[C, T]) MaybeCurrent() (value T, ok bool) {
	return it.current.value, it.current.hasValue
}

func (it *CursorIterWithCurrent[C, T]) Current() T {
	if !it.current.hasValue {
		panic("iterator has no current value")
	}
	return it.current.value
}

func (it *CursorIterWithCurrent[C, T]) HasCurrent() bool {
	return it.current.hasValue
}

func (it *CursorIterWithCurrent[C, T]) Next() T {
	it.takeNext()
	return it.Current()
}

func (it *CursorIterWithCurrent[C, T]) HasNext() bool {
	return it.peekNext().hasValue
}

func (it *CursorIterWithCurrent[C, T]) TryMoveNext() bool {
	return it.takeNext().hasValue
}

type iterState[C any, T any] struct {
	cursor   C
	value    T
	hasValue bool
}

func (it *CursorIterWithCurrent[C, T]) queryNext() iterState[C, T] {
	nextValue, nextCursor, ok := it.current.cursor.Advance()

	return iterState[C, T]{
		cursor:   nextCursor,
		value:    nextValue,
		hasValue: ok,
	}
}

// Return the next iterator state, but do not make it the current state.
func (it *CursorIterWithCurrent[C, T]) peekNext() *iterState[C, T] {
	// Get and cache the next iterator state
	// if not already cached previously.
	if !it.nextValid {
		it.next = it.queryNext()
		it.nextValid = true
	}
	return &it.next
}

// Move the current iterator state to the next element.
func (it *CursorIterWithCurrent[C, T]) takeNext() *iterState[C, T] {
	// Have we already cached the next element (e.g. due to a previous HasNext() call)?
	if it.nextValid {
		// Make the cached `next` state the `current` iterator state.
		it.current = it.next
		it.nextValid = false
		// it.next = iterState[C, T]{}
	} else {
		it.current = it.queryNext()
	}
	return &it.current
}
