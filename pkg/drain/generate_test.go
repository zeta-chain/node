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
	// the drain pins DrainEVMGasLimit (100000), not gas.EVMSend, so a contract safe wallet has room
	require.EqualValues(t, 100_000, tx.GasLimit)
	require.Equal(t, "250000", tx.GasPrice)
	// amount = balance - (100000*250000 + 2_100_000_000)
	require.Equal(t, "999999972900000000", tx.Amount)
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
				require.LessOrEqual(t, len(tx.Inputs), btccommon.MaxNoOfInputsPerTx)
				var sumInputs int64
				for _, in := range tx.Inputs {
					require.False(t, seen[in.Vout], "input reused across txs")
					seen[in.Vout] = true
					sumInputs += in.AmountSats
					totalInputs++
				}
				// output = sum(inputs) - miner-fee-only, right-sized to the input count
				wantSize, err := btccommon.EstimateOutboundSize(int64(len(tx.Inputs)), []btcutil.Address{payee})
				require.NoError(t, err)
				require.Equal(t, in.FeeRate*wantSize, tx.FeeSats)
				require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
			}
			// every UTXO here is economical, so all are swept
			require.Equal(t, tc.numUTXOs, totalInputs)
		})
	}
}

func TestGenerateBTCTxsSkipsDust(t *testing.T) {
	payee := testPayee(t)
	const feeRate = int64(50)

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
	require.Len(t, tx.Inputs, btccommon.MaxNoOfInputsPerTx)
	// sort desc: the three large UTXOs lead the swept group
	require.EqualValues(t, 100_000_000, tx.Inputs[0].AmountSats)
	require.EqualValues(t, 100_000_000, tx.Inputs[1].AmountSats)
	require.EqualValues(t, 100_000_000, tx.Inputs[2].AmountSats)
	wantSize, err := btccommon.EstimateOutboundSize(int64(len(tx.Inputs)), []btcutil.Address{payee})
	require.NoError(t, err)
	require.Equal(t, feeRate*wantSize, tx.FeeSats)
	var sumInputs int64
	for _, in := range tx.Inputs {
		sumInputs += in.AmountSats
	}
	require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
	require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
}

func TestGenerateBTCTxsSkipsHighFeeGroup(t *testing.T) {
	payee := testPayee(t)

	// a fee rate high enough that the right-sized fee exceeds 1/MaxBTCFeeFraction of the input.
	// The poller's validateBTCFee would reject such a sweep, so the generator must not emit it.
	size, err := btccommon.EstimateOutboundSize(1, []btcutil.Address{payee})
	require.NoError(t, err)
	feeRate := int64(1000)
	fee := size * feeRate
	// total = 2*fee: the output (fee) stays well above dust and the fee is coverable, but
	// fee > total/MaxBTCFeeFraction (= fee/2), so the group is uneconomical.
	total := 2 * fee
	require.Greater(t, total-fee, int64(constant.BTCWithdrawalDustAmount))

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332,
		To:      payee,
		FeeRate: feeRate,
		UTXOs:   []drain.UTXO{{TxID: "aaa", Vout: 0, AmountSats: total}},
	})

	// ASSERT: no tx emitted, and no error (skipped like other uneconomical groups)
	require.NoError(t, err)
	require.Empty(t, txs)
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
	p, err := drain.BuildPayload(100, 1, true, drain.NetworkMainnet, []draintx.EVMTx{evmTx}, nil, priv)

	// ASSERT
	require.NoError(t, err)
	require.True(t, p.Final)
	require.Equal(t, drain.NetworkMainnet, p.Network)
	require.NoError(t, p.Verify(ethcrypto.CompressPubkey(&priv.PublicKey)))
}

func TestGenerateEVMTxMaxAmount(t *testing.T) {
	// the uncapped amount for this input, from TestGenerateEVMTx
	const uncapped = "999999972900000000"

	baseInput := func() drain.EVMInput {
		return drain.EVMInput{
			ChainID:        1,
			To:             "0x1111111111111111111111111111111111111111",
			Balance:        sdkmath.NewUint(1e18),
			MedianGasPrice: sdkmath.NewUint(100_000),
			Nonce:          7,
		}
	}

	tests := []struct {
		name       string
		maxAmount  sdkmath.Uint
		wantAmount string
	}{
		{"nil cap drains everything", sdkmath.Uint{}, uncapped},
		{"zero cap drains everything", sdkmath.ZeroUint(), uncapped},
		{"cap below amount is applied", sdkmath.NewUint(1_000), "1000"},
		{"cap above amount is a no-op", sdkmath.NewUintFromString("2000000000000000000"), uncapped},
		{"cap equal to amount is a no-op", sdkmath.NewUintFromString(uncapped), uncapped},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			in := baseInput()
			in.MaxAmount = tc.maxAmount

			// ACT
			tx, err := drain.GenerateEVMTx(in)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, tc.wantAmount, tx.Amount)
			// the cap only lowers the transfer: fee inputs are untouched, so a capped rehearsal
			// signs and broadcasts exactly like the real drain
			require.Equal(t, "250000", tx.GasPrice)
			require.EqualValues(t, 100_000, tx.GasLimit)
			require.EqualValues(t, 7, tx.Nonce)
		})
	}
}

