package lists

type ArrayList[T any] []T

func (l *ArrayList[T]) At(index int) T {
	return (*l)[index]
}

func (l *ArrayList[T]) TryAt(index int) (T, bool) {
	if l.DefinedAt(index) {
		return l.At(index), true
	}
	var defaultT T
	return defaultT, false
}

func (l *ArrayList[T]) AtOrDefault(index int) T {
	if l.DefinedAt(index) {
		return l.At(index)
	}
	var defaultT T
	return defaultT
}

func (l *ArrayList[T]) AtOrElse(index int, fallback T) T {
	if l.DefinedAt(index) {
		return l.At(index)
	}
	return fallback
}

func (l *ArrayList[_]) DefinedAt(index int) bool {
	return index >= 0 && index < len(*l)
}

func (l *ArrayList[_]) IsEmpty() bool {
	return len(*l) == 0
}

func (l *ArrayList[_]) Len() int {
	return len(*l)
}

func (l *ArrayList[T]) Add(elem T) bool {
	*l = append(*l, elem)
	return true
}

func (l *ArrayList[T]) Clear() {
	*l = make([]T, 0)
}

func (l *ArrayList[T]) CopyTo(dst []T) {
	for i, t := range *l {
		dst[i] = t
	}
}
