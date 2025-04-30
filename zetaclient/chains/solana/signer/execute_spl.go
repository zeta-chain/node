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

// createAndSignMsgExecuteSPL creates and batch signs execute spl and increment nonce messages
// for gateway execute_spl_token instruction with TSS.
func (signer *Signer) createAndSignMsgExecuteSPL(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	asset string,
	decimals uint8,
	sender string,
	data []byte,
	remainingAccounts []*solana.AccountMeta,
	executeType contracts.ExecuteType,
	cancelTx bool,
) (*contracts.MsgExecuteSPL, *contracts.MsgIncrementNonce, error) {
	chain := signer.Chain()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 spl through the gateway.
	if cancelTx {
		amount = 0
	}

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// parse mint account
	mintAccount, err := solana.PublicKeyFromBase58(asset)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot parse asset public key %s", asset)
	}

	// get recipient ata
	destinationProgramPda, err := contracts.ComputeConnectedPdaAddress(to)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot decode connected spl pda address")
	}

	destinationProgramPdaAta, _, err := solana.FindAssociatedTokenAddress(destinationProgramPda, mintAccount)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"cannot find ATA for %s and mint account %s",
			destinationProgramPda,
			mintAccount,
		)
	}

	// prepare execute spl msg and compute hash
	msg := contracts.NewMsgExecuteSPL(
		chainID,
		nonce,
		amount,
		decimals,
		mintAccount,
		to,
		destinationProgramPdaAta,
		sender,
		data,
		executeType,
		remainingAccounts,
	)
	msgHash := msg.Hash()

	// prepare increment_nonce msg and compute hash, it will be used as fallback tx in case execute spl fails
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

// createExecuteSPLInstruction wraps the execute spl 'msg' into a Solana instruction.
func (signer *Signer) createExecuteSPLInstruction(msg contracts.MsgExecuteSPL) (*solana.GenericInstruction, error) {
	// create execute spl instruction with program call data
	var dataBytes []byte
	if msg.ExecuteType() == contracts.ExecuteTypeRevert {
		serializedInst, err := borsh.Serialize(contracts.ExecuteSPLRevertInstructionParams{
			Discriminator: contracts.DiscriminatorExecuteSPLRevert,
			Decimals:      msg.Decimals(),
			Amount:        msg.Amount(),
			Sender:        solana.MustPublicKeyFromBase58(msg.Sender()),
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

	pdaAta, _, err := solana.FindAssociatedTokenAddress(signer.pda, msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", signer.pda, msg.MintAccount())
	}

	destinationProgramPda, err := contracts.ComputeConnectedPdaAddress(msg.To())
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode connected spl pda address")
	}

	predefinedAccounts := []*solana.AccountMeta{
		solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
		solana.Meta(signer.pda).WRITE(),
		solana.Meta(pdaAta).WRITE(),
		solana.Meta(msg.MintAccount()),
		solana.Meta(msg.To()),
		solana.Meta(destinationProgramPda).WRITE(),
		solana.Meta(msg.RecipientAta()).WRITE(),
		solana.Meta(solana.TokenProgramID),
		solana.Meta(solana.SPLAssociatedTokenAccountProgramID),
		solana.Meta(solana.SystemProgramID),
	}
	allAccounts := append(predefinedAccounts, msg.RemainingAccounts()...)

	inst := &solana.GenericInstruction{
		ProgID:        signer.gatewayID,
		DataBytes:     dataBytes,
		AccountValues: allAccounts,
	}

	return inst, nil
}
