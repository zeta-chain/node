package signer

import (
	"context"

	"cosmossdk.io/errors"
	sol "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// prepareWithdrawSPLTx prepares withdraw spl outbound
func (signer *Signer) prepareWithdrawSPLTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	cancelTx bool,
	logger zerolog.Logger,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()

	// create msg withdraw spl
	msg, msgIn, err := signer.createMsgWithdrawSPL(
		ctx,
		cctx,
		cancelTx,
	)
	if err != nil {
		return signer.prepareIncrementNonceTx(ctx, cctx, height, logger)
	}

	// TSS sign msg withdraw spl
	msg, msgIn, err = signMsgWithFallback(ctx, signer, height, params.TssNonce, msg, msgIn)
	if err != nil {
		return nil, err
	}

	return func() (*Outbound, error) {
		inst, err := signer.createWithdrawSPLInstruction(*msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating withdraw SPL instruction")
		}

		return signer.createOutboundWithFallback(ctx, inst, msgIn, 0)
	}, nil
}

// createMsgWithdrawSPL creates withdraw spl and increment nonce messages
func (signer *Signer) createMsgWithdrawSPL(
	ctx context.Context,
	cctx *types.CrossChainTx,
	cancelTx bool,
) (*contracts.MsgWithdrawSPL, *contracts.MsgIncrementNonce, error) {
	params := cctx.GetCurrentOutboundParam()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 spl through the gateway.
	if cancelTx {
		amount = 0
	}

	// get mint details to get decimals
	mint, err := signer.decodeMintAccountDetails(ctx, cctx.InboundParams.Asset)
	if err != nil {
		return nil, nil, errors.Wrap(err, "decodeMintAccountDetails error")
	}

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// parse mint account
	mintAccount, err := sol.PublicKeyFromBase58(cctx.InboundParams.Asset)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot parse asset public key %s", cctx.InboundParams.Asset)
	}

	// get recipient ata
	recipientAta, _, err := sol.FindAssociatedTokenAddress(to, mintAccount)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", to, mintAccount)
	}

	msg := contracts.NewMsgWithdrawSPL(chainID, nonce, amount, mint.Decimals, mintAccount, to, recipientAta)
	msgIncrementNonce := contracts.NewMsgIncrementNonce(chainID, nonce, amount)

	return msg, msgIncrementNonce, nil
}

// createWithdrawSPLInstruction wraps the withdraw spl 'msg' into a Solana instruction.
func (signer *Signer) createWithdrawSPLInstruction(msg contracts.MsgWithdrawSPL) (*sol.GenericInstruction, error) {
	// create withdraw spl instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.WithdrawSPLInstructionParams{
		Discriminator: contracts.DiscriminatorWithdrawSPL,
		Decimals:      msg.Decimals(),
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize withdraw instruction")
	}

	pdaAta, _, err := sol.FindAssociatedTokenAddress(signer.pda, msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", signer.pda, msg.MintAccount())
	}

	recipientAta, _, err := sol.FindAssociatedTokenAddress(msg.To(), msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", msg.To(), msg.MintAccount())
	}

	inst := &sol.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*sol.AccountMeta{
			sol.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			sol.Meta(signer.pda).WRITE(),
			sol.Meta(pdaAta).WRITE(),
			sol.Meta(msg.MintAccount()),
			sol.Meta(msg.To()),
			sol.Meta(recipientAta).WRITE(),
			sol.Meta(sol.TokenProgramID),
			sol.Meta(sol.SPLAssociatedTokenAccountProgramID),
			sol.Meta(sol.SystemProgramID),
		},
	}

	return inst, nil
}

func (signer *Signer) decodeMintAccountDetails(ctx context.Context, asset string) (token.Mint, error) {
	mintPk, err := sol.PublicKeyFromBase58(asset)
	if err != nil {
		return token.Mint{}, err
	}

	info, err := signer.solanaClient.GetAccountInfo(ctx, mintPk)
	if err != nil {
		return token.Mint{}, err
	}

	return contracts.DeserializeMintAccountInfo(info)
}
