package types

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
)

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

	MinObserverDelegation = "1000000000000000000"
)

func GetMinObserverDelegation() (sdkmath.Int, bool) {
	return sdkmath.NewIntFromString(MinObserverDelegation)
}

func GetMinObserverDelegationDec() (sdkmath.LegacyDec, error) {
	return sdkmath.LegacyNewDecFromStr(MinObserverDelegation)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func BallotListKeyPrefix(p int64) []byte {
	return fmt.Appendf(nil, "%d", p)
}

func ChainNoncesKeyPrefix(chainID int64) []byte {
	return []byte(strconv.FormatInt(chainID, 10))
}

// Address TODOS for name changes: https://github.com/zeta-chain/node/issues/3098
const (
	BlameKey = "Blame-"
	// TODO change identifier for VoterKey to BallotKey
	VoterKey = "Voter-value-"

	// AllChainParamsKey is the ke prefix for all chain params
	// NOTE: CoreParams is old name for AllChainParams we keep it as key value for backward compatibility
	// TODO rename to ChainParamsListKey
	AllChainParamsKey = "CoreParams"

	ObserverSetKey = "ObserverSet-value-"

	// CrosschainFlagsKey is the key for the crosschain flags
	// NOTE: PermissionFlags is old name for CrosschainFlags we keep it as key value for backward compatibility
	CrosschainFlagsKey = "PermissionFlags-value-"

	// TODO rename to LastObserverCountKey
	LastBlockObserverCountKey = "ObserverCount-value-"
	NodeAccountKey            = "NodeAccount-value-"
	KeygenKey                 = "Keygen-value-"

	// TODO remove unused keys
	BlockHeaderKey      = "BlockHeader-value-"
	BlockHeaderStateKey = "BlockHeaderState-value-"

	// TODO rename to BallotListForHeightKey
	BallotListKey = "BallotList-value-"
	TSSKey        = "TSS-value-"
	TSSHistoryKey = "TSS-History-value-"

	// TODO rename to TssFundMigratorInfoKey
	TssFundMigratorKey = "FundsMigrator-value-"

	// TODO Rename to PendingNoncesListKey
	PendingNoncesKeyPrefix = "PendingNonces-value-"
	ChainNoncesKey         = "ChainNonces-value-"

	// TODO rename to NonceToCctxListKey
	NonceToCctxKeyPrefix = "NonceToCctx-value-"

	ParamsKey = "Params-value-"

	OperationalFlagsKey = "OperationalFlags-value-"
)

func GetBlameIndex(chainID int64, nonce uint64, digest string, height uint64) string {
	return fmt.Sprintf("%d-%d-%s-%d", chainID, nonce, digest, height)
}

func GetBlamePrefix(chainID int64, nonce int64) string {
	return fmt.Sprintf("%d-%d", chainID, nonce)
}
