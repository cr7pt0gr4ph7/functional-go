package lists

type ArrayList[T any] []T

func (l *ArrayList[T]) At(index int) T {
	return (*l)[index]
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

func (l *ArrayList[T]) Add(elem T) {
	*l = append(*l, elem)
}

func (l *ArrayList[T]) Clear() {
	*l = make([]T, 0)
}

func (l *ArrayList[T]) CopyTo(dst []T) {
	for i, t := range *l {
		dst[i] = t
	}
}
