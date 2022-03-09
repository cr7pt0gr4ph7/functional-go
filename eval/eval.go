package eval

type Eval[A any] struct {
	impl  evalImpl[A]
	value A
}

func (e Eval[A]) Memoize() Eval[A] {
	if e.impl != nil {
		return e.impl.Memoize()
	}
	return e
}

func (e Eval[A]) Value() A {
	if e.impl != nil {
		return e.impl.Value()
	}
	return e.value
}

func Defer[A any](deferred func() Eval[A]) Eval[A] {
	return fromImpl[A](&deferImpl[A]{run: deferred})
}

func Now[A any](value A) Eval[A] {
	return Eval[A]{value: value}
}

func Later[A any](valueFactory func() A) Eval[A] {
	return fromImpl[A](&laterImpl[A]{provider: valueFactory})
}

func Always[A any](valueFactory func() A) Eval[A] {
	return fromImpl[A](&alwaysImpl[A]{provider: valueFactory})
}

func Map[A any, B any](e Eval[A], f func(a A) B) Eval[B] {
	return FlatMap(e, func(a A) Eval[B] {
		return Now(f(a))
	})
}

func FlatMap[A any, B any](e Eval[A], f func(a A) Eval[B]) Eval[B] {
	return fromImpl[B](&flatMapImpl[A, B]{
		start: func() Eval[A] {
			return e
		},
		run: f,
	})
}
