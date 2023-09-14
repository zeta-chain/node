package types

import "fmt"

const (
	// ModuleName defines the module name
	ModuleName = "observer"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_observer"

	GroupID1Address = "zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func BallotListKeyPrefix(p int64) []byte {
	return []byte(fmt.Sprintf("%d", p))
}

const (
	BlameKey = "Blame-"
	// TODO change identifier for VoterKey to something more descriptive
	VoterKey                      = "Voter-value-"
	AllCoreParams                 = "CoreParams"
	ObserverMapperKey             = "Observer-value-"
	ObserverParamsKey             = "ObserverParams"
	AdminPolicyParamsKey          = "AdminParams"
	BallotMaturityBlocksParamsKey = "BallotMaturityBlocksParams"

	PermissionFlagsKey        = "PermissionFlags-value-"
	LastBlockObserverCountKey = "ObserverCount-value-"
	NodeAccountKey            = "NodeAccount-value-"
	KeygenKey                 = "Keygen-value-"

	BallotListKey = "BallotList-value-"
)
