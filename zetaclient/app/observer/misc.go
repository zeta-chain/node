package observer

// returns the maximum of two ints
func MaxInt(a int, b int) int {
	if a < b {
		return b
	}
	return a
}
