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
	SuiXGetLatestSuiSystemState(ctx context.Context) (models.SuiSystemStateSummary, error)
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)

	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	SuiDevInspectTransactionBlock(
		ctx context.Context,
		req models.SuiDevInspectTransactionBlockRequest,
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
		err             error
		sig             string
		tx              models.TxnMetaData
		cancelTxBuilder txBuilder
		nonce           = cctx.GetCurrentOutboundParam().TssNonce
	)

	// build outbound txs
	if cctx.IsWithdrawAndCall() {
		tx, cancelTxBuilder, err = s.buildExecuteWithCancelTxBuilder(ctx, cctx, zetaHeight)
		if err != nil {
			return errors.Wrap(err, "unable to build execute tx")
		}
	} else {
		tx, err = s.buildWithdraw(ctx, cctx)
		if err != nil {
			return errors.Wrap(err, "unable to build withdraw tx")
		}
	}

	// sign tx
	sig, err = s.signTx(ctx, tx, zetaHeight, nonce)
	if err != nil {
		return errors.Wrap(err, "unable to sign tx")
	}

	// prepare logger
	logger := s.Logger().Std.With().
		Int64(logs.FieldChain, s.Chain().ChainId).
		Uint64(logs.FieldNonce, nonce).
		Logger()
	ctx = logger.WithContext(ctx)

	// broadcast tx with cancel tx
	txDigest, err := s.broadcastWithCancelTx(ctx, sig, tx, cancelTxBuilder)
	if err != nil {
		// todo we might need additional error handling
		// for the case when the tx is already broadcasted by another zetaclient
		// (e.g. suppress error)
		return errors.Wrap(err, "unable to broadcast tx")
	}

	logger = logger.With().Str(logs.FieldTx, txDigest).Logger()
	ctx = logger.WithContext(ctx)

	bg.Work(ctx,
		func(ctx context.Context) error { return s.reportOutboundTracker(ctx, nonce, txDigest) },
		bg.WithLogger(logger),
		bg.WithName("report_outbound_tracker"),
	)

	return nil
}

// signTx signs a tx with TSS and returns a base64 encoded signature.
func (s *Signer) signTx(ctx context.Context, tx models.TxnMetaData, zetaHeight, nonce uint64) (string, error) {
	digest, err := sui.Digest(tx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get digest")
	}

	// send TSS signature request.
	sig65B, err := s.TSS().Sign(ctx, wrapDigest(digest), zetaHeight, nonce, s.Chain().ChainId)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	// serialize signature
	sig, err := sui.SerializeSignatureECDSA(sig65B, s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", errors.Wrap(err, "unable to serialize tx signature")
	}

	return sig, nil
}

// SignTxWithCancel signs original tx with an optional cancel tx.
// Pointers type is used to be flexible and indicate optional cancel tx.
//
// Note: this function is not used due to tx simulation issue in Sui SDK,
// but we can sign both tx and cancel tx in one go once Sui SDK is updated.
func (s *Signer) SignTxWithCancel(
	ctx context.Context,
	tx models.TxnMetaData,
	txCancel *models.TxnMetaData,
	zetaHeight, nonce uint64,
) (sig string, sigCancel *string, err error) {
	digests := [][]byte{}

	// collect digests
	digest, err := sui.Digest(tx)
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to get digest")
	}
	digests = append(digests, wrapDigest(digest))

	if txCancel != nil {
		digest, err = sui.Digest(*txCancel)
		if err != nil {
			return "", nil, errors.Wrap(err, "unable to get cancel tx digest")
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

	if txCancel != nil {
		sigBase64, err := sui.SerializeSignatureECDSA(sig65Bs[1], s.TSS().PubKey().AsECDSA())
		if err != nil {
			return "", nil, errors.Wrap(err, "unable to serialize tx cancel signature")
		}
		sigCancel = &sigBase64
	}

	return sig, sigCancel, nil
}

// wrapDigest wraps the digest with sha256.
func wrapDigest(digest [32]byte) []byte {
	// another hashing is required for ECDSA.
	// see: https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
	digestWrapped := sha256.Sum256(digest[:])
	return digestWrapped[:]
}
