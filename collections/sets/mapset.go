package sets

type MapSet[K comparable, V any] map[K]V

func (m MapSet[_, V]) placeholder() (value V) {
	return
}

func (m MapSet[K, _]) DefinedAt(key K) bool {
	_, ok := m[key]
	return ok
}

func (m MapSet[_, _]) Empty() bool {
	return len(m) == 0
}

func (m MapSet[_, _]) Len() int {
	return len(m)
}

func (m MapSet[K, V]) Add(elem K) bool {
	if _, ok := m[elem]; !ok {
		m[elem] = m.placeholder()
		return true
	}
	return false
}

func (m MapSet[K, V]) Remove(elem K) bool {
	if _, ok := m[elem]; ok {
		delete(m, elem)
		return true
	}
	return false
}

func (m MapSet[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
}
