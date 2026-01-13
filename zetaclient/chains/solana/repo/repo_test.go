package repo_test

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

// Test_SolanaRepoLive is a phony test to run all live tests
func Test_SolanaRepoLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		return
	}

	t.Run("GetTransactionWithVersion", func(t *testing.T) {
		Run_GetTransactionWithVersion(t)
	})

	t.Run("GetFirstSignatureForAddress", func(t *testing.T) {
		Run_GetFirstSignatureForAddress(t)
	})

	t.Run("GetSignaturesForAddressUntil", func(t *testing.T) {
		Run_GetSignaturesForAddressUntil(t)
	})

	t.Run("GetSignaturesForAddressUntil_Version0", func(t *testing.T) {
		Run_GetSignaturesForAddressUntil_Version0(t)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		Run_HealthCheck(t)
	})
}

func Run_GetTransactionWithVersion(t *testing.T) {
	// create a Solana devnet RPC client
	client := solrpc.New(solrpc.DevNet_RPC)

	repo := repo.New(client)

	// example transaction of version "0"
	// https://explorer.solana.com/tx/Wqgj7hAaUUSfLzieN912G7GxyGHijzBZgY135NtuFtPRjevK8DnYjWwQZy7LAKFQZu582wsjuab2QP27VMUJzAi?cluster=devnet
	txSig := solana.MustSignatureFromBase58(
		"Wqgj7hAaUUSfLzieN912G7GxyGHijzBZgY135NtuFtPRjevK8DnYjWwQZy7LAKFQZu582wsjuab2QP27VMUJzAi",
	)

	t.Run("should get the transaction if the version is supported", func(t *testing.T) {
		ctx := context.Background()
		txResult, err := repo.GetTransaction(ctx, txSig)
		require.NoError(t, err)
		require.NotNil(t, txResult)
	})
}

func Run_GetFirstSignatureForAddress(t *testing.T) {
	// create a Solana devnet RPC client
	client := solrpc.New(solrpc.DevNet_RPC)

	repo := repo.New(client)

	// program address
	address := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")

	// get the first signature for the address (one by one)
	sig, err := repo.GetFirstSignatureForAddress(context.Background(), address, 1)
	require.NoError(t, err)

	// assert
	actualSig := "2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu"
	require.Equal(t, actualSig, sig.String())
}

func Run_GetSignaturesForAddressUntil(t *testing.T) {
	// create a Solana devnet RPC client
	client := solrpc.New(solrpc.DevNet_RPC)

	repo := repo.New(client)

	// program address
	address := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")
	untilSig := solana.MustSignatureFromBase58(
		"2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu",
	)

	// get all signatures for the address until the first signature
	sigs, err := repo.GetSignaturesForAddressUntil(context.Background(), address, untilSig, 100)
	require.NoError(t, err)

	// assert
	require.Greater(t, len(sigs), 0)

	// untilSig should not be in the list
	for _, sig := range sigs {
		require.NotEqual(t, untilSig, sig.Signature)
	}
}

func Run_GetSignaturesForAddressUntil_Version0(t *testing.T) {
	// create a Solana devnet RPC client
	client := solrpc.New(solrpc.DevNet_RPC)

	repo := repo.New(client)

	// program address and signature of version "0"
	chain := chains.SolanaDevnet
	address := solana.MustPublicKeyFromBase58(testutils.GatewayAddresses[chain.ChainId])
	untilSig := solana.MustSignatureFromBase58(
		"39fSgD2nteJCQRQP3ynqEcDMZAFSETCbfb61AUqLU6y7qbzSJL5rn2DHU2oM35zsf94Feb5C5QWd5L5UnncBsAay",
	)

	// should get all signatures for the address until a signature of version "0" successfully
	_, err := repo.GetSignaturesForAddressUntil(context.Background(), address, untilSig, 100)
	require.NoError(t, err)
}

func Run_HealthCheck(t *testing.T) {
	// create a Solana devnet RPC client
	client := solrpc.New(solrpc.DevNet_RPC)

	repo := repo.New(client)

	// check the RPC status
	ctx := context.Background()
	_, err := repo.HealthCheck(ctx)
	require.NoError(t, err)
}
