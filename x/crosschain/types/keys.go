package types

import (
	"fmt"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "crosschain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	// NOTE: named metacore because previous module name was metacore
	// we keep the current value for backward compatibility
	MemStoreKey = "mem_metacore"
)

// GetProtocolFee returns the protocol fee for the cross-chain transaction
// It is no longer used, but the function is kept for backward compatibility with the Zeta Conversion Rate query
func GetProtocolFee() math.Uint {
	return math.ZeroUint()
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// CCTXKey is the key for the cross chain transaction
	// NOTE: Send is the previous name of CCTX and is kept for backward compatibility
	CCTXKey = "Send-value-"

	// CounterValueKey is a static key for storing the cctx counter key for ordering
	CounterValueKey = "ctr-value"

	// CounterIndexKey is the prefix to use for the counter index
	CounterIndexKey = "ctr-idx-"

	LastBlockHeightKey   = "LastBlockHeight-value-"
	FinalizedInboundsKey = "FinalizedInbounds-value-"

	GasPriceKey = "GasPrice-value-"

	// OutboundTrackerKeyPrefix is the prefix to retrieve all OutboundTracker
	// NOTE: OutTxTracker is the previous name of OutboundTracker and is kept for backward compatibility
	OutboundTrackerKeyPrefix = "OutTxTracker-value-"

	// InboundTrackerKeyPrefix is the prefix to retrieve all InboundHashToCctx
	// NOTE: InTxHashToCctx is the previous name of InboundHashToCctx and is kept for backward compatibility
	InboundTrackerKeyPrefix = "InTxTracker-value-"

	// ZetaAccountingKey value is used as prefix for storing ZetaAccountingKey
	// #nosec G101: Potential hardcoded credentials (gosec)
	ZetaAccountingKey = "ZetaAccounting-value-"

	RateLimiterFlagsKey = "RateLimiterFlags-value-"
)

// OutboundTrackerKey returns the store key to retrieve a OutboundTracker from the index fields
func OutboundTrackerKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

func FinalizedInboundKey(inboundHash string, chainID int64, eventIndex uint64) string {
	return fmt.Sprintf("%d-%s-%d", chainID, inboundHash, eventIndex)
}

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
	//ModuleAddressEVM common.EVMAddress
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	//0xB73C0Aac4C1E606C6E495d848196355e6CB30381
)
