package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// createAndSignMsgWithdraw creates and signs a withdraw message for gateway withdraw instruction with TSS.
func (signer *Signer) createAndSignMsgWithdraw(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	cancelTx bool,
) (*contracts.MsgWithdraw, error) {
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

	// prepare withdraw msg and compute hash
	msg := contracts.NewMsgWithdraw(chainID, nonce, amount, to)
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
