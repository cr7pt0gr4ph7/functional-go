package option

type Optional[T any] struct {
	value    T
	hasValue bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, hasValue: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

func FromValueOrFalse[T any](value T, ok bool) Optional[T] {
	if ok {
		return Some(value)
	} else {
		return None[T]()
	}
}

func FromValueOrError[T any](value T, err error) Optional[T] {
	if err != nil {
		return Some(value)
	} else {
		return None[T]()
	}
}

func (o Optional[T]) Value() (value T, ok bool) {
	return o.value, o.hasValue
}

func (o Optional[T]) ValueOrDefault() T {
	if o.hasValue {
		return o.value
	} else {
		var fallback T
		return fallback
	}
}

func (o Optional[T]) ValueOrElse(fallback T) T {
	if o.hasValue {
		return o.value
	} else {
		return fallback
	}
}

func (o Optional[T]) ValueOrError(ifNone error) (T, error) {
	if v, ok := o.Value(); ok {
		return v, nil
	} else {
		// Reuse the default value of v created by Value()
		return v, ifNone
	}
}

func (o Optional[T]) ValueOrPanic(message string) T {
	if v, ok := o.Value(); ok {
		return v
	} else {
		panic(message)
	}
}

func (o Optional[_]) IsPresent() bool {
	return o.hasValue
}

func (o Optional[_]) IsEmpty() bool {
	return !o.hasValue
}

func (o Optional[_]) IsPresent() bool {
	return o.hasValue
}

// Function variant of `o.IsPresent()`.
// Useful as a predicate function.
func IsPresent[T any](o Optional[T]) bool {
	return o.hasValue
}

// Function variant of `o.IsEmpty()`.
// Useful as a predicate function
func IsEmpty[T any](o Optional[T]) bool {
	return !o.hasValue
}

func (o Optional[T]) OrElse(alternative Optional[T]) Optional[T] {
	if o.IsPresent() {
		return o
	} else {
		return alternative
	}
}

func Map[T any, R any](o Optional[T], mapping func(value T) R) Optional[R] {
	if v, ok := o.Value(); ok {
		return Some[R](mapping(v))
	} else {
		return None[R]()
	}
}

func FlatMap[T any, R any](o Optional[T], mapping func(value T) Optional[R]) Optional[R] {
	if v, ok := o.Value(); ok {
		return mapping(v)
	} else {
		return None[R]()
	}
}

func Flatten[T any](o Optional[Optional[T]]) Optional[T] {
	if v, ok := o.Value(); ok {
		return v
	} else {
		return None[T]()
	}
}

func FilterBy[T any](o Optional[T], predicate func(value T) bool) Optional[T] {
	if v, ok := o.Value(); ok && predicate(v) {
		return o
	} else {
		return None[T]()
	}
}
