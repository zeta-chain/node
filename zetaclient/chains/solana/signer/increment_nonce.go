package signer

import (
	"context"

	"cosmossdk.io/errors"
	sol "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/compliance"
)

// prepareIncrementNonceTx prepares increment nonce outbound
func (signer *Signer) prepareIncrementNonceTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	logger zerolog.Logger,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()
	// compliance check
	cancelTx := compliance.IsCCTXRestricted(cctx)
	if cancelTx {
		coinType := coin.CoinType_Gas
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			params.Receiver,
			&coinType,
		)
	}

	// create and sign gateway increment nonce message by TSS
	msg, err := signer.createAndSignMsgIncrementNonce(ctx, params, height, cancelTx)
	if err != nil {
		return nil, err
	}

	return func() (*Outbound, error) {
		// sign the increment_nonce transaction by relayer key
		inst, err := signer.createIncrementNonceInstruction(*msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating increment nonce instruction")
		}

		tx, err := signer.signTx(ctx, inst, 0)
		if err != nil {
			return nil, errors.Wrap(err, "error signing increment nonce instruction")
		}
		return &Outbound{Tx: tx}, nil
	}, nil
}

// createAndSignMsgIncrementNonce creates and signs a increment_nonce message for gateway increment_nonce instruction with TSS.
func (signer *Signer) createAndSignMsgIncrementNonce(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	cancelTx bool,
) (*contracts.MsgIncrementNonce, error) {
	chain := signer.Chain()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 lamports through the gateway.
	if cancelTx {
		amount = 0
	}

	// prepare increment_nonce msg and compute hash
	msg := contracts.NewMsgIncrementNonce(chainID, nonce, amount)
	msgHash := msg.Hash()

	// sign the message with TSS to get an ECDSA signature.
	// the produced signature is in the [R || S || V] format where V is 0 or 1.
	signature, err := signer.TSS().Sign(ctx, msgHash[:], height, nonce, chain.ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "key-sign failed")
	}

	// attach the signature and return
	return msg.SetSignature(signature), nil
}

// createIncrementNonceInstruction wraps the increment_nonce 'msg' into a Solana instruction.
func (signer *Signer) createIncrementNonceInstruction(
	msg contracts.MsgIncrementNonce,
) (*sol.GenericInstruction, error) {
	// create increment_nonce instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.IncrementNonceInstructionParams{
		Discriminator: contracts.DiscriminatorIncrementNonce,
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
		FailureReason: msg.FailureReason(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize increment_nonce instruction")
	}

	inst := &sol.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*sol.AccountMeta{
			sol.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			sol.Meta(signer.pda).WRITE(),
		},
	}

	return inst, nil
}
