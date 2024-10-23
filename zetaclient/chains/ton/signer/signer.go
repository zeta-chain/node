package signer

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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
	client  LiteClient
	gateway *toncontracts.Gateway
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

var _ interfaces.ChainSigner = (*Signer)(nil)

// New Signer constructor.
func New(baseSigner *base.Signer, client LiteClient, gateway *toncontracts.Gateway) *Signer {
	return &Signer{
		Signer:  baseSigner,
		client:  client,
		gateway: gateway,
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
	if err != nil {
		s.Logger().Std.Error().
			Err(err).
			Str("outbound.id", outboundID).
			Uint64("outbound.nonce", cctx.GetCurrentOutboundParam().TssNonce).
			Str("outbound.outcome", string(outcome)).
			Msg("Unable to ProcessOutbound")
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
		Seqno:     uint32(params.TssNonce),
	}

	s.Logger().Std.Info().
		Str("withdrawal.recipient", withdrawal.Recipient.ToRaw()).
		Uint64("withdrawal.amount", withdrawal.Amount.Uint64()).
		Uint32("withdrawal.nonce", withdrawal.Seqno).
		Msg("Signing withdrawal")

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
	if err = s.gateway.SendExternalMessage(ctx, s.client, withdrawal); err != nil {
		return Fail, errors.Wrap(err, "unable to send external message")
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

	chainID := s.Chain().ChainId

	// sig = [65]byte {R, S, V (recovery ID)}
	sig, err := s.TSS().Sign(ctx, hash[:], zetaHeight, nonce, chainID, "")
	if err != nil {
		return errors.Wrap(err, "unable to sign the message")
	}

	msg.SetSignature(sig)

	return nil
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
