package option

func Lift0[A any, B any](f func() B) func(_ Optional[A]) Optional[B] {
	return func(a Optional[A]) Optional[B] {
		if a.hasValue {
			return Some(f())
		}
		return None[B]()
	}
}

func Lift1[A any, B any](f func(a A) B) func(a Optional[A]) Optional[B] {
	return func(a Optional[A]) Optional[B] {
		if a.hasValue {
			return Some(f(a.value))
		}
		return None[B]()
	}
}

func Lift2[A any, B any, C any](f func(a A, b B) C) func(a Optional[A], b Optional[B]) Optional[C] {
	return func(a Optional[A], b Optional[B]) Optional[C] {
		if a.hasValue && b.hasValue {
			return Some(f(a.value, b.value))
		}
		return None[C]()
	}

}

func Lift3[A any, B any, C any, D any](f func(a A, b B, c C) D) func(a Optional[A], b Optional[B], c Optional[C]) Optional[D] {
	return func(a Optional[A], b Optional[B], c Optional[C]) Optional[D] {
		if a.hasValue && b.hasValue && c.hasValue {
			return Some(f(a.value, b.value, c.value))
		}
		return None[D]()
	}
}
