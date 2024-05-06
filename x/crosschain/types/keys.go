package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	MemStoreKey = "mem_metacore"

	ProtocolFee = 2000000000000000000

	//TssMigrationGasMultiplierEVM is multiplied to the median gas price to get the gas price for the tss migration . This is done to avoid the tss migration tx getting stuck in the mempool
	TssMigrationGasMultiplierEVM = "2.5"

	ZetaIndexLength = 66
)

func GetProtocolFee() sdk.Uint {
	return sdk.NewUint(ProtocolFee)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// CCTXKey is the key for the cross chain transaction
	// NOTE: Send is the previous name of CCTX and is kept for backward compatibility
	CCTXKey = "Send-value-"

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

func (m CrossChainTx) LogIdentifierForCCTX() string {
	if len(m.OutboundParams) == 0 {
		return fmt.Sprintf("%s-%d", m.InboundParams.Sender, m.InboundParams.SenderChainId)
	}
	i := len(m.OutboundParams) - 1
	outTx := m.OutboundParams[i]
	return fmt.Sprintf("%s-%d-%d-%d", m.InboundParams.Sender, m.InboundParams.SenderChainId, outTx.ReceiverChainId, outTx.TssNonce)
}

func FinalizedInboundKey(intxHash string, chainID int64, eventIndex uint64) string {
	return fmt.Sprintf("%d-%s-%d", chainID, intxHash, eventIndex)
}

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
	//ModuleAddressEVM common.EVMAddress
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	//0xB73C0Aac4C1E606C6E495d848196355e6CB30381
)
