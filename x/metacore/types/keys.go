package types

const (
	// ModuleName defines the module name
	ModuleName = "metacore"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_metacore"

	// this line is used by starport scaffolding # ibc/keys/name
)

// this line is used by starport scaffolding # ibc/keys/port

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TxinKey = "Txin-value-"
)

const (
	TxinVoterKey = "TxinVoter-value-"
)

const (
	NodeAccountKey = "NodeAccount-value-"
)

const (
	TxoutKey      = "Txout-value-"
	TxoutCountKey = "Txout-count-"
)

const (
	TxoutConfirmationKey = "TxoutConfirmation-value-"
)

const (
	SendVoterKey = "SendVoter-value-"
)

const (
	SendKey = "Send-value-"
)
