package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// prepareWithdrawTx prepares withdraw outbound
func (signer *Signer) prepareWithdrawTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	cancelTx bool,
	logger zerolog.Logger,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()

	// create msg withdraw
	msg, msgIn, err := signer.createMsgWithdraw(params, cancelTx)
	if err != nil {
		return signer.prepareIncrementNonceTx(ctx, cctx, height, logger)
	}

	// TSS sign msg withdraw
	msg, msgIn, err = signMsgWithFallback(ctx, signer, height, params.TssNonce, msg, msgIn)
	if err != nil {
		return nil, err
	}

	return func() (*Outbound, error) {
		inst, err := signer.createWithdrawInstruction(*msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating withdraw instruction")
		}

		return signer.createOutboundWithFallback(ctx, inst, msgIn, 0, nil, nil)
	}, nil
}

// createMsgWithdraw creates a withdraw and increment nonce messages
func (signer *Signer) createMsgWithdraw(
	params *types.OutboundParams,
	cancelTx bool,
) (*contracts.MsgWithdraw, *contracts.MsgIncrementNonce, error) {
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 lamports through the gateway.
	if cancelTx {
		amount = 0
	}

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	msg := contracts.NewMsgWithdraw(chainID, nonce, amount, to)
	msgIncrementNonce := contracts.NewMsgIncrementNonce(chainID, nonce, amount)

	return msg, msgIncrementNonce, nil
}

// createWithdrawInstruction wraps the withdraw 'msg' into a Solana instruction.
func (signer *Signer) createWithdrawInstruction(msg contracts.MsgWithdraw) (*solana.GenericInstruction, error) {
	// create withdraw instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.WithdrawInstructionParams{
		Discriminator: contracts.DiscriminatorWithdraw,
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize withdraw instruction")
	}

	inst := &solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(msg.To()).WRITE(),
		},
	}

	return inst, nil
}
