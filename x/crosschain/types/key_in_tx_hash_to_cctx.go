package types

const (
	// InTxHashToCctxKeyPrefix is the prefix to retrieve all InTxHashToCctx
	InTxHashToCctxKeyPrefix = "InTxHashToCctx/value/"
)

// InTxHashToCctxKey returns the store key to retrieve a InTxHashToCctx from the index fields
func InTxHashToCctxKey(
	inTxHash string,
) []byte {
	var key []byte

	inTxHashBytes := []byte(inTxHash)
	key = append(key, inTxHashBytes...)
	key = append(key, []byte("/")...)

	return key
}