func TestGenerateEVMTxRefusesZeroAmount(t *testing.T) {
	// a balance that is exactly the fee leaves nothing to transfer; pinning it would pay gas
	// to move nothing. fee = 100000*250000 + 2_100_000_000
	fee := sdkmath.NewUint(100_000 * 250_000).Add(sdkmath.NewUintFromString("2100000000"))

	// ACT
	_, err := drain.GenerateEVMTx(drain.EVMInput{
		ChainID:        1,
		To:             "0x1111111111111111111111111111111111111111",
		Balance:        fee,
		MedianGasPrice: sdkmath.NewUint(100_000),
	})

	// ASSERT
	require.ErrorContains(t, err, "zero-amount")
}

func TestGenerateBTCTxsMaxSats(t *testing.T) {
	payee := testPayee(t)

	utxos := []drain.UTXO{
		{TxID: "aaa", Vout: 0, AmountSats: 100_000_000},
		{TxID: "bbb", Vout: 1, AmountSats: 60_000_000},
		{TxID: "ccc", Vout: 2, AmountSats: 30_000_000},
		{TxID: "ddd", Vout: 3, AmountSats: 5_000_000},
	}

	// ARRANGE: a cap that fits the 60M UTXO but not the 100M one; 30M then tops it up, and 5M
	// fits in what remains — largest-first, never exceeding the cap.
	const maxSats = int64(95_000_000)

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332,
		To:      payee,
		FeeRate: 10,
		UTXOs:   utxos,
		MaxSats: maxSats,
	})

	// ASSERT: one sweep, spending only the selected subset, within the cap
	require.NoError(t, err)
	require.Len(t, txs, 1)

	tx := txs[0]
	var totalIn int64
	spent := make([]string, 0, len(tx.Inputs))
	for _, in := range tx.Inputs {
		totalIn += in.AmountSats
		spent = append(spent, in.TxID)
	}
	require.LessOrEqual(t, totalIn, maxSats)
	require.ElementsMatch(t, []string{"bbb", "ccc", "ddd"}, spent)
	// the untouched balance stays at the TSS address for the real drain
	require.NotContains(t, spent, "aaa")
	// the fee math and single-output shape are unchanged by the cap
	require.Equal(t, totalIn-tx.FeeSats, tx.OutputSats)
	require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
}

func TestGenerateBTCTxsMaxSatsPrefersEconomicalInputs(t *testing.T) {
	payee := testPayee(t)
	const feeRate = int64(50)

	// one UTXO that comfortably fits the cap, plus enough dust to fill a whole tx. A
	// smallest-first selection would pack the sweep with dust and the group would be dropped
	// as uneconomical, so a rehearsal would silently produce no BTC tx at all.
	utxos := []drain.UTXO{{TxID: "large", Vout: 0, AmountSats: 10_000_000}}
	for i := 0; i < 40; i++ {
		utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("dust-%02d", i), Vout: uint32(100 + i), AmountSats: 500})
	}

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332,
		To:      payee,
		FeeRate: feeRate,
		UTXOs:   utxos,
		MaxSats: 20_000_000,
	})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Equal(t, "large", txs[0].Inputs[0].TxID)
	require.GreaterOrEqual(t, txs[0].OutputSats, int64(constant.BTCWithdrawalDustAmount))
}

func TestGenerateBTCTxsMaxSatsCapsInputsToOneTx(t *testing.T) {
	payee := testPayee(t)

	// 45 equal UTXOs all within a generous cap: the selection stops at one tx worth of inputs
	// so a rehearsal never fans out into several sweeps.
	utxos := make([]drain.UTXO, 45)
	for i := range utxos {
		utxos[i] = drain.UTXO{TxID: fmt.Sprintf("utxo-%02d", i), Vout: uint32(i), AmountSats: 10_000_000}
	}

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332,
		To:      payee,
		FeeRate: 10,
		UTXOs:   utxos,
		MaxSats: 450_000_000, // fits every UTXO
	})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Len(t, txs[0].Inputs, btccommon.MaxNoOfInputsPerTx)
}

