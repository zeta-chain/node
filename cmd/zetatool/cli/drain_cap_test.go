package cli

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
)

func TestParseEVMMaxAmount(t *testing.T) {
	t.Run("empty means no cap", func(t *testing.T) {
		v, err := parseEVMMaxAmount("")
		require.NoError(t, err)
		require.True(t, v.IsZero())
	})

	t.Run("whitespace is trimmed", func(t *testing.T) {
		v, err := parseEVMMaxAmount("  1000  ")
		require.NoError(t, err)
		require.Equal(t, sdkmath.NewUint(1000), v)
	})

	t.Run("parses a wei amount beyond int64", func(t *testing.T) {
		v, err := parseEVMMaxAmount("100000000000000000000")
		require.NoError(t, err)
		require.Equal(t, "100000000000000000000", v.String())
	})

	t.Run("an explicit zero is rejected", func(t *testing.T) {
		// a zero cap would pin a transfer that pays gas to move nothing; the operator almost
		// certainly meant to omit the flag
		_, err := parseEVMMaxAmount("0")
		require.ErrorContains(t, err, "must be positive")
	})

	t.Run("a value above 256 bits is rejected, not panicked on", func(t *testing.T) {
		// sdkmath.Uint's constructor panics past 256 bits, so an extra-zeros typo would take
		// down zetatool mid-drain instead of reporting a bad flag
		maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

		v, err := parseEVMMaxAmount(maxUint256.String())
		require.NoError(t, err)
		require.Equal(t, maxUint256.String(), v.String())

		require.NotPanics(t, func() {
			_, err = parseEVMMaxAmount(new(big.Int).Add(maxUint256, big.NewInt(1)).String())
		})
		require.ErrorContains(t, err, "exceeds the maximum 256-bit")
	})

	t.Run("rejects non-decimal and negative values", func(t *testing.T) {
		for _, s := range []string{"0x64", "1e18", "abc", "-1", "1.5"} {
			_, err := parseEVMMaxAmount(s)
			require.ErrorContains(t, err, "invalid --evm-max-amount", "input %q", s)
		}
	})
}

func TestWarnIfCapped(t *testing.T) {
	t.Run("silent when nothing is capped", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 0)
		require.Empty(t, buf.String())
	})

	t.Run("an unset cap reads as uncapped, not as zero", func(t *testing.T) {
		// "btc 0 sats total" would read as "sends no BTC", which is the opposite of the truth:
		// BTC is swept in full when only the EVM cap is set
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 0)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL PAYLOAD")
		require.Contains(t, out, "evm: 1000 wei per chain")
		require.Contains(t, out, "btc: uncapped")
	})

	t.Run("btc-only cap", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 100_000)

		out := buf.String()
		require.Contains(t, out, "evm: uncapped")
		require.Contains(t, out, "btc: 100000 sats total")
	})

	t.Run("does not claim nothing is drained", func(t *testing.T) {
		// a cap is an upper bound: a chain holding less than it is swept in full
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 100_000)
		require.Contains(t, buf.String(), "except on chains already below the cap")
	})
}

func TestReportUnviableBTCCap(t *testing.T) {
	payee, err := btcutil.NewAddressWitnessPubKeyHash(make([]byte, 20), &chaincfg.RegressionNetParams)
	require.NoError(t, err)

	const feeRate = int64(10)
	minViable, err := drain.MinViableSweepSats(feeRate, payee)
	require.NoError(t, err)

	t.Run("reports the floor on the sweep total", func(t *testing.T) {
		// mainnet-shaped: large UTXOs far above the cap, the rest dust
		utxos := []drain.UTXO{
			{TxID: "large", Vout: 0, AmountSats: 8_828_175},
			{TxID: "dust", Vout: 1, AmountSats: 931},
		}

		var buf bytes.Buffer
		reportUnviableBTCCap(&buf, utxos, 100_000, feeRate, payee, 8332)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL SWEPT NO BTC")
		require.Contains(t, out, fmt.Sprintf("must total at least %d sats", minViable))
		require.Contains(t, out, fmt.Sprintf("raise --btc-max-sats to at least %d", minViable))
		// the floor is on the total, not on any one UTXO: pointing at the smallest sweepable UTXO
		// would overstate the cap needed, since smaller UTXOs can clear the floor as a group
		require.NotContains(t, out, "8828175")
	})

	t.Run("only claims impossibility against the whole balance", func(t *testing.T) {
		// the entire balance is below the floor, which is the one sound impossibility test
		utxos := []drain.UTXO{
			{TxID: "dust-0", Vout: 0, AmountSats: 931},
			{TxID: "dust-1", Vout: 1, AmountSats: 500},
		}

		var buf bytes.Buffer
		reportUnviableBTCCap(&buf, utxos, 100_000, feeRate, payee, 8332)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL SWEPT NO BTC")
		require.Contains(t, out, "the entire TSS balance is 1431 sats")
		require.Contains(t, out, "no cap can rehearse BTC")
	})

	t.Run("does not claim impossibility when a group could clear the floor", func(t *testing.T) {
		// no single UTXO reaches the floor, but together they exceed it — the old per-UTXO test
		// declared this wallet un-rehearsable, which was wrong
		utxos := []drain.UTXO{
			{TxID: "a", Vout: 0, AmountSats: 15_000},
			{TxID: "b", Vout: 1, AmountSats: 15_000},
		}
		require.Less(t, utxos[0].AmountSats, minViable)
		require.Greater(t, utxos[0].AmountSats+utxos[1].AmountSats, minViable)

		var buf bytes.Buffer
		reportUnviableBTCCap(&buf, utxos, 10_000, feeRate, payee, 8332)

		out := buf.String()
		require.NotContains(t, out, "no cap can rehearse BTC")
		require.Contains(t, out, fmt.Sprintf("raise --btc-max-sats to at least %d", minViable))
	})
}

func TestReportNonceState(t *testing.T) {
	t.Run("silent when quiesced", func(t *testing.T) {
		var buf bytes.Buffer
		reportNonceState(&buf, 1, 42, 42)
		require.Empty(t, buf.String())
	})

	t.Run("reports in-flight txs when pending is ahead", func(t *testing.T) {
		var buf bytes.Buffer
		reportNonceState(&buf, 1, 42, 45)

		out := buf.String()
		require.Contains(t, out, "NONCE NOT QUIESCED")
		require.Contains(t, out, "3 tx(s) are still in flight")
		require.Contains(t, out, "confirmed nonce 42")
	})

	t.Run("a lagging pending view does not underflow", func(t *testing.T) {
		// pending < pinned means a stale RPC view (a load-balanced endpoint answering the two
		// calls from different backends). Subtracting would wrap uint64 and claim ~1.8e19 txs in
		// flight, and "wait for confirmations" would be the wrong advice.
		var buf bytes.Buffer
		reportNonceState(&buf, 1, 45, 42)

		out := buf.String()
		require.NotContains(t, out, "18446744073709551")
		require.NotContains(t, out, "NONCE NOT QUIESCED")
		require.Contains(t, out, "inconsistent nonce view")
		require.Contains(t, out, "load-balanced")
	})
}
