package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

var (
	ctx                              = context.Background()
	errTest                          = errors.New("test error")
	slot      uint64                 = 1000
	blockTime solana.UnixTimeSeconds = 12345
	gatewayID solana.PublicKey       // defined in the init() function
)

func init() {
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}
	gatewayID = privateKey.PublicKey()
}

// repo instantiates a new SolanaRepo using the default gatewayID.
func repo(client SolanaClient) *SolanaRepo {
	return New(client, gatewayID)
}

// newSigs creates a slice of transaction-signatures to be used by the tests.
func newSigs(i, j int) []*rpc.TransactionSignature {
	y := j - i + 1
	sigs := make([]*rpc.TransactionSignature, y)
	for x := range y {
		sig := solana.Signature([64]byte{0: byte(x + i), 1: 42})
		sigs[x] = &rpc.TransactionSignature{Signature: sig}
	}
	return sigs
}

func TestHealthCheck(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetSlot", ctx, rpc.CommitmentFinalized).Return(slot, nil)
		client.On("GetBlockTime", ctx, slot).Return(&blockTime, nil)

		gotBlockTime, err := repo(client).HealthCheck(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(blockTime), gotBlockTime.Unix())
	})

	t.Run("Error/Client", func(t *testing.T) {
		t.Run("GetSlot", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetSlot", ctx, rpc.CommitmentFinalized).Return(uint64(0), errTest)

			gotBlockTime, err := repo(client).HealthCheck(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetSlot)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotBlockTime)
		})

		t.Run("GetBlockTime", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetSlot", ctx, rpc.CommitmentFinalized).Return(slot, nil)
			client.On("GetBlockTime", ctx, slot).Return(nil, errTest)

			gotBlockTime, err := repo(client).HealthCheck(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetBlockTime)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotBlockTime)
		})
	})
}

func TestGetSlot(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetSlot", ctx, rpc.CommitmentFinalized).Return(slot, nil)

		gotSlot, err := repo(client).GetSlot(ctx, rpc.CommitmentFinalized)
		require.NoError(t, err)
		require.Equal(t, slot, gotSlot)
	})

	t.Run("Error/Client/GetSlot", func(t *testing.T) {
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetSlot", ctx, rpc.CommitmentConfirmed).Return(uint64(0), errTest)

		gotSlot, err := repo(client).GetSlot(ctx, rpc.CommitmentConfirmed)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrClientGetSlot)
		require.ErrorIs(t, err, errTest)
		require.Zero(t, gotSlot)
	})
}

func TestGetTransaction(t *testing.T) {
	tx := &rpc.GetTransactionResult{Slot: 10}
	sig := solana.Signature([64]byte{5})
	opts := rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0,
	}

	t.Run("Ok", func(t *testing.T) {
		t.Run("WithTransaction", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &opts).Return(tx, nil)

			gotTx, err := repo(client).GetTransaction(ctx, sig, rpc.CommitmentFinalized)
			require.NoError(t, err)
			require.Equal(t, tx, gotTx)
		})

		// TODO: don't know if this case should exist
		t.Run("WithoutTransaction", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &opts).Return(nil, nil)

			gotTx, err := repo(client).GetTransaction(ctx, sig, rpc.CommitmentFinalized)
			require.NoError(t, err)
			require.Equal(t, (*rpc.GetTransactionResult)(nil), gotTx)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Client/GetTransaction", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &opts).Return(nil, errTest)

			gotTx, err := repo(client).GetTransaction(ctx, sig, rpc.CommitmentFinalized)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetTransaction)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotTx)
		})

		t.Run("UnsupportedTxVersion", func(t *testing.T) {
			errTest := fmt.Errorf("%w (code %s)", errTest, errorCodeUnsupportedTransactionVersion)

			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &opts).Return(nil, errTest)

			gotTx, err := repo(client).GetTransaction(ctx, sig, rpc.CommitmentFinalized)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrUnsupportedTxVersion)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotTx)
		})
	})
}

