package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// createAndSignMsgWhitelist creates and signs a whitelist message (for gateway whitelist_spl_mint instruction) with TSS.
func (signer *Signer) createAndSignMsgWhitelist(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	whitelistCandidate solana.PublicKey,
	whitelistEntry solana.PublicKey,
) (*contracts.MsgWhitelist, error) {
	chain := signer.Chain()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce

	// prepare whitelist msg and compute hash
	msg := contracts.NewMsgWhitelist(whitelistCandidate, whitelistEntry, chainID, nonce)
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

// createWhitelistInstruction wraps the whitelist 'msg' into a Solana instruction.
func (signer *Signer) createWhitelistInstruction(msg *contracts.MsgWhitelist) (*solana.GenericInstruction, error) {
	// create whitelist_spl_mint instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.WhitelistInstructionParams{
		Discriminator: contracts.DiscriminatorWhitelistSplMint,
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize whitelist_spl_mint instruction")
	}

	inst := &solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(msg.WhitelistEntry()).WRITE(),
			solana.Meta(msg.WhitelistCandidate()),
			solana.Meta(solana.SystemProgramID),
		},
	}

	return inst, nil
}
