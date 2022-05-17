package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ZetaConversionRateKeyPrefix is the prefix to retrieve all ZetaConversionRate
	ZetaConversionRateKeyPrefix = "ZetaConversionRate/value/"
)

// ZetaConversionRateKey returns the store key to retrieve a ZetaConversionRate from the index fields
func ZetaConversionRateKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
