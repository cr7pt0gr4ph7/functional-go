package maps

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) At(key K) V {
	return m[key]
}

func (m Map[K, V]) TryAt(key K) (value V, ok bool) {
	value, ok = m[key]
	return
}

func (m Map[K, V]) AtOrDefault(key K) V {
	v, _ := m[key]
	return v
}

func (m Map[K, V]) AtOrElse(key K, fallback V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}

func (m Map[K, _]) DefinedAt(key K) bool {
	_, ok := m[key]
	return ok
}

func (m Map[_, _]) Empty() bool {
	return len(m) == 0
}

func (m Map[_, _]) Len() int {
	return len(m)
}
