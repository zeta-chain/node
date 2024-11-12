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

// createAndSignMsgWithdrawSPL creates and signs a withdraw spl message for gateway withdraw_spl instruction with TSS.
func (signer *Signer) createAndSignMsgWithdrawSPL(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	asset string,
	decimals uint8,
	cancelTx bool,
) (*contracts.MsgWithdrawSPL, error) {
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
		return nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// parse token account
	tokenAccount, err := solana.PublicKeyFromBase58(asset)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse asset public key %s", asset)
	}

	// get recipient ata
	recipientAta, _, err := solana.FindAssociatedTokenAddress(to, tokenAccount)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and token account %s", to, tokenAccount)
	}

	// prepare withdraw spl msg and compute hash
	msg := contracts.NewMsgWithdrawSPL(chainID, nonce, amount, decimals, tokenAccount, to, recipientAta)
	msgHash := msg.Hash()

	// sign the message with TSS to get an ECDSA signature.
	// the produced signature is in the [R || S || V] format where V is 0 or 1.
	signature, err := signer.TSS().Sign(ctx, msgHash[:], height, nonce, chain.ChainId, "")
	if err != nil {
		return nil, errors.Wrap(err, "Key-sign failed")
	}

	// attach the signature and return
	return msg.SetSignature(signature), nil
}

// signWithdrawSPLTx wraps the withdraw spl 'msg' into a Solana transaction and signs it with the relayer key.
func (signer *Signer) signWithdrawSPLTx(
	ctx context.Context,
	msg contracts.MsgWithdrawSPL,
) (*solana.Transaction, error) {
	// create withdraw spl instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.WithdrawSPLInstructionParams{
		Discriminator: contracts.DiscriminatorWithdrawSPL,
		Decimals:      msg.Decimals(),
		Amount:        msg.Amount(),
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize withdraw instruction")
	}

	pdaAta, _, err := solana.FindAssociatedTokenAddress(signer.pda, msg.TokenAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and token account %s", signer.pda, msg.TokenAccount())
	}

	recipientAta, _, err := solana.FindAssociatedTokenAddress(msg.To(), msg.TokenAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and token account %s", msg.To(), msg.TokenAccount())
	}

	inst := solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(pdaAta).WRITE(),
			solana.Meta(msg.TokenAccount()),
			solana.Meta(msg.To()),
			solana.Meta(recipientAta).WRITE(),
			solana.Meta(signer.rentPayerPda).WRITE(),
			solana.Meta(solana.TokenProgramID),
			solana.Meta(solana.SPLAssociatedTokenAccountProgramID),
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
