package signer

import (
	"context"
	"crypto/sha256"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
)

// Signer Sui outbound transaction signer.
type Signer struct {
	*base.Signer
	client      RPC
	gateway     *sui.Gateway
	withdrawCap *withdrawCap
}

// RPC represents Sui rpc.
type RPC interface {
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)

	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}

// New Signer constructor.
func New(baseSigner *base.Signer, client RPC, gateway *sui.Gateway) *Signer {
	return &Signer{
		Signer:      baseSigner,
		client:      client,
		gateway:     gateway,
		withdrawCap: &withdrawCap{},
	}
}

// ProcessCCTX schedules outbound cross-chain transaction.
func (s *Signer) ProcessCCTX(ctx context.Context, cctx *cctypes.CrossChainTx, zetaHeight uint64) error {
	// todo ... vote if confirmed, etc ...

	nonce := cctx.GetCurrentOutboundParam().TssNonce

	tx, err := s.buildWithdrawal(ctx, cctx)
	if err != nil {
		return errors.Wrap(err, "unable to build withdrawal tx")
	}

	sig, err := s.signTx(ctx, tx, zetaHeight, nonce)
	if err != nil {
		return errors.Wrap(err, "unable to sign tx")
	}

	if err := s.broadcast(ctx, tx, sig); err != nil {
		return errors.Wrap(err, "unable to broadcast tx")
	}

	// todo ...

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
