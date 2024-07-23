package runner

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	solanacontract "github.com/zeta-chain/zetacore/pkg/contract/solana"
)

// ComputePdaAddress computes the PDA address for the gateway program
func (r *E2ERunner) ComputePdaAddress() solana.PublicKey {
	seed := []byte(solanacontract.PDASeed)
	GatewayProgramID := solana.MustPublicKeyFromBase58(solanacontract.SolanaGatewayProgramID)
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, GatewayProgramID)
	require.NoError(r, err)

	r.Logger.Info("computed pda: %s, bump %d\n", pdaComputed, bump)

	return pdaComputed
}

// CreateDepositInstruction creates a 'deposit' instruction
func (r *E2ERunner) CreateDepositInstruction(
	signer solana.PublicKey,
	receiver ethcommon.Address,
	amount uint64,
) solana.Instruction {
	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgram

	// create 'deposit' instruction
	inst := &solana.GenericInstruction{}
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(programID))
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	var err error
	inst.DataBytes, err = borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: solanacontract.DiscriminatorDeposit(),
		Amount:        amount,
		Memo:          receiver.Bytes(),
	})
	require.NoError(r, err)

	return inst
}

// CreateSignedTransaction creates a signed transaction from instructions
func (r *E2ERunner) CreateSignedTransaction(
	instructions []solana.Instruction,
	privateKey solana.PrivateKey,
) *solana.Transaction {
	// get a recent blockhash
	recent, err := r.SolanaClient.GetRecentBlockhash(r.Ctx, rpc.CommitmentFinalized)
	require.NoError(r, err)

	// create the initialize transaction
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	require.NoError(r, err)

	// sign the initialize transaction
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			return nil
		},
	)
	require.NoError(r, err)

	return tx
}

// BroadcastTxSync broadcasts a transaction and waits for it to be finalized
func (r *E2ERunner) BroadcastTxSync(tx *solana.Transaction) (solana.Signature, *rpc.GetTransactionResult) {
	// broadcast the transaction
	sig, err := r.SolanaClient.SendTransactionWithOpts(r.Ctx, tx, rpc.TransactionOpts{})
	require.NoError(r, err)
	r.Logger.Info("broadcast success! tx sig %s; waiting for confirmation...", sig)

	// wait for the transaction to be finalized
	var out *rpc.GetTransactionResult
	for {
		time.Sleep(1 * time.Second)
		out, err = r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{})
		if err == nil {
			break
		}
	}

	return sig, out
}