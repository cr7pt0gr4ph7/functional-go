package effects

type Eff[E any, T any] interface {
	eff(eff E) T // marker method - do not call
}

func (_ Pure[E, T]) eff(eff E) T    { panic("marker method") }
func (_ Cont[E, _, T]) eff(eff E) T { panic("marker method") }

type Pure[E any, T any] struct {
	Value T
}

type Cont[E any, A any, B any] struct {
}

func Lift[E any, A any, B any](f func(arg A) B) func(arg Eff[E, A]) Eff[E, B] {
	return func(arg Eff[E, A]) Eff[E, B] {
		return Return[E](f(arg.(Pure[E, A]).Value))
	}
}

func Return[E any, T any](value T) Eff[E, T] {
	return Pure[E, T]{Value: value}
}

func Map[E any, A any, B any](arg Eff[E, A], f func(arg A) B) Eff[E, B] {
	return Return[E](f(arg.(Pure[E, A]).Value))
}

func FlatMap[E any, A any, B any](arg Eff[E, A], f func(arg A) Eff[E, B]) Eff[E, B] {
	return f(arg.(Pure[E, A]).Value)
}
