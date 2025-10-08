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
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
)

type TONClient interface {
	GetTransactionsSince(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)

	GetAccountState(context.Context, ton.AccountID) (rpc.Account, error)

	SendMessage(context.Context, []byte) (uint32, error)
}

// Signer represents TON signer.
type Signer struct {
	*base.Signer
	tonClient TONClient
	gateway   *toncontracts.Gateway
}

// Outcome possible outbound processing outcomes.
type Outcome string

const (
	Invalid Outcome = "invalid"
	Fail    Outcome = "fail"
	Success Outcome = "success"
)

// New Signer constructor.
func New(baseSigner *base.Signer, tonClient TONClient, gateway *toncontracts.Gateway) *Signer {
	return &Signer{
		Signer:    baseSigner,
		tonClient: tonClient,
		gateway:   gateway,
	}
}

// TryProcessOutbound tries to process an outbound CCTX.
func (s *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetaRepo *zrepo.ZetaRepo,
	zetaBlockHeight uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	s.MarkOutbound(outboundID, true)
	defer s.MarkOutbound(outboundID, false)

	outcome, err := s.processOutbound(ctx, cctx, zetaRepo, zetaBlockHeight)

	logger := s.Logger().Std.With().
		Str(logs.FieldOutboundID, outboundID).
		Uint64(logs.FieldNonce, cctx.GetCurrentOutboundParam().TssNonce).
		Str("outcome", string(outcome)).
		Logger()

	if err != nil {
		logger.Error().Err(err).Msg("error calling ProcessOutbound")
		return
	}

	if outcome != Success {
		logger.Warn().Msg("unsuccessful outcome for ProcessOutbound")
		return
	}

	logger.Info().Msg("processed outbound")
}

// processOutbound signs and broadcasts an outbound cross-chain transaction.
func (s *Signer) processOutbound(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetaRepo *zrepo.ZetaRepo,
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

	logger := s.Logger().Std.With().Fields(outbound.logFields).Logger()

	if s.ClientMode.IsDryMode() {
		logger.Info().
			Stringer(logs.FieldMode, mode.DryMode).
			Msg("skipping TON signing, sending, and tracking")
		return Success, nil
	}

	logger.Info().Msg("signing outbound")

	err = s.SignMessage(ctx, outbound.message, zetaHeight, nonce)
	if err != nil {
		return Fail, errors.Wrap(err, "unable to sign withdrawal message")
	}

	gwState, err := s.tonClient.GetAccountState(ctx, s.gateway.AccountID())
	if err != nil {
		return Fail, errors.Wrap(err, "unable to get gateway state")
	}

	// Publishes signed message to Gateway
	// Note that max(tx fee) is hardcoded in the contract.
	//
	// Example: If a cctx has amount of 5 TON, the recipient will receive 5 TON,
	// and gateway's balance will be decreased by 5 TON + txFees.
	exitCode, err := s.gateway.SendExternalMessage(ctx, s.tonClient, outbound.message)
	if err != nil || exitCode != 0 {
		return s.handleSendError(exitCode, err, outbound.logFields)
	}

	// it's okay to run this in the same goroutine
	// because TryProcessOutbound method should be called in a goroutine
	err = s.trackOutbound(ctx, zetaRepo, outbound, gwState)
	if err != nil {
		return Fail, errors.Wrap(err, "unable to track outbound")
	}

	return Success, nil
}

// SignMessage signs TON external message using TSS
// Note that TSS has in-mem cache for existing signatures to abort duplicate signing requests.
func (s *Signer) SignMessage(ctx context.Context,
	msg toncontracts.ExternalMsg,
	zetaHeight,
	nonce uint64,
) error {
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
			s.Logger().Std.Warn().Fields(logFields).Msg("message already sent")
			return Invalid, nil
		}
	}

	switch {
	case exitCode == uint32(toncontracts.ExitCodeInvalidSeqno):
		// Might be possible if zeta clients send several seq. numbers concurrently.
		// In the current implementation, Gateway supports only 1 nonce per block.
		logFields["outbound_error_exit_code"] = exitCode
		s.Logger().Std.Warn().Fields(logFields).Msg("invalid nonce, retry later")
		return Invalid, nil
	case err != nil:
		return Fail, errors.Wrap(err, "unable to send external message")
	default:
		return Fail, errors.Errorf("unable to send external message: exit code %d", exitCode)
	}
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
		Str("signer_old_gateway_address", s.gateway.AccountID().ToRaw()).
		Str("signer_new_gateway_address", acc.ToRaw()).
		Msg("Updated gateway address")

	s.Lock()
	s.gateway = toncontracts.NewGateway(acc)
	s.Unlock()
}