func TestGenerateBTCTxsMaxSatsBelowSmallestUTXO(t *testing.T) {
	payee := testPayee(t)

	// ACT: no UTXO fits under the cap
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332,
		To:      payee,
		FeeRate: 10,
		UTXOs:   []drain.UTXO{{TxID: "aaa", Vout: 0, AmountSats: 10_000_000}},
		MaxSats: 1_000,
	})

	// ASSERT: nothing to sweep, and no error — the operator sees an empty BTC section rather
	// than a failed payload
	require.NoError(t, err)
	require.Empty(t, txs)
}

func TestGenerateBTCTxsMaxSatsDeterministic(t *testing.T) {
	payee := testPayee(t)

	// a capped selection must be as order-independent as the uncapped partitioning, or nodes
	// would select different UTXOs and the keysign ceremony would never form
	base := []drain.UTXO{
		{TxID: "aaa", Vout: 0, AmountSats: 50_000_000},
		{TxID: "bbb", Vout: 1, AmountSats: 50_000_000}, // ties on amount → tie-break by TxID/Vout
		{TxID: "ccc", Vout: 2, AmountSats: 90_000_000},
		{TxID: "ddd", Vout: 3, AmountSats: 10_000_000},
		{TxID: "aaa", Vout: 4, AmountSats: 50_000_000}, // same TxID as [0], tie-break by Vout
	}
	shuffled := []drain.UTXO{base[3], base[0], base[4], base[2], base[1]}

	txs1, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: 10, UTXOs: base, MaxSats: 110_000_000,
	})
	require.NoError(t, err)
	txs2, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: 10, UTXOs: shuffled, MaxSats: 110_000_000,
	})
	require.NoError(t, err)

	require.NotEmpty(t, txs1)
	require.Equal(t, txs1, txs2)
}

func TestGenerateBTCTxsMaxSatsDoesNotMutateCallerSlice(t *testing.T) {
	payee := testPayee(t)

	// the capped path sorts before selecting; the caller's slice must survive it untouched
	utxos := []drain.UTXO{
		{TxID: "aaa", Vout: 0, AmountSats: 10_000_000},
		{TxID: "bbb", Vout: 1, AmountSats: 90_000_000},
		{TxID: "ccc", Vout: 2, AmountSats: 50_000_000},
	}
	original := append([]drain.UTXO(nil), utxos...)

	// ACT
	_, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: 10, UTXOs: utxos, MaxSats: 60_000_000,
	})

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, original, utxos)
}

// mainnetShapedUTXOs mirrors the shape measured on mainnet during the live run: 20 UTXOs holding
// ~99.75% of the balance, the remainder spread over hundreds of dust outputs. That shape is what
// makes a capped selection hard — there is a wide gap between "large" and "dust" with nothing in
// between for a small cap to land on.
func mainnetShapedUTXOs() []drain.UTXO {
	const totalSats = int64(177_000_000)
	large := totalSats * 9975 / 10000 / 20

	utxos := make([]drain.UTXO, 0, 500)
	for i := 0; i < 20; i++ {
		utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("large-%03d", i), Vout: uint32(i), AmountSats: large})
	}
	for i := 0; i < 480; i++ {
		utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("dust-%03d", i), Vout: uint32(1000 + i), AmountSats: 931})
	}
	return utxos
}

func TestGenerateBTCTxsCapNeverEmitsUneconomicalSweep(t *testing.T) {
	payee := testPayee(t)

	// A cap below every large UTXO used to slide past all of them and fill the group with dust:
	// fitting under the cap was the only test, so the walk kept going down the list. The dust
	// group's fee then dwarfed its value and the whole sweep was dropped, leaving the operator
	// with an empty BTC section. Whatever the cap, an emitted sweep must satisfy the poller's
	// fee bound — the alternative is a rehearsal that silently skips BTC.
	for _, feeRate := range []int64{10, 50} {
		for _, maxSats := range []int64{100_000, 1_000_000, 5_000_000, 9_000_000, 20_000_000} {
			txs, err := drain.GenerateBTCTxs(drain.BTCInput{
				ChainID: 8332,
				To:      payee,
				FeeRate: feeRate,
				UTXOs:   mainnetShapedUTXOs(),
				MaxSats: maxSats,
			})
			require.NoError(t, err)

			for _, tx := range txs {
				var totalIn int64
				for _, in := range tx.Inputs {
					totalIn += in.AmountSats
				}
				require.LessOrEqual(t, totalIn, maxSats, "feeRate=%d cap=%d", feeRate, maxSats)
				// exactly the poller's validateBTCFee bound
				require.LessOrEqual(
					t, tx.FeeSats, totalIn/drain.MaxBTCFeeFraction,
					"feeRate=%d cap=%d emitted a sweep the poller would reject", feeRate, maxSats,
				)
				require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
				require.Equal(t, totalIn-tx.FeeSats, tx.OutputSats)
			}
		}
	}
}

