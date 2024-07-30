package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Map applies a function to each item in a slice and returns a new slice with the results.
func Map[T, V any](items []T, f func(T) V) []V {
	result := make([]V, len(items))

	for i, item := range items {
		result[i] = f(item)
	}

	return result
}

// ElementsMatch returns true if two slices have the same elements in the same order.
// Note that this function SORTS the slices before comparing them.
func ElementsMatch[T constraints.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Diff returns the elements in `a` that are not in `b`
func Diff[T comparable](a, b []T) []T {
	var (
		cache  = map[T]struct{}{}
		result []T
	)

	for _, v := range b {
		cache[v] = struct{}{}
	}

	for _, v := range a {
		if _, ok := cache[v]; !ok {
			result = append(result, v)
		}
	}

	return result
}
