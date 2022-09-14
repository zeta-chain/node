package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "zetacore"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_metacore"

	ProtocolFee = 1000000000000000000
)

func GetProtocolFee() sdk.Uint {
	return sdk.NewUint(ProtocolFee)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TxinKey                     = "Txin-value-"
	TxinVoterKey                = "TxinVoter-value-"
	NodeAccountKey              = "NodeAccount-value-"
	TxoutKey                    = "Txout-value-"
	TxoutCountKey               = "Txout-count-"
	TxoutConfirmationKey        = "TxoutConfirmation-value-"
	SendKey                     = "Send-value-"
	VoteCounterKey              = "VoteCounter-value-"
	ReceiveKey                  = "Receive-value-"
	LastBlockHeightKey          = "LastBlockHeight-value-"
	ChainNoncesKey              = "ChainNonces-value-"
	GasPriceKey                 = "GasPrice-value-"
	SupportedChainsKey          = "SupportedChains-value-"
	AllSupportedChainsKey       = "AllSupportedChains"
	GasBalanceKey               = "GasBalance-value-"
	TxListKey                   = "TxList-value-"
	InTxKey                     = "InTx-value-"
	TSSKey                      = "TSS-value-"
	TSSVoterKey                 = "TSSVoter-value-"
	KeygenKey                   = "Keygen-value-"
	OutTxTrackerKeyPrefix       = "OutTxTracker/value/"
	ZetaConversionRateKeyPrefix = "ZetaConversionRate/value/"
)

// OutTxTrackerKey returns the store key to retrieve a OutTxTracker from the index fields
func OutTxTrackerKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// ZetaConversionRateKey returns the store key to retrieve a ZetaConversionRate from the index fields
func ZetaConversionRateKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// events follow here
const (
	SendEventKey         = "NewSendCreated" // Indicates what key to listen to
	StoredGameEventIndex = "Index"          // What game is relevant
	StoredGameEventRed   = "Red"            // Is it relevant to me?
	StoredGameEventBlack = "Black"          // Is it relevant to me?
)

func (cctx CrossChainTx) LogIdentifierForCCTX() string {
	return fmt.Sprintf("%s-%s-%s", cctx.InBoundTxParams.Sender, cctx.InBoundTxParams.SenderChain, cctx.OutBoundTxParams.ReceiverChain)
}
