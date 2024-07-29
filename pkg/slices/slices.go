package slices

// Map applies a function to each item in a slice and returns a new slice with the results.
func Map[T, V any](items []T, f func(T) V) []V {
	result := make([]V, len(items))

	for i, item := range items {
		result[i] = f(item)
	}

	return result
}
