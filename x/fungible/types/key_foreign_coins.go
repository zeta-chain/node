package types

const (
	// ForeignCoinsKeyPrefix is the prefix to retrieve all ForeignCoins
	ForeignCoinsKeyPrefix = "ForeignCoins/value/"
)

// ForeignCoinsKey returns the store key to retrieve a ForeignCoins from the index fields
func ForeignCoinsKey(
	index string,
) []byte {
	indexBytes := []byte(index)
	key := make([]byte, 0, len(indexBytes)+1)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
