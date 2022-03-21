package views

type Indexed[T any] Keyed[int, T]

type Keyed[K any, V any] interface {
	At(key K) V
	AtOrDefault(key K) V
	AtOrElse(key K, fallback V) V
	TryAt(key K) (value V, ok bool)
	DefinedAt(key K) bool
}

type Sized interface {
	Empty() bool
	Len() int
}
