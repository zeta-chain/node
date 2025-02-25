package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
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

// createIncrementNonceInstruction wraps the increment_nonce 'msg' into a Solana instruction.
func (signer *Signer) createIncrementNonceInstruction(
	msg contracts.MsgIncrementNonce,
) (*solana.GenericInstruction, error) {
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

	inst := &solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
		},
	}

	return inst, nil
}
