package signer

import (
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	contract "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/compliance"
)

// CancelReason represents a reason for canceling an outbound (increase_seqno)
type CancelReason uint32

const (
	// InvalidWorkchain tx cancelled due to attempt to withdraw to masterchain account
	// https://docs.ton.org/v3/documentation/smart-contracts/shards/shards-intro
	InvalidWorkchain CancelReason = 1

	// ComplianceViolation tx cancelled due to attempt to withdraw to a restricted address
	ComplianceViolation CancelReason = 2
)

type outbound struct {
	message   contract.ExternalMsg
	seqno     uint32
	logFields map[string]any
}

// composeOutbound constructs outbound message based on the CCTX for further signing.
func (s *Signer) composeOutbound(cctx *cctypes.CrossChainTx) (outbound, error) {
	params := cctx.GetCurrentOutboundParam()

	// #nosec G115 always in range
	seqno := uint32(params.TssNonce)

	recipient, err := ton.ParseAccountID(params.Receiver)
	if err != nil {
		return outbound{}, errors.Wrapf(err, "unable to parse recipient %q", params.Receiver)
	}

	logFields := map[string]any{
		"outbound_recipient": recipient.ToRaw(),
		"outbound_amount":    params.Amount.Uint64(),
		"outbound_nonce":     seqno,
	}

	var cancelReason CancelReason

	// Compliance check (with different address variations)
	if compliance.IsCCTXRestricted(
		cctx,
		recipient.ToRaw(),
		recipient.ToHuman(false, false),
		recipient.ToHuman(true, false),
	) {
		cancelReason = ComplianceViolation

		compliance.PrintComplianceLog(
			s.Logger().Std,
			s.Logger().Compliance,
			true,
			s.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			params.Receiver,
			params.CoinType.String(),
		)
	}

	// Restrict masterchain (`-1:...`) withdrawals
	if recipient.Workchain != 0 {
		cancelReason = InvalidWorkchain
	}

	var message contract.ExternalMsg

	if cancelReason == 0 {
		message = &contract.Withdrawal{
			Recipient: recipient,
			Amount:    params.Amount,
			Seqno:     seqno,
		}
	} else {
		// proceed with a cancellation tx that would only increase the seqno
		// without withdrawing any funds.
		logFields["outbound_cancel"] = true
		logFields["outbound_cancel_reason"] = cancelReason
		message = &contract.IncreaseSeqno{
			Seqno:      seqno,
			ReasonCode: uint32(cancelReason),
		}
	}

	return outbound{
		message,
		seqno,
		logFields,
	}, nil
}
