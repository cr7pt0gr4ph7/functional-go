package collections

type Sized interface {
	Empty() bool
	Len() int
}

type ReadOnlyKeyed[K any, V any] interface {
	At(key K) V
	AtOrDefault(key K) V
	AtOrElse(key K, fallback V) V
	TryAt(key K) (value V, ok bool)
	DefinedAt(key K) bool
}

type ReadOnlyIndexed[T any] ReadOnlyKeyed[int, T]

// ===================
// Mutable collections
// ===================

// Represents a mutable collection with the ability to add elements.
type Builder[T any] interface {
	Add(elem T) bool
}

// Represents a mutable collection or set.
type Unordered[T any] interface {
	Add(elem T) bool    // Add `elem` to the collection. Returns `true` when the element was actually added.
	Remove(elem T) bool // Removes `elem` from the collection. Returns `true` when the element was actually removed.
	Clear()             // Removes all elements from the collection.
}

type Keyed[K any, V any] interface {
	InsertAt(key K, value V)
	ReplaceAt(key K, value V)
	RemoveAt(key K) bool
	Clear()
}

type Indexed[T any] Keyed[int, T]
