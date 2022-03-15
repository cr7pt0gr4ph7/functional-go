package list

// Represents a simple immutable list.
type List[T any] interface {
	Push(item T) List[T]
	Pop() (head T, tail List[T], ok bool)
	Head() (T, bool)
	Tail() List[T]
}

func Empty[T]() List[T] {
	return Nil[T]{}
}

type Nil[T any] struct{}

func (n Nil[T]) Push(item T) List[T] {
	return Cons[T]{head: item, tail: n}
}

func (_ Nil[T]) Pop() (head T, tail List[T], ok bool) {
	ok = false
	return
}

func (_ Nil[T]) Head() (_ T, ok bool) {
	ok = false
	return
}

func (_ Nil[T]) Tail() List[T] {
	return nil
}

type Cons[T any] struct {
	head T
	tail List[T]
}

func (c Cons[T]) Push(item T) List[T] {
	return Cons[T]{head: item, tail: c}
}

func (c Cons[T]) Pop() (head T, tail List[T], ok bool) {
	return c.head, c.tail, true
}

func (c Cons[T]) Head() (_ T, ok bool) {
	return c.head, true
}

func (c Cons[T]) Tail() List[T] {
	return c.tail
}
