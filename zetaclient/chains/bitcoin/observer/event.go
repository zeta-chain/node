package observer

import (
	"bytes"

	"github.com/zeta-chain/node/pkg/constant"
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

// BTCInboundEvent represents an incoming transaction event
type BTCInboundEvent struct {
	// FromAddress is the first input address
	FromAddress string

	// ToAddress is the TSS address
	ToAddress string

	// Value is the amount of BTC
	Value float64

	// DepositorFee is the deposit fee
	DepositorFee float64

	// MemoBytes is the memo of inbound
	MemoBytes []byte

	// MemoStd is the standard inbound memo if it can be decoded
	MemoStd *memo.InboundMemo

	// BlockNumber is the block number of the inbound
	BlockNumber uint64

	// TxHash is the hash of the inbound
	TxHash string
}

// IsProcessable checks if the inbound event is processable
func (event *BTCInboundEvent) CheckProcessability() InboundProcessability {
	// compliance check on sender and receiver addresses
	if config.ContainRestrictedAddress(event.FromAddress, event.ToAddress) {
		return InboundProcessabilityComplianceViolation
	}

	// compliance check on receiver, revert/abort addresses in standard memo
	if event.MemoStd != nil {
		if config.ContainRestrictedAddress(
			event.MemoStd.Receiver.Hex(),
			event.MemoStd.RevertOptions.RevertAddress,
			event.MemoStd.RevertOptions.AbortAddress,
		) {
			return InboundProcessabilityComplianceViolation
		}
	}

	// donation check
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return InboundProcessabilityDonation
	}

	return InboundProcessabilityGood
}