func TestGenerateBTCTxsCapDoesNotTopUpWithDust(t *testing.T) {
	payee := testPayee(t)
	const feeRate = int64(50)

	// A single 90k-sat UTXO is viable on its own at 50 sat/vB (fee 8550 vs a 9000 allowance).
	// Adding one 300-sat input costs ~68 vB more fee while adding almost no value, pushing the
	// fee past the bound and taking the whole sweep down with it. Every input must cover its own
	// marginal fee under the same 1/MaxBTCFeeFraction rule.
	utxos := []drain.UTXO{{TxID: "big", Vout: 0, AmountSats: 90_000}}
	for i := 0; i < 19; i++ {
		utxos = append(utxos, drain.UTXO{TxID: fmt.Sprintf("d%02d", i), Vout: uint32(100 + i), AmountSats: 300})
	}

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: feeRate, UTXOs: utxos, MaxSats: 100_000,
	})

	// ASSERT: the viable sweep survives, and the dust that would have killed it is left behind
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Len(t, txs[0].Inputs, 1)
	require.EqualValues(t, 90_000, txs[0].Inputs[0].AmountSats)
	require.LessOrEqual(t, txs[0].FeeSats, txs[0].Inputs[0].AmountSats/drain.MaxBTCFeeFraction)
}

func TestGenerateBTCTxsCapTopsUpWithInputsThatPayForThemselves(t *testing.T) {
	payee := testPayee(t)

	// the rejection above is about marginal value, not about using one input: an input large
	// enough to cover its own added fee is still taken
	utxos := []drain.UTXO{
		{TxID: "aaa", Vout: 0, AmountSats: 60_000_000},
		{TxID: "bbb", Vout: 1, AmountSats: 30_000_000},
		{TxID: "ccc", Vout: 2, AmountSats: 5_000_000},
	}

	// ACT
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: 10, UTXOs: utxos, MaxSats: 95_000_000,
	})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Len(t, txs[0].Inputs, 3)
}

func TestMinViableSweepSats(t *testing.T) {
	payee := testPayee(t)

	// a sweep below this threshold cannot satisfy both the fee bound and the dust floor, so no
	// cap under it can ever produce a BTC tx
	for _, feeRate := range []int64{1, 10, 50, 200} {
		minViable, err := drain.MinViableSweepSats(feeRate, payee)
		require.NoError(t, err)

		// just under the threshold: nothing is emitted
		txs, err := drain.GenerateBTCTxs(drain.BTCInput{
			ChainID: 8332,
			To:      payee,
			FeeRate: feeRate,
			UTXOs:   []drain.UTXO{{TxID: "aaa", Vout: 0, AmountSats: minViable - 1}},
			MaxSats: minViable - 1,
		})
		require.NoError(t, err)
		require.Empty(t, txs, "feeRate=%d", feeRate)

		// exactly at the threshold: a single-input sweep is emitted and passes the poller's bounds
		txs, err = drain.GenerateBTCTxs(drain.BTCInput{
			ChainID: 8332,
			To:      payee,
			FeeRate: feeRate,
			UTXOs:   []drain.UTXO{{TxID: "aaa", Vout: 0, AmountSats: minViable}},
			MaxSats: minViable,
		})
		require.NoError(t, err)
		require.Len(t, txs, 1, "feeRate=%d minViable=%d", feeRate, minViable)
		require.LessOrEqual(t, txs[0].FeeSats, minViable/drain.MaxBTCFeeFraction)
		require.GreaterOrEqual(t, txs[0].OutputSats, int64(constant.BTCWithdrawalDustAmount))
	}
}

func TestGenerateBTCTxsCapStopsAtTheFirstUnaffordableInput(t *testing.T) {
	payee := testPayee(t)
	const feeRate = int64(50)

	// Selection stops at the first candidate that can't cover its own marginal fee rather than
	// scanning on. With a descending list that is not merely an optimisation: the fee is fixed by
	// the current input count, so no smaller candidate can pass the same test. Anything that did
	// get taken after such a candidate would mean the ordering or the bound had been broken.
	utxos := []drain.UTXO{
		{TxID: "big", Vout: 0, AmountSats: 90_000},
		{TxID: "mid", Vout: 1, AmountSats: 300},
		// a UTXO large enough to pay its way, but ordered behind the one that stops the walk
		{TxID: "late", Vout: 2, AmountSats: 200},
	}

	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: 8332, To: payee, FeeRate: feeRate, UTXOs: utxos, MaxSats: 100_000,
	})
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Len(t, txs[0].Inputs, 1)
	require.EqualValues(t, 90_000, txs[0].Inputs[0].AmountSats)
}