func TestGetPriorityFee(t *testing.T) {
	fees := []rpc.PriorizationFeeResult{
		{Slot: 1, PrioritizationFee: 2},
		{Slot: 2, PrioritizationFee: 4},
		{Slot: 3, PrioritizationFee: 17},
		{Slot: 4, PrioritizationFee: 0},
		{Slot: 5, PrioritizationFee: 50},
		{Slot: 6, PrioritizationFee: 1000},
		{Slot: 7, PrioritizationFee: 6},
		{Slot: 8, PrioritizationFee: 15},
		{Slot: 9, PrioritizationFee: 0},
		{Slot: 10, PrioritizationFee: 0},
	}
	priorityFee := uint64(15) // based on the median (remember that we ignore fees that are zero)

	t.Run("Ok", func(t *testing.T) {
		t.Run("SomeRecentPrioritizationFees", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			var accounts solana.PublicKeySlice = nil
			client.On("GetRecentPrioritizationFees", ctx, accounts).Return(fees, nil)

			gotPriorityFee, err := repo(client).GetPriorityFee(ctx)
			require.NoError(t, err)
			require.Equal(t, priorityFee, gotPriorityFee)
		})

		t.Run("NoRecentPrioritizationFees", func(t *testing.T) {
			fees := []rpc.PriorizationFeeResult{}

			client := mocks.NewSolanaRPCClient(t)
			var accounts solana.PublicKeySlice = nil
			client.On("GetRecentPrioritizationFees", ctx, accounts).Return(fees, nil)

			gotPriorityFee, err := repo(client).GetPriorityFee(ctx)
			require.NoError(t, err)
			require.Equal(t, uint64(0), gotPriorityFee)
		})
	})

	t.Run("Error/Client/GetRecentPrioritizationFees", func(t *testing.T) {
		client := mocks.NewSolanaRPCClient(t)
		var accounts solana.PublicKeySlice = nil
		client.On("GetRecentPrioritizationFees", ctx, accounts).Return(nil, errTest)

		gotPriorityFee, err := repo(client).GetPriorityFee(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrClientGetRecentPrioritizationFees)
		require.ErrorIs(t, err, errTest)
		require.Equal(t, uint64(0), gotPriorityFee)
	})
}

func TestGetFirstSignature(t *testing.T) {
	opts := rpc.GetSignaturesForAddressOpts{Commitment: rpc.CommitmentFinalized}

	t.Run("Ok", func(t *testing.T) {
		t.Run("OneCallToGetSignatures", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)

			sigs := newSigs(0, 2)

			opts1, opts2 := opts, opts
			opts2.Before = sigs[2].Signature

			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts1).Return(sigs, nil)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts2).Return(nil, nil)

			gotSig, err := repo(client).GetFirstSignature(ctx)
			require.NoError(t, err)
			require.Equal(t, sigs[2].Signature, gotSig)
		})

		t.Run("MultipleCallsToGetSignatures", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)

			sigs1 := newSigs(0, 2)
			sigs2 := newSigs(3, 3)

			opts1, opts2, opts3 := opts, opts, opts
			opts2.Before = sigs1[2].Signature
			opts3.Before = sigs2[0].Signature

			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts1).Return(sigs1, nil)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts2).Return(sigs2, nil)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts3).Return(nil, nil)

			gotSig, err := repo(client).GetFirstSignature(ctx)
			require.NoError(t, err)
			require.Equal(t, sigs2[0].Signature, gotSig)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Client/GetSignaturesForAddress", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts).Return(nil, errTest)

			gotSig, err := repo(client).GetFirstSignature(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetSignaturesForAddress)
			require.ErrorIs(t, err, errTest)
			require.True(t, gotSig.IsZero())
		})

		t.Run("FoundNoSignatures", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &opts).Return(nil, nil)

			gotSig, err := repo(client).GetFirstSignature(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrFoundNoSignatures)
			require.True(t, gotSig.IsZero())
		})
	})
}

