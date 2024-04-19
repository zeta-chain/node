package types

const (
	// ModuleName defines the module name
	ModuleName = "lightclient"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName
)

const (
	BlockHeaderKey       = "BlockHeader-value-"
	ChainStateKey        = "ChainState-value-"
	VerificationFlagsKey = "VerificationFlags-value-"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
