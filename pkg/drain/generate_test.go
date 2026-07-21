package drain_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
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
				// output = sum(inputs) - miner-fee-only (no RBF / nonce reserve)
				require.Equal(t, in.FeeRate*btccommon.OutboundBytesMax, tx.FeeSats)
				require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
			}
			// every UTXO here is economical, so all are swept
			require.Equal(t, tc.numUTXOs, totalInputs)
		})
	}
}

func TestGenerateBTCTxsSkipsDust(t *testing.T) {
	payee := testPayee(t)
	const feeRate = int64(50) // miner fee = 50 * 1543 = 77_150 sats per tx

	// 3 large + 40 dust: sorted desc, the large ones head the first group (swept) and the dust
	// clusters into trailing groups that cannot cover the fee (skipped, not aborted).
	makeUTXOs := func() []drain.UTXO {
		utxos := make([]drain.UTXO, 0, 43)
		for i := 0; i < 3; i++ {
			utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("large-%02d", i), Vout: uint32(i), AmountSats: 100_000_000})
		}
		for i := 0; i < 40; i++ {
			utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("dust-%02d", i), Vout: uint32(100 + i), AmountSats: 500})
		}
		return utxos
	}

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{ChainID: 8332, To: payee, FeeRate: feeRate, UTXOs: makeUTXOs()})

	// ASSERT: skipping dust groups must not abort the sweep
	require.NoError(t, err)
	require.Len(t, txs, 1)

	tx := txs[0]
	require.Len(t, tx.Inputs, drain.MaxInputsPerTx)
	// sort desc: the three large UTXOs lead the swept group
	require.EqualValues(t, 100_000_000, tx.Inputs[0].AmountSats)
	require.EqualValues(t, 100_000_000, tx.Inputs[1].AmountSats)
	require.EqualValues(t, 100_000_000, tx.Inputs[2].AmountSats)
	require.Equal(t, feeRate*btccommon.OutboundBytesMax, tx.FeeSats)
	var sumInputs int64
	for _, in := range tx.Inputs {
		sumInputs += in.AmountSats
	}
	require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
	require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
}

func TestGenerateBTCTxsDeterministic(t *testing.T) {
	payee := testPayee(t)

	// same UTXOs, two different input orderings, must produce byte-identical txs
	base := []drain.UTXO{
		{TxID: "aaa", Vout: 0, AmountSats: 50_000_000},
		{TxID: "bbb", Vout: 1, AmountSats: 50_000_000}, // ties on amount → tie-break by TxID/Vout
		{TxID: "ccc", Vout: 2, AmountSats: 90_000_000},
		{TxID: "ddd", Vout: 3, AmountSats: 10_000_000},
		{TxID: "aaa", Vout: 4, AmountSats: 50_000_000}, // same TxID as [0], tie-break by Vout
	}
	shuffled := []drain.UTXO{base[3], base[0], base[4], base[2], base[1]}

	in1 := drain.BTCInput{ChainID: 8332, To: payee, FeeRate: 10, UTXOs: base}
	in2 := drain.BTCInput{ChainID: 8332, To: payee, FeeRate: 10, UTXOs: shuffled}

	txs1, err := drain.GenerateBTCTxs(in1)
	require.NoError(t, err)
	txs2, err := drain.GenerateBTCTxs(in2)
	require.NoError(t, err)

	require.Equal(t, txs1, txs2)
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