func TestGetSignaturesSince(t *testing.T) {
	sig := solana.Signature([64]byte{42})

	txOpts := rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0,
	}

	sigs1 := newSigs(0, 0)   // 0
	sigs2 := newSigs(1, 4)   // 1, 2, 3, 4
	sigs3 := newSigs(5, 6)   // 5, 6
	sigs3[1].Signature = sig // note that we are not able to test the "until" feature properly

	sigOpts := rpc.GetSignaturesForAddressOpts{Until: sig, Commitment: rpc.CommitmentFinalized}
	sigOpts1, sigOpts2, sigOpts3, sigOpts4 := sigOpts, sigOpts, sigOpts, sigOpts
	sigOpts2.Before = sigs1[0].Signature
	sigOpts3.Before = sigs2[3].Signature
	sigOpts4.Before = sigs3[1].Signature

	tx := &rpc.GetTransactionResult{Slot: 10}

	t.Run("Ok", func(t *testing.T) {
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetTransaction", ctx, sig, &txOpts).Return(tx, nil)
		client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &sigOpts1).Return(sigs1, nil)
		client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &sigOpts2).Return(sigs2, nil)
		client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &sigOpts3).Return(sigs3, nil)
		client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &sigOpts4).Return(nil, nil)

		gotSigs, err := repo(client).GetSignaturesSince(ctx, sig)
		require.NoError(t, err)
		require.Len(t, gotSigs, len(sigs1)+len(sigs2)+len(sigs3))

		require.Equal(t, sigs1[0], gotSigs[0])
		require.Equal(t, sigs2[0], gotSigs[1])
		require.Equal(t, sigs2[1], gotSigs[2])
		require.Equal(t, sigs2[2], gotSigs[3])
		require.Equal(t, sigs2[3], gotSigs[4])
		require.Equal(t, sigs3[0], gotSigs[5])
		require.Equal(t, sigs3[1], gotSigs[6])
	})

	t.Run("Error/Client", func(t *testing.T) {
		t.Run("GetTransaction", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &txOpts).Return(nil, errTest)

			gotSigs, err := repo(client).GetSignaturesSince(ctx, sig)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetTransaction)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotSigs)
		})

		t.Run("GetSignaturesForAddress", func(t *testing.T) {
			client := mocks.NewSolanaRPCClient(t)
			client.On("GetTransaction", ctx, sig, &txOpts).Return(tx, nil)
			client.On("GetSignaturesForAddressWithOpts", ctx, gatewayID, &sigOpts1).Return(nil, errTest)

			gotSigs, err := repo(client).GetSignaturesSince(ctx, sig)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrClientGetSignaturesForAddress)
			require.ErrorIs(t, err, errTest)
			require.Nil(t, gotSigs)
		})
	})
}

// ------------------------------------------------------------------------------------------------
// Live tests
// ------------------------------------------------------------------------------------------------

func liveRepo(gatewayID solana.PublicKey) *SolanaRepo {
	// Create a Solana devnet RPC client.
	client := rpc.New(rpc.DevNet_RPC)

	return New(client, gatewayID)
}

func TestSolanaRepoLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		return
	}

	gatewayID := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")

	t.Run("HealthCheck", func(t *testing.T) {
		_, err := liveRepo(gatewayID).HealthCheck(ctx)
		require.NoError(t, err)
	})

	t.Run("GetTransaction", func(t *testing.T) {
		// Example transaction (V0): https://explorer.solana.com/tx/BASE58?cluster=devnet
		base58 := "Wqgj7hAaUUSfLzieN912G7GxyGHijzBZgY135NtuFtPRjevK8DnYjWwQZy7LAKFQZu582wsjuab2QP27VMUJzAi"
		sig := solana.MustSignatureFromBase58(base58)

		// Should get the transaction if the version is supported.
		tx, err := liveRepo(gatewayID).GetTransaction(ctx, sig, rpc.CommitmentFinalized)
		require.NoError(t, err)
		require.NotNil(t, tx)
	})

	t.Run("GetFirstSignature", func(t *testing.T) {
		expectedSig := "2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu"
		firstSig, err := liveRepo(gatewayID).GetFirstSignature(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedSig, firstSig.String())
	})

	t.Run("GetSignaturesSince", func(t *testing.T) {
		sinceSig := solana.MustSignatureFromBase58(
			"2tUQtcrXxtNFtV9kZ4kQsmY7snnEoEEArmu9pUptr4UCy8UdbtjPD6UtfEtPJ2qk5CTzZTmLwsbmZdLymcwSUcHu",
		)

		sigs, err := liveRepo(gatewayID).GetSignaturesSince(ctx, sinceSig)
		require.NoError(t, err)
		require.NotEmpty(t, sigs)
		for _, sig := range sigs { // "since" should not be in the list
			require.NotEqual(t, sinceSig, sig.Signature)
		}
	})

	t.Run("GetSignaturesSince_Version0", func(t *testing.T) {
		chainID := chains.SolanaDevnet.ChainId
		gatewayID := solana.MustPublicKeyFromBase58(testutils.GatewayAddresses[chainID])
		sig := solana.MustSignatureFromBase58(
			"39fSgD2nteJCQRQP3ynqEcDMZAFSETCbfb61AUqLU6y7qbzSJL5rn2DHU2oM35zsf94Feb5C5QWd5L5UnncBsAay",
		)
		sigs, err := liveRepo(gatewayID).GetSignaturesSince(ctx, sig)
		require.NoError(t, err)
		require.NotEmpty(t, sigs)
	})
}
