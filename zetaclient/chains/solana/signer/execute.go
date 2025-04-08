package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// createAndSignMsgExecute creates and batch signs execute and increment nonce messages
// for gateway execute instruction with TSS.
func (signer *Signer) createAndSignMsgExecute(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	sender string,
	data []byte,
	remainingAccounts []*solana.AccountMeta,
	revert bool,
	cancelTx bool,
) (*contracts.MsgExecute, *contracts.MsgIncrementNonce, error) {
	chain := signer.Chain()
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

	// prepare execute msg and compute hash
	msg := contracts.NewMsgExecute(chainID, nonce, amount, to, sender, data, revert, remainingAccounts)
	msgHash := msg.Hash()

	// prepare increment_nonce msg and compute hash, it will be used as fallback tx in case execute fails
	msgIncrementNonce := contracts.NewMsgIncrementNonce(chainID, nonce, amount)
	msgHashIncrementNonce := msgIncrementNonce.Hash()

	// sign the message with TSS to get an ECDSA signature.
	// the produced signature is in the [R || S || V] format where V is 0 or 1.
	signature, err := signer.TSS().
		SignBatch(ctx, [][]byte{msgHash[:], msgHashIncrementNonce[:]}, height, nonce, chain.ChainId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "key-sign failed")
	}

	// attach the signature and return
	return msg.SetSignature(signature[0]), msgIncrementNonce.SetSignature(signature[1]), nil
}

// createExecuteInstruction wraps the execute 'msg' into a Solana instruction.
func (signer *Signer) createExecuteInstruction(msg contracts.MsgExecute) (*solana.GenericInstruction, error) {
	// create execute instruction with program call data
	discriminator := contracts.DiscriminatorExecute
	dataBytes := []byte{}
	if msg.Revert() {
		discriminator = contracts.DiscriminatorExecuteRevert
		serializedInst, err := borsh.Serialize(contracts.ExecuteRevertInstructionParams{
			Discriminator: discriminator,
			Amount:        msg.Amount(),
			Sender:        solana.MustPublicKeyFromBase58(msg.Sender()),
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

	predefinedAccounts := []*solana.AccountMeta{
		solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
		solana.Meta(signer.pda).WRITE(),
		solana.Meta(msg.To()).WRITE(),
		solana.Meta(destinationProgramPda),
	}
	allAccounts := append(predefinedAccounts, msg.RemainingAccounts()...)

	inst := &solana.GenericInstruction{
		ProgID:        signer.gatewayID,
		DataBytes:     dataBytes,
		AccountValues: allAccounts,
	}

	return inst, nil
}
