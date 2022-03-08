package eval

func Defer[A any](deferred func() Eval[A]) Eval[A] {
	return deferred()
}

func Now[A any](value A) Eval[A] {
	return Eval[A]{value: value}
}

func Later[A any](valueFactory func() A) Eval[A] {
	return Eval[A]{value: valueFactory()}
}

func Always[A any](valueFactory func() A) Eval[A] {
	return Eval[A]{value: valueFactory()}
}

func Map[A any, B any](e Eval[A], f func(a A) B) Eval[B] {
	return Later(func() B {
		return f(e.Value())
	})
}

func FlatMap[A any, B any](e Eval[A], f func(a A) Eval[B]) Eval[B] {
	return Defer(func() Eval[B] {
		return f(e.Value())
	})
}
