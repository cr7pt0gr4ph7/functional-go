package maps

import (
	"github.com/cr7pt0gr4ph7/functional-go/collections/views"
)

// Ensure that all expected interfaces are implemented.
func _[K comparable, V any]() {
	var m Map[K, V]
	var _ views.Keyed[K, V] = m
	var _ views.Sized = m
}
