package types

import (
	"bytes"
	"encoding/hex"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/zetaclient/config"
)

// InboundProcessability is an enum representing the processability of an inbound
type InboundProcessability int

const (
	// InboundProcessabilityGood represents a processable inbound
	InboundProcessabilityGood InboundProcessability = iota

	// InboundProcessabilityDonation represents a donation inbound
	InboundProcessabilityDonation

	// InboundProcessabilityComplianceViolation represents a compliance violation
	InboundProcessabilityComplianceViolation
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

// DecodeMemo decodes the receiver from the memo bytes
func (event *InboundEvent) DecodeMemo() error {
	// skip decoding donation tx as it won't go through zetacore
	if bytes.Equal(event.Memo, []byte(constant.DonationMessage)) {
		return nil
	}

	// decode receiver address from memo
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.Memo))
	if err != nil { // unreachable code
		return errors.Wrap(err, "invalid memo hex")
	}

	// ensure the receiver is valid
	if crypto.IsEmptyAddress(parsedAddress) {
		return errors.New("got empty receiver address from memo")
	}
	event.Receiver = parsedAddress.Hex()

	return nil
}

// Processability returns the processability of the inbound event
func (event *InboundEvent) Processability() InboundProcessability {
	// parse memo-specified receiver
	receiver := ""
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.Memo))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}

	// check restricted addresses
	if config.ContainRestrictedAddress(event.Sender, event.Receiver, event.TxOrigin, receiver) {
		return InboundProcessabilityComplianceViolation
	}

	// donation check
	if bytes.Equal(event.Memo, []byte(constant.DonationMessage)) {
		return InboundProcessabilityDonation
	}

	return InboundProcessabilityGood
}
