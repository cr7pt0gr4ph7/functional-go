

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

func (c Chain[T]) Cursor() Cursor[T] {
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

type fromSlice[T any] struct {
	slice []T
}

type fromCursor[T any] struct {
	cursor Cursor[T]
}

type chainCursor[T any] struct {
	impl chainImpl[T]
}

func (c chainCursor[T]) Advance() (T, Cursor[T], bool) {
	switch i := c.impl.(type) {
	case nil:
		var t T
		return t, c, false
	case one[T]:
		return i.item, emptyCursor[T]{}, true
	case concat[T]:
		return chainedCursor[T]{chainCursor[T]{i.left}, chainCursor[T]{i.right}}.Advance()
	case fromSlice[T]:
		return sliceCursor[T]{i.slice}.Advance()
	case fromCursor[T]:
		return i.cursor.Advance()
	default:
		fmt.Println(reflect.TypeOf(i))
		panic("unreachable")
	}
}
