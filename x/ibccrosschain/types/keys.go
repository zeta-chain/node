package types

const (
	// ModuleName defines the module name
	// NOTE: module name can't have the name of another module as a prefix
	// because of potential store key conflicts
	// ibcblockchain or crosschainibc can't be used as module name
	ModuleName = "zetaibccrosschain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName
)
