package signer

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	sol "github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// prepareExecuteTx prepares execute outbound
func (signer *Signer) prepareExecuteTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	cancelTx bool,
	logger zerolog.Logger,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()

	// create msg execute
	msg, msgIn, err := signer.createMsgExecute(ctx, cctx, cancelTx)
	if err != nil {
		return signer.prepareIncrementNonceTx(ctx, cctx, height, logger)
	}

	// TSS sign msg execute
	msg, msgIn, err = signMsgWithFallback(ctx, signer, height, params.TssNonce, msg, msgIn)
	if err != nil {
		return nil, err
	}

	return func() (*Outbound, error) {
		inst, err := signer.createExecuteInstruction(*msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating execute instruction")
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

func (signer *Signer) prepareExecuteMsg(
	cctx *types.CrossChainTx,
) (contracts.ExecuteType, *contracts.GenericExecuteMsg, error) {
	var executeType contracts.ExecuteType
	if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.RevertOptions.CallOnRevert {
		executeType = contracts.ExecuteTypeRevert
	} else {
		executeType = contracts.ExecuteTypeCall
	}

	var message []byte
	if executeType == contracts.ExecuteTypeRevert {
		message = cctx.RevertOptions.RevertMessage
	} else {
		messageToDecode, err := hex.DecodeString(cctx.RelayedMessage)
		if err != nil {
			return executeType, nil, errors.Wrapf(err, "decodeString %s error", cctx.RelayedMessage)
		}
		message = messageToDecode
	}

	msg, err := contracts.DecodeExecuteMsg(message)
	if err != nil {
		return executeType, nil, errors.Wrapf(err, "decode ExecuteMsg %s error", cctx.RelayedMessage)
	}

	return executeType, msg, nil
}

func (signer *Signer) prepareExecuteMsgParams(
	ctx context.Context,
	msg *contracts.GenericExecuteMsg,
) ([]*solana.AccountMeta, solana.PublicKeySlice, error) {
	remainingAccounts := []*solana.AccountMeta{}
	if msg.ALTAddress() == nil {
		for _, a := range msg.Legacy.Accounts {
			remainingAccounts = append(remainingAccounts, &solana.AccountMeta{
				PublicKey:  solana.PublicKey(a.PublicKey),
				IsWritable: a.IsWritable,
			})
		}

		return remainingAccounts, nil, nil
	}

	alt, err := addresslookuptable.GetAddressLookupTableStateWithOpts(
		ctx,
		signer.client.(*rpc.Client),
		*msg.ALTAddress(),
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentProcessed},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get alt")
	}

	writableSet := make(map[int]struct{}, len(msg.Alt.WriteableIndexes))
	for _, j := range msg.Alt.WriteableIndexes {
		writableSet[int(j)] = struct{}{}
	}

	for i, a := range alt.Addresses {
		_, isWritable := writableSet[i]
		remainingAccounts = append(remainingAccounts, &solana.AccountMeta{
			PublicKey:  solana.PublicKey(a),
			IsWritable: isWritable,
		})
	}

	return remainingAccounts, alt.Addresses, nil
}

// createMsgExecute creates execute and increment nonce messages
func (signer *Signer) createMsgExecute(
	ctx context.Context,
	cctx *types.CrossChainTx,
	cancelTx bool,
) (*contracts.MsgExecute, *contracts.MsgIncrementNonce, error) {
	params := cctx.GetCurrentOutboundParam()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 lamports through the gateway.
	if cancelTx {
		amount = 0
	}

	// prepare data for msg execute
	executeType, msg, err := signer.prepareExecuteMsg(cctx)
	if err != nil {
		return nil, nil, err
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

	remainingAccounts, altAddresses, err := signer.prepareExecuteMsgParams(ctx, msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot prepare execute msg params")
	}

	msgExecute := contracts.NewMsgExecute(
		chainID,
		nonce,
		amount,
		to,
		sender,
		msg.Data(),
		executeType,
		remainingAccounts,
		msg.ALTAddress(),
		altAddresses,
	)
	msgIncrementNonce := contracts.NewMsgIncrementNonce(chainID, nonce, amount)

	return msgExecute, msgIncrementNonce, nil
}

// createExecuteInstruction wraps the execute 'msg' into a Solana instruction.
func (signer *Signer) createExecuteInstruction(msg contracts.MsgExecute) (*sol.GenericInstruction, error) {
	// create execute instruction with program call data
	discriminator := contracts.DiscriminatorExecute
	var dataBytes []byte
	if msg.ExecuteType() == contracts.ExecuteTypeRevert {
		discriminator = contracts.DiscriminatorExecuteRevert
		serializedInst, err := borsh.Serialize(contracts.ExecuteRevertInstructionParams{
			Discriminator: discriminator,
			Amount:        msg.Amount(),
			Sender:        sol.MustPublicKeyFromBase58(msg.Sender()),
			Data:          msg.Data(),
			Signature:     msg.SigRS(),
			RecoveryID:    msg.SigV(),
			MessageHash:   msg.Hash(),
			Nonce:         msg.Nonce(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "cannot serialize execute_revert instruction")
		}

		dataBytes = serializedInst
	} else {
		serializedInst, err := borsh.Serialize(contracts.ExecuteInstructionParams{
			Discriminator: discriminator,
			Amount:        msg.Amount(),
			Sender:        common.HexToAddress(msg.Sender()),
			Data:          msg.Data(),
			Signature:     msg.SigRS(),
			RecoveryID:    msg.SigV(),
			MessageHash:   msg.Hash(),
			Nonce:         msg.Nonce(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "cannot serialize execute instruction")
		}

		dataBytes = serializedInst
	}

	destinationProgramPda, err := contracts.ComputeConnectedPdaAddress(msg.To())
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode connected pda address")
	}

	predefinedAccounts := []*sol.AccountMeta{
		sol.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
		sol.Meta(signer.pda).WRITE(),
		sol.Meta(msg.To()).WRITE(),
		sol.Meta(destinationProgramPda).WRITE(),
	}
	allAccounts := append(predefinedAccounts, msg.RemainingAccounts()...)

	inst := &sol.GenericInstruction{
		ProgID:        signer.gatewayID,
		DataBytes:     dataBytes,
		AccountValues: allAccounts,
	}

	return inst, nil
}

// validateSender validates and formats the sender address based on execute type
func validateSender(sender string, executeType contracts.ExecuteType) (string, error) {
	if executeType == contracts.ExecuteTypeCall {
		// for regular execute, sender should be an Ethereum address
		senderEth := common.HexToAddress(sender)
		if senderEth == (common.Address{}) {
			return "", fmt.Errorf("invalid execute sender %s", sender)
		}
		return senderEth.Hex(), nil
	}

	// for revert execute, sender should be a Solana address
	senderSol, err := sol.PublicKeyFromBase58(sender)
	if err != nil {
		return "", errors.Wrapf(err, "invalid execute revert sender %s", sender)
	}
	return senderSol.String(), nil
}
