package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

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

// signIncrementNonceTx wraps the increment_nonce 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) signIncrementNonceTx(
	ctx context.Context,
	msg contracts.MsgIncrementNonce,
) (*solana.Transaction, error) {
	// create increment_nonce instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.IncrementNonceInstructionParams{
		Discriminator: contracts.DiscriminatorIncrementNonce,
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize increment_nonce instruction")
	}

	inst := solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
		},
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
		return nil, errors.Wrap(err, "unable to create new tx")
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
