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

	// parse mint account
	mintAccount, err := solana.PublicKeyFromBase58(asset)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse asset public key %s", asset)
	}

	// get recipient ata
	recipientAta, _, err := solana.FindAssociatedTokenAddress(to, mintAccount)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", to, mintAccount)
	}

	// prepare withdraw spl msg and compute hash
	msg := contracts.NewMsgWithdrawSPL(chainID, nonce, amount, decimals, mintAccount, to, recipientAta)
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

// createWithdrawSPLInstruction wraps the withdraw spl 'msg' into a Solana instruction.
func (signer *Signer) createWithdrawSPLInstruction(msg contracts.MsgWithdrawSPL) (*solana.GenericInstruction, error) {
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

	pdaAta, _, err := solana.FindAssociatedTokenAddress(signer.pda, msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", signer.pda, msg.MintAccount())
	}

	recipientAta, _, err := solana.FindAssociatedTokenAddress(msg.To(), msg.MintAccount())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find ATA for %s and mint account %s", msg.To(), msg.MintAccount())
	}

	inst := &solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(pdaAta).WRITE(),
			solana.Meta(msg.MintAccount()),
			solana.Meta(msg.To()),
			solana.Meta(recipientAta).WRITE(),
			solana.Meta(solana.TokenProgramID),
			solana.Meta(solana.SPLAssociatedTokenAccountProgramID),
			solana.Meta(solana.SystemProgramID),
		},
	}

	return inst, nil
}
