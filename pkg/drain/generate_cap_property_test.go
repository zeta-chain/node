package drain_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/drain"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

// viableSubsets exhaustively answers, for a small UTXO set, whether any non-empty subset could be
// swept at all — cap, the poller's fee bound, and the dust floor — and what the most valuable such
// subset is worth. Exponential, so callers keep the sets small.
func viableSubsets(
	t *testing.T,
	utxos []drain.UTXO,
	maxSats, feeRate int64,
	payee btcutil.Address,
) (exists bool, bestTotal int64) {
	t.Helper()

	for mask := 1; mask < (1 << len(utxos)); mask++ {
		var total, count int64
		for i := range utxos {
			if mask&(1<<i) != 0 {
				total += utxos[i].AmountSats
				count++
			}
		}
		if total > maxSats || count > int64(btccommon.MaxNoOfInputsPerTx) {
			continue
		}
		size, err := btccommon.EstimateOutboundSize(count, []btcutil.Address{payee})
		require.NoError(t, err)
		fee := size * feeRate
		if fee > total/drain.MaxBTCFeeFraction || total-fee < int64(constant.BTCWithdrawalDustAmount) {
			continue
		}
		exists = true
		if total > bestTotal {
			bestTotal = total
		}
	}
	return exists, bestTotal
}

// TestGenerateBTCTxsCapAgainstExhaustiveSearch is the guard that the cap selection heuristic keeps
// pace with what is actually possible. Three separate bugs in this logic all had the same shape —
// a viable subset existed under the cap and the selection emitted nothing, so a rehearsal silently
// covered no BTC — and each was found by review rather than by a test, because the unit tests only
// pinned the shapes already thought of.
//
// So this compares against brute force over every subset on small random sets, asserting:
//
//   - anything emitted is valid: within the cap and inside the poller's fee and dust bounds;
//   - nothing is emitted when no subset could have been swept;
//   - something IS emitted whenever some subset could have been.
//
// The last is the property the bugs violated. Selection is a heuristic, so this is evidence rather
// than proof, but it is the check that fails when the next variant appears. The amounts mix dust,
// medium and wide ranges because the reported failures came from medium UTXOs and dust tails.
//
// Not asserted: that the *most* valuable viable subset is chosen. Across this sweep the heuristic
// picks a smaller viable subset in a few percent of cases, which is a rehearsal sweeping less than
// it could — harmless, and not worth an exhaustive search in production code.
func TestGenerateBTCTxsCapAgainstExhaustiveSearch(t *testing.T) {
	payee := testPayee(t)
	rng := rand.New(rand.NewSource(42)) // #nosec G404 deterministic test vectors, not security

	amount := func() int64 {
		switch rng.Intn(3) {
		case 0:
			return int64(1 + rng.Intn(2_000)) // dust
		case 1:
			return int64(10_000 + rng.Intn(100_000)) // medium: where the reported bugs lived
		default:
			return int64(1 + rng.Intn(10_000_000)) // wide
		}
	}

	var viable, missed, suboptimal int
	const iterations = 5000

	for i := 0; i < iterations; i++ {
		utxos := make([]drain.UTXO, 2+rng.Intn(11))
		for j := range utxos {
			utxos[j] = drain.UTXO{TxID: fmt.Sprintf("u%02d", j), Vout: uint32(j), AmountSats: amount()}
		}
		feeRate := int64(1 + rng.Intn(200))
		maxSats := int64(1 + rng.Intn(400_000))

		txs, err := drain.GenerateBTCTxs(drain.BTCInput{
			ChainID: 8332, To: payee, FeeRate: feeRate, UTXOs: utxos, MaxSats: maxSats,
		})
		require.NoError(t, err)

		var swept int64
		for _, tx := range txs {
			var totalIn int64
			for _, in := range tx.Inputs {
				totalIn += in.AmountSats
			}
			swept += totalIn
			require.LessOrEqual(t, totalIn, maxSats, "sweep exceeds the cap")
			require.LessOrEqual(t, tx.FeeSats, totalIn/drain.MaxBTCFeeFraction, "poller would reject this fee")
			require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
			require.Equal(t, totalIn-tx.FeeSats, tx.OutputSats)
		}

		exists, bestTotal := viableSubsets(t, utxos, maxSats, feeRate, payee)
		if !exists {
			require.Empty(t, txs, "emitted a sweep where no viable subset exists")
			continue
		}

		viable++
		switch {
		case len(txs) == 0:
			missed++
			t.Errorf("no sweep emitted though a viable subset exists: feeRate=%d cap=%d utxos=%v",
				feeRate, maxSats, utxos)
		case swept < bestTotal:
			suboptimal++
		}
	}

	require.Zero(t, missed, "%d of %d viable cases emitted nothing", missed, viable)
	// guards the sample itself: if a change made almost nothing viable, "0 misses" would be hollow
	require.Greater(t, viable, iterations/10, "too few viable cases to be meaningful")
	t.Logf("viable=%d/%d missed=%d suboptimal=%d", viable, iterations, missed, suboptimal)
}
