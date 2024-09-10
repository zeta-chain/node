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

// SignMsgWithdraw signs a withdraw message (for gateway withdraw/withdraw_spl instruction) with TSS.
func (signer *Signer) SignMsgWithdraw(
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

	// zero out the amount if cancelTx is set. It's legal to withdraw 0 lamports thru the gateway.
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
	signature, err := signer.TSS().Sign(ctx, msgHash[:], height, nonce, chain.ChainId, "")
	if err != nil {
		return nil, errors.Wrap(err, "Key-sign failed")
	}
	signer.Logger().Std.Info().Msgf("Key-sign succeed for chain %d nonce %d", chainID, nonce)

	// attach the signature and return
	return msg.SetSignature(signature), nil
}

// SignWithdrawTx wraps the withdraw 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) SignWithdrawTx(ctx context.Context, msg contracts.MsgWithdraw) (*solana.Transaction, error) {
	// create withdraw instruction with program call data
	var err error
	var inst solana.GenericInstruction
	inst.DataBytes, err = borsh.Serialize(contracts.WithdrawInstructionParams{
		Discriminator: contracts.DiscriminatorWithdraw(),
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize withdraw instruction")
	}

	// attach required accounts to the instruction
	privkey := signer.relayerKey
	attachWithdrawAccounts(&inst, privkey.PublicKey(), signer.pda, msg.To(), signer.gatewayID)

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

// attachWithdrawAccounts attaches the required accounts for the gateway withdraw instruction.
func attachWithdrawAccounts(
	inst *solana.GenericInstruction,
	signer solana.PublicKey,
	pda solana.PublicKey,
	to solana.PublicKey,
	gatewayID solana.PublicKey,
) {
	// attach required accounts to the instruction
	var accountSlice []*solana.AccountMeta
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pda).WRITE())
	accountSlice = append(accountSlice, solana.Meta(to).WRITE())
	inst.ProgID = gatewayID

	inst.AccountValues = accountSlice
}
