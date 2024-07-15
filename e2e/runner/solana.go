package runner

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solanacontract "github.com/zeta-chain/zetacore/zetaclient/chains/solana/contract"
)

// GatewayProgramID is the program ID for the gateway program
func (r *E2ERunner) GatewayProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(solanacontract.GatewayProgramID)
}

// ComputePdaAddress computes the PDA address for the gateway program
func (r *E2ERunner) ComputePdaAddress() solana.PublicKey {
	seed := []byte(solanacontract.PDASeed)
	GatewayProgramID := solana.MustPublicKeyFromBase58(solanacontract.GatewayProgramID)
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, GatewayProgramID)
	require.NoError(r, err)

	r.Logger.Info("computed pda: %s, bump %d\n", pdaComputed, bump)

	return pdaComputed
}

// CreateSignedTransaction creates a signed transaction from instructions
func (r *E2ERunner) CreateSignedTransaction(
	instructions []solana.Instruction,
	privateKey solana.PrivateKey,
) *solana.Transaction {
	// get a recent blockhash
	recent, err := r.SolanaClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
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
	sig, err := r.SolanaClient.SendTransactionWithOpts(
		context.TODO(),
		tx,
		rpc.TransactionOpts{},
	)
	require.NoError(r, err)
	r.Logger.Info("broadcast success! tx sig %s; waiting for confirmation...", sig)

	// wait for the transaction to be finalized
	var out *rpc.GetTransactionResult
	for {
		time.Sleep(1 * time.Second)
		out, err = r.SolanaClient.GetTransaction(context.TODO(), sig, &rpc.GetTransactionOpts{})
		if err == nil {
			break
		}
	}

	return sig, out
}
