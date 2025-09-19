package observer

import (
	"bytes"
	"encoding/hex"
	"fmt"

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
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
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

	// Status is the status of the inbound event
	Status crosschaintypes.InboundStatus

	// ErrorMessage carries error information that caused non-SUCCESS 'Status'
	ErrorMessage string
}

// SetStatusAndErrMessage attaches the status and error message to the inbound event
func (event *BTCInboundEvent) SetStatusAndErrMessage(status crosschaintypes.InboundStatus, errorMessage string) {
	event.Status = status
	event.ErrorMessage = errorMessage
}

// Category returns the category of the inbound event
func (event *BTCInboundEvent) Category() clienttypes.InboundCategory {
	// compliance check on sender and receiver addresses
	if config.ContainRestrictedAddress(event.FromAddress, event.ToAddress) {
		return clienttypes.InboundCategoryRestricted
	}

	// compliance check on receiver, revert/abort addresses in standard memo
	if event.MemoStd != nil {
		if config.ContainRestrictedAddress(
			event.MemoStd.Receiver.Hex(),
			event.MemoStd.RevertOptions.RevertAddress,
			event.MemoStd.RevertOptions.AbortAddress,
		) {
			return clienttypes.InboundCategoryRestricted
		}
	}

	// donation check
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return clienttypes.InboundCategoryDonation
	}

	return clienttypes.InboundCategoryProcessable
}

// DecodeMemoBytes decodes the contained memo bytes as either standard or legacy memo
// It updates the event object with the decoded data
func (event *BTCInboundEvent) DecodeMemoBytes(chainID int64) error {
	var (
		err            error
		isStandardMemo bool
		memoStd        *memo.InboundMemo
		receiver       ethcommon.Address
	)

	// skip decoding if no memo is found, returning error to revert the inbound
	if bytes.Equal(event.MemoBytes, []byte(noMemoFound)) {
		event.MemoBytes = []byte{}
		return errors.New("no memo found in inbound")
	}

	// skip decoding donation tx as it won't go through zetacore
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return nil
	}

	// try to decode the standard memo as the preferred format
	// then process standard memo or fallback to legacy memo
	// note: err is guaranteed to be nil when 'isStandardMemo == false',
	// so a non-nil error indicates the standard memo contains improper data
	memoStd, isStandardMemo, err = memo.DecodeFromBytes(event.MemoBytes)
	if err != nil {
		return errors.Wrap(err, "standard memo contains improper data")
	}

	if isStandardMemo {
		// validate the content of the standard memo
		err = ValidateStandardMemo(*memoStd, chainID)
		if err != nil {
			return errors.Wrap(err, "invalid standard memo for bitcoin")
		}

		event.MemoStd = memoStd
		receiver = memoStd.Receiver
	} else {
		// legacy memo, ensure the it is no less than ZEVM address length (20-byte receiver)
		// checking upfront is to return more informative error message in the CCTX struct
		if len(event.MemoBytes) < ethcommon.AddressLength {
			return errors.New("legacy memo length must be at least 20 bytes")
		}

		parsedAddress, payload, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.MemoBytes))
		if err != nil { // unreachable code
			return errors.Wrap(err, "invalid legacy memo")
		}
		receiver = parsedAddress

		// update the memo bytes to only contain the data
		event.MemoBytes = payload
	}

	// ensure the receiver is valid
	if crypto.IsEmptyAddress(receiver) {
		return errors.New("got empty receiver address from memo")
	}
	event.ToAddress = receiver.Hex()

	return nil
}

// ValidateStandardMemo validates the standard memo in Bitcoin context
func ValidateStandardMemo(memoStd memo.InboundMemo, chainID int64) error {
	// ensure the revert address is a valid and supported BTC address
	revertAddress := memoStd.RevertOptions.RevertAddress
	if revertAddress != "" {
		btcAddress, err := chains.DecodeBtcAddress(revertAddress, chainID)
		if err != nil {
			return errors.Wrapf(err, "invalid revert address in memo: %s", revertAddress)
		}
		if !chains.IsBtcAddressSupported(btcAddress) {
			return fmt.Errorf("unsupported revert address in memo: %s", revertAddress)
		}
	}

	return nil
}

// IsEventProcessable checks if the inbound event is processable
func (ob *Observer) IsEventProcessable(event BTCInboundEvent) bool {
	logger := ob.Logger().Inbound.With().Str(logs.FieldTx, event.TxHash).Logger()

	switch category := event.Category(); category {
	case clienttypes.InboundCategoryProcessable:
		return true
	case clienttypes.InboundCategoryDonation:
		logger.Info().Msg("thank you rich folk for your donation!")
		return false
	case clienttypes.InboundCategoryRestricted:
		coinType := coin.CoinType_Gas
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance, false,
			ob.Chain().ChainId, event.TxHash, event.FromAddress, event.ToAddress, &coinType)
		return false
	default:
		logger.Error().Any("category", category).Msg("unreachable code got InboundCategory")
		return false
	}
}
