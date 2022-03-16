package immutable

import (
	"github.com/cr7pt0gr4ph7/functional-go/collections/immutable/cursor"
)

// Represents an immutable list.
type List[L List_[L, T], T any] interface {
	List_[L, T]
}

// Defines the methods for `List[L, T]`.
//
// This is a separate type due to compiler restrictions.
// This type is public so that it is visible in the documentation.
type List_[L any, T any] interface {
	Prepend(item T) L
	Append(item T) L
	Concat(other L) L
	Empty() bool
	Len() int
	Cursor() cursor.Cursor[T]
}

// NOTE: Generic type aliases are not possible (as of Go 1.18),
//       therefore we cannot provide the following alias:
// type Cursor[T any] = cursor.Cursor[T]
