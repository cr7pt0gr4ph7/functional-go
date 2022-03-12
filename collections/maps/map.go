package maps

type Map[K comparable, V any] map[K]V

type Entry[K any, V any] struct {
	Key   K
	Value V
}

type Overwrite[K comparable, V any] struct {
	Map[K, V]
}

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

func (m Map[K, V]) Add(entry Entry[K, V]) bool {
	if !m.DefinedAt(entry.Key) {
		m[entry.Key] = entry.Value
		return true
	}
	return false
}

func (o Overwrite[K, V]) Add(entry Entry[K, V]) bool {
	o.Map.Set(entry)
	return true
}

func (m Map[K, V]) Set(entry Entry[K, V]) {
	m[entry.Key] = entry.Value
}

func (m Map[K, V]) Remove(entry Entry[K, V]) bool {
	_, existed := m[entry.Key]
	// TODO(lw) Compare v and entry.Value?
	delete(m, entry.Key)
	return existed
}

func (m Map[K, V]) InsertAt(key K, value V) {
	if !m.DefinedAt(key) {
		m[key] = value
	}
}

func (o Overwrite[K, V]) InsertAt(key K, value V) {
	o.Map.ReplaceAt(key, value)
}

func (m Map[K, V]) ReplaceAt(key K, value V) {
	m[key] = value
}

func (m Map[K, V]) RemoveAt(key K) bool {
	_, existed := m[key]
	delete(m, key)
	return existed
}
