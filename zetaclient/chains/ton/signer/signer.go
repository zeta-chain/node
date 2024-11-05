package signer

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

// LiteClient represents a TON client
// see https://github.com/ton-blockchain/ton/blob/master/tl/generate/scheme/tonlib_api.tl
//
//go:generate mockery --name LiteClient --structname SignerLiteClient --filename ton_signerliteclient.go --case underscore --output ../../../testutils/mocks
type LiteClient interface {
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) ([]ton.Transaction, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

// Signer represents TON signer.
type Signer struct {
	*base.Signer
	client          LiteClient
	gateway         *toncontracts.Gateway
	signaturesCache *lru.Cache
}

// Signable represents a message that can be signed.
type Signable interface {
	Hash() ([32]byte, error)
	SetSignature(sig [65]byte)
}

// Outcome possible outbound processing outcomes.
type Outcome string

const (
	Invalid Outcome = "invalid"
	Fail    Outcome = "fail"
	Success Outcome = "success"
)

const signaturesHashSize = 1024

var _ interfaces.ChainSigner = (*Signer)(nil)

// New Signer constructor.
func New(baseSigner *base.Signer, client LiteClient, gateway *toncontracts.Gateway) *Signer {
	sigCache, _ := lru.New(signaturesHashSize)

	return &Signer{
		Signer:          baseSigner,
		client:          client,
		gateway:         gateway,
		signaturesCache: sigCache,
	}
}

// TryProcessOutbound tries to process outbound cctx.
// Note that this API signature will be refactored in orchestrator V2
func (s *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *cc.CrossChainTx,
	proc *outboundprocessor.Processor,
	outboundID string,
	_ interfaces.ChainObserver,
	zetacore interfaces.ZetacoreClient,
	zetaBlockHeight uint64,
) {
	proc.StartTryProcess(outboundID)

	defer func() {
		proc.EndTryProcess(outboundID)
	}()

	outcome, err := s.ProcessOutbound(ctx, cctx, zetacore, zetaBlockHeight)

	lf := map[string]any{
		"outbound.id":      outboundID,
		"outbound.nonce":   cctx.GetCurrentOutboundParam().TssNonce,
		"outbound.outcome": string(outcome),
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
	cctx *cc.CrossChainTx,
	zetacore interfaces.ZetacoreClient,
	zetaHeight uint64,
) (Outcome, error) {
	// TODO: note that *InboundParams* are use used on purpose due to legacy reasons.
	// https://github.com/zeta-chain/node/issues/1949
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return Invalid, errors.New("only gas coin outbounds are supported")
	}

	params := cctx.GetCurrentOutboundParam()

	// TODO: add compliance check
	// https://github.com/zeta-chain/node/issues/2916

	receiver, err := ton.ParseAccountID(params.Receiver)
	if err != nil {
		return Invalid, errors.Wrapf(err, "unable to parse recipient %q", params.Receiver)
	}

	withdrawal := &toncontracts.Withdrawal{
		Recipient: receiver,
		Amount:    params.Amount,
		// #nosec G115 always in range
		Seqno: uint32(params.TssNonce),
	}

	lf := map[string]any{
		"outbound.recipient": withdrawal.Recipient.ToRaw(),
		"outbound.amount":    withdrawal.Amount.Uint64(),
		"outbound.nonce":     withdrawal.Seqno,
	}

	s.Logger().Std.Info().Fields(lf).Msg("Signing withdrawal")

	if err = s.SignMessage(ctx, withdrawal, zetaHeight, params.TssNonce); err != nil {
		return Fail, errors.Wrap(err, "unable to sign withdrawal message")
	}

	gwState, err := s.client.GetAccountState(ctx, s.gateway.AccountID())
	if err != nil {
		return Fail, errors.Wrap(err, "unable to get gateway state")
	}

	// Publishes signed message to Gateway
	// Note that max(tx fee) is hardcoded in the contract.
	//
	// Example: If a cctx has amount of 5 TON, the recipient will receive 5 TON,
	// and gateway's balance will be decreased by 5 TON + txFees.
	exitCode, err := s.gateway.SendExternalMessage(ctx, s.client, withdrawal)
	if err != nil || exitCode != 0 {
		return s.handleSendError(exitCode, err, lf)
	}

	// it's okay to run this in the same goroutine
	// because TryProcessOutbound method should be called in a goroutine
	if err = s.trackOutbound(ctx, zetacore, withdrawal, gwState); err != nil {
		return Fail, errors.Wrap(err, "unable to track outbound")
	}

	return Success, nil
}

// SignMessage signs TON external message using TSS
func (s *Signer) SignMessage(ctx context.Context, msg Signable, zetaHeight, nonce uint64) error {
	hash, err := msg.Hash()
	if err != nil {
		return errors.Wrap(err, "unable to hash message")
	}

	// cache hit
	if sig, ok := s.getSignature(hash); ok {
		msg.SetSignature(sig)
		return nil
	}

	chainID := s.Chain().ChainId

	// sig = [65]byte {R, S, V (recovery ID)}
	sig, err := s.TSS().Sign(ctx, hash[:], zetaHeight, nonce, chainID, "")
	if err != nil {
		return errors.Wrap(err, "unable to sign the message")
	}

	msg.SetSignature(sig)
	s.setSignature(hash, sig)

	return nil
}

// because signed msg might fail due to high nonce,
// we need to make sure that signature is cached to avoid redundant TSS calls
func (s *Signer) getSignature(hash [32]byte) ([65]byte, bool) {
	sig, ok := s.signaturesCache.Get(hash)
	if !ok {
		return [65]byte{}, false
	}

	return sig.([65]byte), true
}

// caches signature
func (s *Signer) setSignature(hash [32]byte, sig [65]byte) {
	s.signaturesCache.Add(hash, sig)
}

// Sample (from local ton):
// error code: 0 message: cannot apply external message to current state:
// External message was not accepted Cannot run message on account:
// inbound external message rejected by transaction ...: exitcode=109, steps=108, gas_used=0\
// VM Log (truncated): ...
var exitCodeErrorRegex = regexp.MustCompile(`exitcode=(\d+)`)

// handleSendError tries to figure out the reason of the send error.
func (s *Signer) handleSendError(exitCode uint32, err error, logFields map[string]any) (Outcome, error) {
	if err != nil {
		// Might be possible if 2 concurrent zeta clients
		// are trying to broadcast the same message.
		if strings.Contains(err.Error(), "duplicate message") {
			s.Logger().Std.Warn().Fields(logFields).Msg("Message already sent")
			return Invalid, nil
		}

		var errLiteClient liteclient.LiteServerErrorC
		if errors.As(err, &errLiteClient) {
			logFields["outbound.error.message"] = errLiteClient.Message
			exitCode = errLiteClient.Code
		}

		if code, ok := extractExitCode(err.Error()); ok {
			exitCode = code
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

func extractExitCode(text string) (uint32, bool) {
	match := exitCodeErrorRegex.FindStringSubmatch(text)
	if len(match) < 2 {
		return 0, false
	}

	exitCode, err := strconv.ParseUint(match[1], 10, 32)
	if err != nil {
		return 0, false
	}

	return uint32(exitCode), true
}

// GetGatewayAddress returns gateway address as raw TON address "0:ABC..."
func (s *Signer) GetGatewayAddress() string {
	return s.gateway.AccountID().ToRaw()
}

// SetGatewayAddress sets gateway address. Has a check for noop.
func (s *Signer) SetGatewayAddress(addr string) {
	// noop
	if s.gateway.AccountID().ToRaw() == addr {
		return
	}

	acc, err := ton.ParseAccountID(addr)
	if err != nil {
		s.Logger().Std.Error().Err(err).Str("addr", addr).Msg("unable to parse gateway address")
		return
	}

	s.Lock()
	s.gateway = toncontracts.NewGateway(acc)
	s.Unlock()
}

// not used

func (s *Signer) GetZetaConnectorAddress() (_ ethcommon.Address) { return }
func (s *Signer) GetERC20CustodyAddress() (_ ethcommon.Address)  { return }
func (s *Signer) SetZetaConnectorAddress(_ ethcommon.Address)    {}
func (s *Signer) SetERC20CustodyAddress(_ ethcommon.Address)     {}
