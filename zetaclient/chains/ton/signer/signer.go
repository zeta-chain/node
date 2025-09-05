package signer

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/logs"
)

type RPC interface {
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) ([]ton.Transaction, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (rpc.Account, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

// Signer represents TON signer.
type Signer struct {
	*base.Signer
	rpc     RPC
	gateway *toncontracts.Gateway
}

// Outcome possible outbound processing outcomes.
type Outcome string

const (
	Invalid Outcome = "invalid"
	Fail    Outcome = "fail"
	Success Outcome = "success"
)

// New Signer constructor.
func New(baseSigner *base.Signer, rpc RPC, gateway *toncontracts.Gateway) *Signer {
	return &Signer{
		Signer:  baseSigner,
		rpc:     rpc,
		gateway: gateway,
	}
}

// TryProcessOutbound tries to process outbound cctx.
// Note that this API signature will be refactored in orchestrator V2
func (s *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetacore interfaces.ZetacoreClient,
	zetaBlockHeight uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	s.MarkOutbound(outboundID, true)
	defer s.MarkOutbound(outboundID, false)

	outcome, err := s.ProcessOutbound(ctx, cctx, zetacore, zetaBlockHeight)

	lf := map[string]any{
		logs.FieldOutboundID: outboundID,
		logs.FieldNonce:      cctx.GetCurrentOutboundParam().TssNonce,
		"outcome":            string(outcome),
	}

	switch {
	case err != nil:
		s.Logger().Std.Error().Err(err).Fields(lf).Msg("Unable to ProcessOutbound")
	case outcome != Success:
		s.Logger().Std.Warn().Fields(lf).Msg("Unsuccessful outcome for ProcessOutbound")
	default:
		s.Logger().Std.Info().Fields(lf).Msg("Processed outbound")
	}
}

// ProcessOutbound signs and broadcasts an outbound cross-chain transaction.
func (s *Signer) ProcessOutbound(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetacore interfaces.ZetacoreClient,
	zetaHeight uint64,
) (Outcome, error) {
	// TODO: note that *InboundParams* are use used on purpose due to legacy reasons.
	// https://github.com/zeta-chain/node/issues/1949
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return Invalid, errors.New("only gas coin outbounds are supported")
	}

	nonce := cctx.GetCurrentOutboundParam().TssNonce

	outbound, err := s.composeOutbound(cctx)
	if err != nil {
		return Invalid, errors.Wrap(err, "unable to compose message")
	}

	s.Logger().Std.Info().Fields(outbound.logFields).Msg("Signing outbound")

	if err = s.SignMessage(ctx, outbound.message, zetaHeight, nonce); err != nil {
		return Fail, errors.Wrap(err, "unable to sign withdrawal message")
	}

	gwState, err := s.rpc.GetAccountState(ctx, s.gateway.AccountID())
	if err != nil {
		return Fail, errors.Wrap(err, "unable to get gateway state")
	}

	// Publishes signed message to Gateway
	// Note that max(tx fee) is hardcoded in the contract.
	//
	// Example: If a cctx has amount of 5 TON, the recipient will receive 5 TON,
	// and gateway's balance will be decreased by 5 TON + txFees.
	exitCode, err := s.gateway.SendExternalMessage(ctx, s.rpc, outbound.message)
	if err != nil || exitCode != 0 {
		return s.handleSendError(exitCode, err, outbound.logFields)
	}

	// it's okay to run this in the same goroutine
	// because TryProcessOutbound method should be called in a goroutine
	if err = s.trackOutbound(ctx, zetacore, outbound, gwState); err != nil {
		return Fail, errors.Wrap(err, "unable to track outbound")
	}

	return Success, nil
}

// SignMessage signs TON external message using TSS
// Note that TSS has in-mem cache for existing signatures to abort duplicate signing requests.
func (s *Signer) SignMessage(ctx context.Context, msg toncontracts.ExternalMsg, zetaHeight, nonce uint64) error {
	hash, err := msg.Hash()
	if err != nil {
		return errors.Wrap(err, "unable to hash message")
	}

	chainID := s.Chain().ChainId

	// sig = [65]byte {R, S, V (recovery ID)}
	sig, err := s.TSS().Sign(ctx, hash[:], zetaHeight, nonce, chainID)
	if err != nil {
		return errors.Wrap(err, "unable to sign the message")
	}

	msg.SetSignature(sig)

	return nil
}

// handleSendError tries to figure out the reason of the send error.
func (s *Signer) handleSendError(exitCode uint32, err error, logFields map[string]any) (Outcome, error) {
	if err != nil {
		// Might be possible if 2 concurrent zeta clients
		// are trying to broadcast the same message.
		if strings.Contains(err.Error(), "duplicate") {
			s.Logger().Std.Warn().Fields(logFields).Msg("Message already sent")
			return Invalid, nil
		}
	}

	switch {
	case exitCode == uint32(toncontracts.ExitCodeInvalidSeqno):
		// Might be possible if zeta clients send several seq. numbers concurrently.
		// In the current implementation, Gateway supports only 1 nonce per block.
		logFields["outbound.error.exit_code"] = exitCode
		s.Logger().Std.Warn().Fields(logFields).Msg("Invalid nonce, retry later")
		return Invalid, nil
	case err != nil:
		return Fail, errors.Wrap(err, "unable to send external message")
	default:
		return Fail, errors.Errorf("unable to send external message: exit code %d", exitCode)
	}
}

// GetGatewayAddress returns gateway address as raw TON address "0:ABC..."
func (s *Signer) GetGatewayAddress() string {
	return s.gateway.AccountID().ToRaw()
}

// SetGatewayAddress sets gateway address. Has a check for noop.
func (s *Signer) SetGatewayAddress(addr string) {
	// noop
	if addr == "" || s.gateway.AccountID().ToRaw() == addr {
		return
	}

	acc, err := ton.ParseAccountID(addr)
	if err != nil {
		s.Logger().Std.Error().Err(err).Str("addr", addr).Msg("unable to parse gateway address")
		return
	}

	s.Logger().Std.Info().
		Str("signer.old_gateway_address", s.gateway.AccountID().ToRaw()).
		Str("signer.new_gateway_address", acc.ToRaw()).
		Msg("Updated gateway address")

	s.Lock()
	s.gateway = toncontracts.NewGateway(acc)
	s.Unlock()
}
