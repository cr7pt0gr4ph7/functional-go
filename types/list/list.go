package list

// Represents a simple immutable list.
type List[T any] struct{ *entry[T] }

type entry[T any] struct {
	head T
	tail *entry[T]
}

func Empty[T any]() List[T] {
	return List[T]{}
}

func New[T any](s ...T) List[T] {
	var r List[T]
	for i := len(s) - 1; i >= 0; i-- {
		r = r.Push(s[i])
	}
	return r
}

func NewReverse[T any](s ...T) List[T] {
	var r List[T]
	for _, x := range s {
		r = r.Push(x)
	}
	return r
}

func (l *entry[T]) Prepend(item T) List[T] {
	return List[T]{&entry[T]{head: item, tail: l}}
}

func (l *entry[T]) Push(item T) List[T] {
	return List[T]{&entry[T]{head: item, tail: l}}
}

func (l *entry[T]) Pop() (head T, tail List[T], ok bool) {
	if l == nil {
		ok = false
		return
	} else {
		return l.head, List[T]{l.tail}, true
	}
}

func (l *entry[T]) Empty() bool {
	return l == nil
}

func (l *entry[T]) Len() int {
	i := 0
	for ; l != nil; l = l.tail {
		i++
	}
	return i
}

func (l *entry[T]) Head() (_ T, ok bool) {
	if l == nil {
		ok = false
		return
	} else {
		return l.head, true
	}
}

func (l *entry[T]) Tail() List[T] {
	if l == nil {
		return List[T]{nil}
	} else {
		return List[T]{l.tail}
	}
}
