package hmap

// ==========
// :: Maps ::
// ==========

type UntypedEntry = TypeProof

// Represents a heterogeneous map.
type Map struct {
	items map[keyInstance]any
}

// Wrapper type to make casts more strict.
// Relevant for nil and interface types.
type valueHolder[T any] struct {
	value T
}

func Get[T any](from Map, key Key[T]) (item T, ok bool) {
	if x, ok := from.items[key.actualKey()]; ok {
		return x.(valueHolder[T]).value, true
	}
	return
}

func Put[T any](into Map, key Key[T], value T) {
	into.items[key.actualKey()] = valueHolder[T]{value}
}

func Remove[T any](from Map, key Key[T]) {
	delete(from.items, key.actualKey())
}

func (m Map) Clear() {
	for k, _ := range m.items {
		delete(m.items, k)
	}
}

func (m Map) Clone() Map {
	var clone Map
	for k, v := range m.items {
		clone.items[k] = v
	}
	return clone
}

// ==========
// :: Keys ::
// ==========

type keyInstance any

type KeyBase interface {
	actualKey() keyInstance
	isType(value any) bool
	checkType(value any)
	getFrom(from Map) (any, bool)
	putInto(into Map, value any)
}

type Key[T any] interface {
	KeyBase
	Get(from Map) (T, bool)
	Put(into Map, value T)
}

type keyImpl[T any] struct {
	name string
}

func NewKey[T any](name string) Key[T] {
	return &keyImpl[T]{name: name}
}

//
// KeyBase methods
//

func (k *keyImpl[T]) actualKey() keyInstance {
	return k
}

func (k *keyImpl[T]) isType(value any) bool {
	_, ok := value.(T)
	return ok
}

func (k *keyImpl[T]) checkType(value any) {
	_ = value.(T)
}

func (k *keyImpl[T]) getFrom(from Map) (any, bool) {
	return k.Get(from)
}

func (k *keyImpl[T]) putInto(into Map, value any) {
	k.Put(into, value.(T))
}

func (k *keyImpl[T]) String() string {
	return k.name
}

//
// Key[T] methods
//

func (k *keyImpl[T]) Get(from Map) (T, bool) {
	return Get[T](from, k)
}

func (k *keyImpl[T]) Put(into Map, value T) {
	Put[T](into, k, value)
}

// ============
// :: Proofs ::
// ============

// A valid proof can only be constructed from
// a `Key[T]` and a value of type `T`.
type TypeProof struct {
	isValid bool
	key     KeyBase
	value   any
}

func NewProof[T any](key Key[T], value T) TypeProof {
	return newProofUnchecked(key, value)
}

func NewProofChecked[T any](key KeyBase, value any) TypeProof {
	key.checkType(value)
	return newProofUnchecked(key, value)
}

func newProofUnchecked(key KeyBase, value any) TypeProof {
	return TypeProof{
		isValid: true,
		key:     key,
		value:   value,
	}
}

func (p TypeProof) IsValid() bool {
	return p.isValid
}

func (p TypeProof) CheckValid() {
	if !p.isValid {
		panic("not a valid type proof instance")
	}
}

func (p TypeProof) Extract() (KeyBase, any) {
	p.CheckValid()
	return p.key, p.value
}
