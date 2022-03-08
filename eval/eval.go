package eval

type Eval[A any] struct {
	value A
}

func (e Eval[A]) Value() A {
	return e.value
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
