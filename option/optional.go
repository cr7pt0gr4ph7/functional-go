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

func (o Optional[T]) Get() (value T, ok bool) {
	return o.value, o.hasValue
}

func (o Optional[T]) GetOrDefault() T {
	if o.hasValue {
		return o.value
	} else {
		var fallback T
		return fallback
	}
}

func (o Optional[T]) GetOrFallback(fallback T) T {
	if o.hasValue {
		return o.value
	} else {
		return fallback
	}
}

func (o Optional[_]) IsPresent() bool {
	return o.hasValue
}

func (o Optional[_]) IsEmpty() bool {
	return !o.hasValue
}
