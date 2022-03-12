package maps

// Ensure that all expected interfaces are implemented.
func _[K comparable, V any]() {
	var m Map[K, V]
	var _ Keyed[K, V] = m
	var _ Sized = m
}
