// Package typelists allows representing list of types on the type level.
package typelists

// A type whose existence only matters for type-level operations,
// and whose value-level representation is simply empty.
type phantomType = struct{]

// TypeList can be used as a constraint for type-level lists.
type TypeList interface{
 	typeList() // marker method
}

// NonEmpty can be used as a constraint for non-empty type-level lists.
type NonEmpty interface {
	TypeList
	nonEmpty() // marker method
}

func (_ Nil) typeList()        {}
func (_ Cons[_, _]) typeList() {}
func (_ Cons[_, _]) nonEmpty() {}

// Nil represents the empty list.
type Nil phantomType

// Cons represents a list with a head element and a tail list.
type Cons[Head any, Tail TypeList] phantomType

func (_ Cons[Head, _]) Head() Head { return *new(Head) }
func (_ Cons[_, Tail]) Tail() Tail { return *new(Tail) }
