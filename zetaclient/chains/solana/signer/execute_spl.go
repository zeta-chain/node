package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	sol "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// prepareExecuteSPLTx prepares execute spl outbound
func (signer *Signer) prepareExecuteSPLTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	cancelTx bool,
	logger zerolog.Logger,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()

	// create msg execute spl
	msg, msgIn, err := signer.createMsgExecuteSPL(ctx, cctx, cancelTx)
	if err != nil {
		return signer.prepareIncrementNonceTx(ctx, cctx, height, logger)
	}

	// TSS sign msg execute spl
	msg, msgIn, err = signMsgWithFallback(ctx, signer, height, params.TssNonce, msg, msgIn)
	if err != nil {
		return nil, err
	}

	return func() (*Outbound, error) {
		inst, err := signer.createExecuteSPLInstruction(*msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating execute SPL instruction")
		}

		return signer.createOutboundWithFallback(
			ctx,
			inst,
			msgIn,
			params.CallOptions.GasLimit,
			msg.ALT(),
			msg.ALTStateAddresses(),
		)
	}, nil
}

// createMsgExecuteSPL creates execute spl and increment nonce messages
func (signer *Signer) createMsgExecuteSPL(
	ctx context.Context,
	cctx *types.CrossChainTx,
	cancelTx bool,
) (*contracts.MsgExecuteSPL, *contracts.MsgIncrementNonce, error) {
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
		return nil, nil, errors.Wrap(err, "decoding mint account details")
	}

	executeType, msg, err := signer.prepareExecuteMsg(cctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "prepare ExecuteMsg error")
	}

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// check sender based on execute type
	sender, err := validateSender(cctx.InboundParams.Sender, executeType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot validate sender")
	}

	// parse mint account
	mintAccount, err := sol.PublicKeyFromBase58(cctx.InboundParams.Asset)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot parse asset public key %s", cctx.InboundParams.Asset)
	}

	// get recipient ata
	destinationProgramPda, err := contracts.ComputeConnectedPdaAddress(to)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot decode connected spl pda address")
	}

	destinationProgramPdaAta, _, err := sol.FindAssociatedTokenAddress(destinationProgramPda, mintAccount)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"cannot find ATA for %s and mint account %s",
			destinationProgramPda,
			mintAccount,
		)
	}

	remainingAccounts, altAddresses, err := signer.prepareExecuteMsgParams(ctx, msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot prepare execute msg params")
	}

	// prepare execute spl and increment nonce messages
	msgExecuteSPL := contracts.NewMsgExecuteSPL(
		chainID,
		nonce,
		amount,
		mint.Decimals,
		mintAccount,
		to,
		destinationProgramPdaAta,
		sender,
		msg.Data(),
		executeType,
		remainingAccounts,
		msg.ALTAddress(),
		altAddresses,
	)

	msgIncrementNonce := contracts.NewMsgIncrementNonce(chainID, nonce, amount)

	return msgExecuteSPL, msgIncrementNonce, nil
}

// createExecuteSPLInstruction wraps the execute spl 'msg' into a Solana instruction.
func (signer *Signer) createExecuteSPLInstruction(msg contracts.MsgExecuteSPL) (*sol.GenericInstruction, error) {
	// create execute spl instruction with program call data
	var dataBytes []byte
	if msg.ExecuteType() == contracts.ExecuteTypeRevert {
		serializedInst, err := borsh.Serialize(contracts.ExecuteSPLRevertInstructionParams{
			Discriminator: contracts.DiscriminatorExecuteSPLRevert,
			Decimals:      msg.Decimals(),
			Amount:        msg.Amount(),
			Sender:        sol.MustPublicKeyFromBase58(msg.Sender()),
			Data:          msg.Data(),
			Signature:     msg.SigRS(),
			RecoveryID:    msg.SigV(),
			MessageHash:   msg.Hash(),
			Nonce:         msg.Nonce(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "cannot serialize execute spl instruction")
		}

		dataBytes = serializedInst
	} else {
		serializedInst, err := borsh.Serialize(contracts.ExecuteSPLInstructionParams{
			Discriminator: contracts.DiscriminatorExecuteSPL,
			Decimals:      msg.Decimals(),
			Amount:        msg.Amount(),
			Sender:        common.HexToAddress(msg.Sender()),
			Data:          msg.Data(),
			Signature:     msg.SigRS(),
			RecoveryID:    msg.SigV(),
			MessageHash:   msg.Hash(),
			Nonce:         msg.Nonce(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "cannot serialize execute spl instruction")
		}

		dataBytes = serializedInst
	}

	pdaAta, _, err := sol.FindAssociatedTokenAddress(signer.pda, msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", signer.pda, msg.MintAccount())
	}

	destinationProgramPda, err := contracts.ComputeConnectedPdaAddress(msg.To())
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode connected spl pda address")
	}

	predefinedAccounts := []*sol.AccountMeta{
		sol.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
		sol.Meta(signer.pda).WRITE(),
		sol.Meta(pdaAta).WRITE(),
		sol.Meta(msg.MintAccount()),
		sol.Meta(msg.To()),
		sol.Meta(destinationProgramPda).WRITE(),
		sol.Meta(msg.RecipientAta()).WRITE(),
		sol.Meta(sol.TokenProgramID),
		sol.Meta(sol.SPLAssociatedTokenAccountProgramID),
		sol.Meta(sol.SystemProgramID),
	}
	allAccounts := append(predefinedAccounts, msg.RemainingAccounts()...)

	inst := &sol.GenericInstruction{
		ProgID:        signer.gatewayID,
		DataBytes:     dataBytes,
		AccountValues: allAccounts,
	}

	return inst, nil
}
