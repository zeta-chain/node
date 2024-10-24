package signer

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

// SignWhitelistTx wraps the whitelist 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) SignWhitelistTx(ctx context.Context, msg *contracts.MsgWhitelist) (*solana.Transaction, error) {
	// create whitelist_spl_mint instruction with program call data
	var err error
	var inst solana.GenericInstruction
	inst.DataBytes, err = borsh.Serialize(contracts.WhitelistInstructionParams{
		Discriminator: contracts.DiscriminatorWhitelistSplMint(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize whitelist_spl_mint instruction")
	}

	// attach required accounts to the instruction
	privkey := signer.relayerKey
	attachWhitelistAccounts(
		&inst,
		privkey.PublicKey(),
		signer.pda,
		msg.WhitelistCandidate(),
		msg.WhitelistEntry(),
		signer.gatewayID,
	)

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
		solana.TransactionPayer(privkey.PublicKey()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewTransaction error")
	}

	// relayer signs the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privkey.PublicKey()) {
			return privkey
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "signer unable to sign transaction")
	}

	return tx, nil
}

// attachWhitelistAccounts attaches the required accounts for the gateway whitelist instruction.
func attachWhitelistAccounts(
	inst *solana.GenericInstruction,
	signer solana.PublicKey,
	pda solana.PublicKey,
	whitelistCandidate solana.PublicKey,
	whitelistEntry solana.PublicKey,
	gatewayID solana.PublicKey,
) {
	// attach required accounts to the instruction
	var accountSlice []*solana.AccountMeta
	accountSlice = append(accountSlice, solana.Meta(whitelistEntry).WRITE())
	accountSlice = append(accountSlice, solana.Meta(whitelistCandidate))
	accountSlice = append(accountSlice, solana.Meta(pda).WRITE())
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	inst.ProgID = gatewayID

	inst.AccountValues = accountSlice
}
