// Package ptr provides helper functions for working with pointers.
package ptr

// Ptr returns a pointer to the value passed in.
func Ptr[T any](value T) *T {
	return &value
}

// Deref returns the value of the pointer passed in, or the zero value of the type if the pointer is nil.
func Deref[T any](value *T) T {
	var out T
	if value != nil {
		out = *value
	}

	return out
}
