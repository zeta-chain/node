package types

const (
	// ModuleName defines the module name
	ModuleName = "authority"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName
)

// KeyPrefix returns the store key prefix for the policies store
func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// PoliciesKey is the key for the policies store
	PoliciesKey = "Policies-value-"
)
