// Package typelists allows representing list of types on the type level.
package typelists

// TypeList can be used as a constraint for type-level lists.
type TypeList interface{ typeList() }

func (_ Nil) typeList()        {}
func (_ Cons[_, _]) typeList() {}

// Nil represents the empty list.
type Nil struct{}

// Cons represents a list with a head element and a tail list.
type Cons[Head any, Tail TypeList] struct{}

func (_ Cons[Head, _]) Head() Head { return *new(Head) }
func (_ Cons[_, Tail]) Tail() Tail { return *new(Tail) }
