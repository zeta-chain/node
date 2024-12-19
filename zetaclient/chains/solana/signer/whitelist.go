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
		return nil, errors.Wrap(err, "Key-sign failed")
	}

	// attach the signature and return
	return msg.SetSignature(signature), nil
}

// signWhitelistTx wraps the whitelist 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) signWhitelistTx(ctx context.Context, msg *contracts.MsgWhitelist) (*solana.Transaction, error) {
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

	inst := solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(msg.WhitelistEntry()).WRITE(),
			solana.Meta(msg.WhitelistCandidate()),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(solana.SystemProgramID),
		},
	}

	// get a recent blockhash
	recent, err := signer.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, errors.Wrap(err, "GetLatestBlockhash error")
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
		return nil, errors.Wrap(err, "NewTransaction error")
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
