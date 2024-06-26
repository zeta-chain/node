package types

const (
	// InboundHashToCctxKeyPrefix is the prefix to retrieve all InboundHashToCctx
	// NOTE: InTxHashToCctx is the previous name of InboundHashToCctx and is kept for backward compatibility
	InboundHashToCctxKeyPrefix = "InTxHashToCctx/value/"
)

// InboundHashToCctxKey returns the store key to retrieve a InboundHashToCctx from the index fields
func InboundHashToCctxKey(
	inboundHash string,
) []byte {
	var key []byte

	inboundHashBytes := []byte(inboundHash)
	key = append(key, inboundHashBytes...)
	key = append(key, []byte("/")...)

	return key
}
