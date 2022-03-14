package effects

type Eff[E any, T any] struct {
	Value T
}

func Lift[E any, A any, B any](f func(arg A) B) func(arg Eff[E, A]) Eff[E, B] {
	return func(arg Eff[E, A]) Eff[E, B] {
		return Return[E](f(arg.Value))
	}
}

func Return[E any, T any](value T) Eff[E, T] {
	return Eff[E, T]{Value: value}
}

func Map[E any, A any, B any](arg Eff[E, A], f func(arg A) B) Eff[E, B] {
	return Return[E](f(arg.Value))
}

func FlatMap[E any, A any, B any](arg Eff[E, A], f func(arg A) Eff[E, B]) Eff[E, B] {
	return f(arg.Value)
}
