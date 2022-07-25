package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OutTxTrackerKeyPrefix is the prefix to retrieve all OutTxTracker
	OutTxTrackerKeyPrefix = "OutTxTracker/value/"
)

// OutTxTrackerKey returns the store key to retrieve a OutTxTracker from the index fields
func OutTxTrackerKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
