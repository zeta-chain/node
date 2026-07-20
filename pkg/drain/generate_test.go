package drain_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
)

func testPayee(t *testing.T) btcutil.Address {
	addr, err := btcutil.NewAddressWitnessPubKeyHash(make([]byte, 20), &chaincfg.RegressionNetParams)
	require.NoError(t, err)
	return addr
}

func TestGenerateEVMTx(t *testing.T) {
	// ARRANGE
	in := drain.EVMInput{
		ChainID:        1,
		To:             "0x1111111111111111111111111111111111111111",
		Balance:        sdkmath.NewUint(1e18),
		MedianGasPrice: sdkmath.NewUint(100_000),
		Nonce:          7,
	}

	// ACT
	tx, err := drain.GenerateEVMTx(in)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, in.To, tx.To)
	require.EqualValues(t, 7, tx.Nonce)
	require.EqualValues(t, 21_000, tx.GasLimit)
	require.Equal(t, "250000", tx.GasPrice)
	// amount = balance - (21000*250000 + 2_100_000_000)
	require.Equal(t, "999999992650000000", tx.Amount)
}

func TestGenerateBTCTxsPartitioning(t *testing.T) {
	payee := testPayee(t)

	makeUTXOs := func(n int) []drain.UTXO {
		utxos := make([]drain.UTXO, n)
		for i := 0; i < n; i++ {
			utxos[i] = drain.UTXO{TxID: string(rune('a' + i)), Vout: uint32(i), AmountSats: 10_000_000}
		}
		return utxos
	}

	tests := []struct {
		name       string
		numUTXOs   int
		expectTxns int
	}{
		{"empty", 0, 0},
		{"single group", 20, 1},
		{"exactly two groups", 40, 2},
		{"partial last group", 45, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			in := drain.BTCInput{ChainID: 8332, To: payee, FeeRate: 10, UTXOs: makeUTXOs(tc.numUTXOs)}

			// ACT
			txs, err := drain.GenerateBTCTxs(in)

			// ASSERT
			require.NoError(t, err)
			require.Len(t, txs, tc.expectTxns)

			// inputs are disjoint, each group <= 20, and together cover the full set
			seen := map[uint32]bool{}
			totalInputs := 0
			for _, tx := range txs {
				require.LessOrEqual(t, len(tx.Inputs), drain.MaxInputsPerTx)
				var sumInputs int64
				for _, in := range tx.Inputs {
					require.False(t, seen[in.Vout], "input reused across txs")
					seen[in.Vout] = true
					sumInputs += in.AmountSats
					totalInputs++
				}
				// output = sum(inputs) - fee
				require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
				require.Positive(t, tx.FeeSats)
			}
			require.Equal(t, tc.numUTXOs, totalInputs)
		})
	}
}

func TestBuildPayloadSigns(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	evmTx, err := drain.GenerateEVMTx(drain.EVMInput{
		ChainID: 1, To: "0x1111111111111111111111111111111111111111",
		Balance: sdkmath.NewUint(1e18), MedianGasPrice: sdkmath.NewUint(100_000), Nonce: 0,
	})
	require.NoError(t, err)

	// ACT
	p, err := drain.BuildPayload(100, 1, true, []draintx.EVMTx{evmTx}, nil, priv)

	// ASSERT
	require.NoError(t, err)
	require.True(t, p.Final)
	require.NoError(t, p.Verify(ethcrypto.CompressPubkey(&priv.PublicKey)))
}
