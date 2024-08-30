package types

import (
	"github.com/zeta-chain/node/pkg/coin"
)

// InboundEvent represents an inbound event
// TODO: we should consider using this generic struct when it applies (e.g. for Bitcoin, Solana, etc.)
// https://github.com/zeta-chain/node/issues/2495
type InboundEvent struct {
	// SenderChainID is the chain ID of the sender
	SenderChainID int64

	// Sender is the sender address
	Sender string

	// Receiver is the receiver address
	Receiver string

	// TxOrigin is the origin of the transaction
	TxOrigin string

	// Value is the amount of token
	Amount uint64

	// Memo is the memo attached to the inbound
	Memo []byte

	// BlockNumber is the block number of the inbound
	BlockNumber uint64

	// TxHash is the hash of the inbound
	TxHash string

	// Index is the index of the event
	Index uint32

	// CoinType is the coin type of the inbound
	CoinType coin.CoinType

	// Asset is the asset of the inbound
	Asset string
}
