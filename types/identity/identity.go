package identity

type Identity[T any] struct {
	Value T
}

func Lift[A any, B any](f func(arg A) B) func(arg Identity[A]) Identity[B] {
	return func(arg Identity[A]) Identity[B] {
		return Return(f(arg.Value))
	}
}

func Return[T any](value T) Identity[T] {
	return Identity[T]{Value: value}
}

func Map[A any, B any](arg Identity[A], f func(arg A) B) Identity[B] {
	return Return(f(arg.Value))
}

func FlatMap[A any, B any](arg Identity[A], f func(arg A) Identity[B]) Identity[B] {
	return f(arg.Value)
}
