package helpers

// SafeDereference dereferences and returns the value of a non-nil pointer.
// For a nil pointer, the zero-value of its type is returned.
func SafeDereference[T any](pointer *T) T {
	if pointer != nil {
		return *pointer
	} else {
		return *new(T)
	}
}
