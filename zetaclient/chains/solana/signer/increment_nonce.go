package signer

import (
	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

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
