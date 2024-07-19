package rpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/zetaclient/chains/solana/rpc"
)

// Test_SolanaRPCLive is a phony test to run all live tests
func Test_SolanaRPCLive(t *testing.T) {
	// LiveTest_GetFirstSignatureForAddress(t)
	// LiveTest_GetSignaturesForAddressUntil(t)
}

func LiveTest_GetFirstSignatureForAddress(t *testing.T) {
	// create a Solana devnet RPC client
	urlDevnet := os.Getenv("TEST_SOL_URL_DEVNET")
	client := solanarpc.New(urlDevnet)

	// program address
	address := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")

	// get the first signature for the address (one by one)
	sig, err := rpc.GetFirstSignatureForAddress(context.TODO(), client, address, 1)
	require.NoError(t, err)

	// assert
	actualSig := "2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu"
	require.Equal(t, actualSig, sig.String())
}

func LiveTest_GetSignaturesForAddressUntil(t *testing.T) {
	// create a Solana devnet RPC client
	urlDevnet := os.Getenv("TEST_SOL_URL_DEVNET")
	client := solanarpc.New(urlDevnet)

	// program address
	address := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")
	untilSig := solana.MustSignatureFromBase58(
		"2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu",
	)

	// get all signatures for the address until the first signature (one by one)
	sigs, err := rpc.GetSignaturesForAddressUntil(context.TODO(), client, address, untilSig, 1)
	require.NoError(t, err)

	// assert
	require.Greater(t, len(sigs), 0)

	// untilSig should not be in the list
	for _, sig := range sigs {
		require.NotEqual(t, untilSig, sig.Signature)
	}
}
