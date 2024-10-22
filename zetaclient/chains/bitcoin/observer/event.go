package observer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
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

	// ToAddress is the ZEVM receiver address
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

// DecodeEventMemoBytes decodes the memo bytes as either standard or legacy memo
func (ob *Observer) DecodeEventMemoBytes(event *BTCInboundEvent) error {
	// skip decoding donation tx as it won't go through zetacore
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return nil
	}

	// try to decode the standard memo as the preferred format
	var receiver ethcommon.Address
	memoStd, err := memo.DecodeFromBytes(event.MemoBytes)
	switch {
	case err == nil:
		// NoAssetCall will be disabled for Bitcoin until full V2 support
		if memoStd.OpCode == memo.OpCodeCall {
			return errors.New("NoAssetCall is disabled")
		}

		// ensure the revert address is valid and supported
		revertAddress := memoStd.RevertOptions.RevertAddress
		if revertAddress != "" {
			btcAddress, err := chains.DecodeBtcAddress(revertAddress, ob.Chain().ChainId)
			if err != nil {
				return errors.Wrapf(err, "invalid revert address in memo: %s", revertAddress)
			}
			if !chains.IsBtcAddressSupported(btcAddress) {
				return fmt.Errorf("unsupported revert address in memo: %s", revertAddress)
			}
		}
		event.MemoStd = memoStd
		receiver = memoStd.Receiver
	default:
		// fallback to legacy memo if standard memo decoding fails
		parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.MemoBytes))
		if err != nil { // never happens
			return errors.Wrap(err, "invalid legacy memo")
		}
		receiver = parsedAddress
	}

	// ensure the receiver is valid
	if crypto.IsEmptyAddress(receiver) {
		return errors.New("got empty receiver address from memo")
	}
	event.ToAddress = receiver.Hex()

	return nil
}

// CheckEventProcessability checks if the inbound event is processable
func (ob *Observer) CheckEventProcessability(event *BTCInboundEvent) bool {
	// check if the event is processable
	result := event.CheckProcessability()
	switch result {
	case InboundProcessabilityGood:
		return true
	case InboundProcessabilityDonation:
		logFields := map[string]any{
			logs.FieldChain: ob.Chain().ChainId,
			logs.FieldTx:    event.TxHash,
		}
		ob.Logger().Inbound.Info().Fields(logFields).Msgf("thank you rich folk for your donation!")
		return false
	case InboundProcessabilityComplianceViolation:
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.Chain().ChainId, event.TxHash, event.FromAddress, event.ToAddress, "BTC")
		return false
	}

	return true
}

// NewInboundVoteV1 creates a MsgVoteInbound message for inbound that uses legacy memo
func (ob *Observer) NewInboundVoteV1(event *BTCInboundEvent, amountSats *big.Int) *crosschaintypes.MsgVoteInbound {
	message := hex.EncodeToString(event.MemoBytes)

	return crosschaintypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		event.FromAddress,
		ob.Chain().ChainId,
		event.FromAddress,
		event.ToAddress,
		ob.ZetacoreClient().Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountSats),
		message,
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V1,
		false, // not relevant for v1
	)
}

// NewInboundVoteMemoStd creates a MsgVoteInbound message for inbound that uses standard memo
// TODO: rename this function as 'EventToInboundVoteV2' to use ProtocolContractVersion_V2 and enable more options
// https://github.com/zeta-chain/node/issues/2711
func (ob *Observer) NewInboundVoteMemoStd(event *BTCInboundEvent, amountSats *big.Int) *crosschaintypes.MsgVoteInbound {
	// replace 'sender' with 'revertAddress' if specified in the memo, so that
	// zetacore will refund to the address specified by the user in the revert options.
	sender := event.FromAddress
	if event.MemoStd.RevertOptions.RevertAddress != "" {
		sender = event.MemoStd.RevertOptions.RevertAddress
	}

	// make a legacy message so that zetacore can process it as V1
	msgBytes := append(event.MemoStd.Receiver.Bytes(), event.MemoStd.Payload...)
	message := hex.EncodeToString(msgBytes)

	return crosschaintypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		sender,
		ob.Chain().ChainId,
		event.FromAddress,
		event.ToAddress,
		ob.ZetacoreClient().Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountSats),
		message,
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V1,
		false, // not relevant for v1
	)
}
