package signer

import (
	"context"
	"crypto/sha256"
	"fmt"

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
	outboundID := base.OutboundIDFromCCTX(cctx)
	s.MarkOutbound(outboundID, true)
	defer func() { s.MarkOutbound(outboundID, false) }()

	var (
		err         error
		tx          models.TxnMetaData
		txFallback  *models.TxnMetaData
		sig         string
		sigFallback *string
		nonce       = cctx.GetCurrentOutboundParam().TssNonce
	)

	// build outbound txs
	if cctx.IsWithdrawAndCall() {
		tx, txFallback, err = s.buildExecute(ctx, cctx)
		if err != nil {
			return errors.Wrap(err, "unable to build execute tx")
		}
	} else {
		tx, err = s.buildWithdraw(ctx, cctx)
		if err != nil {
			return errors.Wrap(err, "unable to build withdraw tx")
		}
	}

	// sign outbound txs
	sig, sigFallback, err = s.signTxWithFallback(ctx, tx, txFallback, zetaHeight, nonce)
	if err != nil {
		return errors.Wrap(err, "unable to sign tx with fallback")
	}

	// broadcast outbound txs
	txDigest, err := s.broadcastWithFallback(ctx, tx, txFallback, sig, sigFallback)
	if err != nil {
		// todo we might need additional error handling
		// for the case when the tx is already broadcasted by another zetaclient
		// (e.g. suppress error)
		return errors.Wrap(err, "unable to broadcast tx")
	}

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

// signTxWithFallback signs original tx with an optional fallback tx.
// Pointers types are used to be flexible and indicate optional fallback tx.
func (s *Signer) signTxWithFallback(
	ctx context.Context,
	tx models.TxnMetaData,
	txFallback *models.TxnMetaData,
	zetaHeight, nonce uint64,
) (sig string, sigFallback *string, err error) {
	digests := [][]byte{}

	// collect digests
	digest, err := sui.Digest(tx)
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to get digest")
	}
	digests = append(digests, wrapDigest(digest))

	if txFallback != nil {
		digest, err := sui.Digest(*txFallback)
		if err != nil {
			return "", nil, errors.Wrap(err, "unable to get fallback digest")
		}
		digests = append(digests, wrapDigest(digest))
	}

	// sign digests with TSS
	sig65Bs, err := s.TSS().SignBatch(ctx, digests, zetaHeight, nonce, s.Chain().ChainId)
	if err != nil {
		return "", nil, errors.Wrapf(err, "unable to sign %d tx(s) with TSS", len(digests))
	}

	// should never mismatch
	if len(sig65Bs) != len(digests) {
		return "", nil, fmt.Errorf("expected %d signatures, got %d", len(digests), len(sig65Bs))
	}

	// serialize signatures
	sig, err = sui.SerializeSignatureECDSA(sig65Bs[0], s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to serialize tx signature")
	}

	if txFallback != nil {
		sigBase64, err := sui.SerializeSignatureECDSA(sig65Bs[1], s.TSS().PubKey().AsECDSA())
		if err != nil {
			return "", nil, errors.Wrap(err, "unable to serialize tx fallback signature")
		}
		sigFallback = &sigBase64
	}

	return sig, sigFallback, nil
}

// wrapDigest wraps the digest with sha256.
// another hashing is required for ECDSA.
// see: https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
func wrapDigest(digest [32]byte) []byte {
	digestWrapped := sha256.Sum256(digest[:])
	return digestWrapped[:]
}
