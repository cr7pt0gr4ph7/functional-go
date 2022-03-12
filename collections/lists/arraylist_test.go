package lists

import (
	"fmt"
	"github.com/cr7pt0gr4ph7/functional-go/collections/views"
)

// Ensure that all expected interfaces are implemented.
func _[T any]() {
	var l ArrayList[T]
	var _ views.Indexed[T] = l
	var _ views.Keyed[int, T] = l
	var _ views.Sized = l
}

func ExampleArrayList() {
	var l ArrayList[int]
	fmt.Println(l)

	l.Add(10)
	l.Add(20)
	l.Add(30)
	fmt.Println(l)

	l.Clear()
	fmt.Println(l)
}
