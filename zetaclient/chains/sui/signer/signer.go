package signer

import (
	"context"
	"crypto/sha256"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// Signer Sui outbound transaction signer.
type Signer struct {
	*base.Signer
	client      RPC
	gateway     *sui.Gateway
	withdrawCap *withdrawCap

	zetacore interfaces.ZetacoreClient
}

// RPC represents Sui rpc.
type RPC interface {
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)

	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}

// New Signer constructor.
func New(
	baseSigner *base.Signer,
	client RPC,
	gateway *sui.Gateway,
	zetacore interfaces.ZetacoreClient,
) *Signer {
	return &Signer{
		Signer:      baseSigner,
		client:      client,
		gateway:     gateway,
		zetacore:    zetacore,
		withdrawCap: &withdrawCap{},
	}
}

// ProcessCCTX schedules outbound cross-chain transaction.
// Build --> Sign --> Broadcast --(async)--> Wait for execution --> PostOutboundTracker
func (s *Signer) ProcessCCTX(ctx context.Context, cctx *cctypes.CrossChainTx, zetaHeight uint64) error {
	var (
		outboundID = base.OutboundIDFromCCTX(cctx)
		nonce      = cctx.GetCurrentOutboundParam().TssNonce
	)

	s.MarkOutbound(outboundID, true)
	defer func() { s.MarkOutbound(outboundID, false) }()

	tx, err := s.buildWithdrawal(ctx, cctx)
	if err != nil {
		return errors.Wrap(err, "unable to build withdrawal tx")
	}

	sig, err := s.signTx(ctx, tx, zetaHeight, nonce)
	if err != nil {
		return errors.Wrap(err, "unable to sign tx")
	}

	txDigest, err := s.broadcast(ctx, tx, sig)
	if err != nil {
		// todo we might need additional error handling
		// for the case when the tx is already broadcasted by another zetaclient
		// (e.g. suppress error)
		return errors.Wrap(err, "unable to broadcast tx")
	}

	s.Logger().Std.Info().Str(logs.FieldTx, txDigest).Msg("Broadcasted transaction")

	logger := s.Logger().Std.With().
		Str(logs.FieldMethod, "reportToOutboundTracker").
		Int64(logs.FieldChain, s.Chain().ChainId).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txDigest).
		Logger()

	ctx = logger.WithContext(ctx)

	bg.Work(ctx,
		func(ctx context.Context) error { return s.reportOutboundTracker(ctx, nonce, txDigest) },
		bg.WithLogger(logger),
		bg.WithName("report_outbound_tracker"),
	)

	return nil
}

func (s *Signer) signTx(ctx context.Context, tx models.TxnMetaData, zetaHeight, nonce uint64) ([65]byte, error) {
	digest, err := sui.Digest(tx)
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "unable to get digest")
	}

	// Another hashing is required for ECDSA.
	// https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
	digestWrapped := sha256.Sum256(digest[:])

	// Send TSS signature request.
	return s.TSS().Sign(
		ctx,
		digestWrapped[:],
		zetaHeight,
		nonce,
		s.Chain().ChainId,
	)
}
