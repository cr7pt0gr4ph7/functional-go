package builders

type Indexed[T any] Keyed[int, T]

type Keyed[K any, V any] interface {
	InsertAt(key K, value V)
	ReplaceAt(key K, value V)
	RemoveAt(key K) bool
	Clear()
}

type Unordered[T any] interface {
	Add(elem T) bool
	Remove(elem T) bool
	Clear()
}
