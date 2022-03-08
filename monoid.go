package functional

type canAdd interface {
	constraints.Complex | constraints.Integer | constraints.Float | string
}

type nativeSemigroup[A canAdd] struct{}

func (s nativeSemigroup[A]) Combine(x A, y A) A {
	return x + y
}

func (s nativeSemigroup[A]) Empty() A {
	var x A
	return x // gives the right result for all allowed types
}

type Semigroup[A any] interface {
	Combine(x A, y A) A
}

type Monoid[A any] interface {
	Semigroup[A]
	Empty() A
}
