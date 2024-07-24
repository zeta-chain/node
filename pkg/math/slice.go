package math

import "slices"

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// SliceMedianValue returns the median value of the given slice.
// Returns 0 for an empty slice. If inPlace is true, the input slice will be sorted in place.
// Otherwise, a copy of the input slice will be sorted.
func SliceMedianValue[T number](items []T, inPlace bool) T {
	if inPlace {
		return sliceMedianValue(items)
	}

	copied := make([]T, len(items))
	copy(copied, items)

	out := sliceMedianValue(copied)

	// We don't need the copy anymore
	//nolint:ineffassign // let's help the garbage collector :)
	copied = nil

	return out
}

func sliceMedianValue[T number](items []T) T {
	switch len(items) {
	case 0:
		return 0
	case 1:
		return items[0]
	}

	slices.Sort(items)

	// note that int division is used here e.g. 5/2 => 2

	// []int{1 2 3 4 5} => items[(5/2)] => items[2] => 3
	if len(items)%2 == 1 {
		return items[len(items)/2]
	}

	// odd number of items
	rightIndex := len(items) / 2
	leftIndex := rightIndex - 1

	// []int{1 2 3 4} => (items[1] + items[2]) / 2 => 5/2 => 2
	return (items[leftIndex] + items[rightIndex]) / 2
}
