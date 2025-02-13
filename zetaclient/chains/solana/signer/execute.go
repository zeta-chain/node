package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// createAndSignMsgExecute creates and signs a execute message for gateway execute instruction with TSS.
func (signer *Signer) createAndSignMsgExecute(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	sender [20]byte,
	data []byte,
	remainingAccounts []*solana.AccountMeta,
	cancelTx bool,
) (*contracts.MsgExecute, error) {
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
		return nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// prepare execute msg and compute hash
	msg := contracts.NewMsgExecute(chainID, nonce, amount, to, sender, data, remainingAccounts)
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

// signExecuteTx wraps the execute 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) signExecuteTx(ctx context.Context, msg contracts.MsgExecute) (*solana.Transaction, error) {
	// create execute instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.ExecuteInstructionParams{
		Discriminator: contracts.DiscriminatorExecute,
		Amount:        msg.Amount(),
		Sender:        msg.Sender(),
		Data:          msg.Data(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize execute instruction")
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

	inst := solana.GenericInstruction{
		ProgID:        signer.gatewayID,
		DataBytes:     dataBytes,
		AccountValues: allAccounts,
	}

	// get a recent blockhash
	recent, err := signer.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, errors.Wrap(err, "getLatestBlockhash error")
	}

	// create a transaction that wraps the instruction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			// TODO: outbound now uses 5K lamports as the fixed fee, we could explore priority fee and compute budget
			// https://github.com/zeta-chain/node/issues/2599
			// programs.ComputeBudgetSetComputeUnitLimit(computeUnitLimit),
			// programs.ComputeBudgetSetComputeUnitPrice(computeUnitPrice),
			&inst},
		recent.Value.Blockhash,
		solana.TransactionPayer(signer.relayerKey.PublicKey()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "newTransaction error")
	}

	// relayer signs the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(signer.relayerKey.PublicKey()) {
			return signer.relayerKey
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "signer unable to sign transaction")
	}

	return tx, nil
}
