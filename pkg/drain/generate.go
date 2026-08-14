// Package drain builds the signed drain payload from resolved per-chain inputs.
// The logic is network-agnostic: callers (the zetatool subcommand and the e2e test)
// fetch balances, gas, nonces and UTXOs however they like, then hand them here so the
// partitioning, fee math and signing live in one deterministic place.
package drain

import (
	"crypto/ecdsa"
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/draintx"
	"github.com/zeta-chain/node/pkg/migration"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

// MaxBTCFeeFraction bounds a sweep's fee to at most 1/10 of its input total. It is the single
// source of truth shared with the poller's validateBTCFee, so the generator never emits a group
// the poller would reject.
const MaxBTCFeeFraction = 10

// MaxEVMGasPriceGwei caps the pinned drain gas price so a leaked operator key can't burn EVM funds
// via an absurd gas price.
const MaxEVMGasPriceGwei int64 = 10_000

// EVMInput is the resolved state needed to build a single EVM drain tx.
type EVMInput struct {
	ChainID        int64
	To             string // hardcoded safe receiver
	Balance        sdkmath.Uint
	MedianGasPrice sdkmath.Uint
	Nonce          uint64
	// MaxAmount caps the transferred amount (wei) so operators can rehearse the drain with a
	// small value before committing the full balance. Nil or zero means no cap (full drain).
	// The fee still comes out of the balance; only the transfer is capped, so the remainder
	// stays at the TSS address.
	MaxAmount sdkmath.Uint
}

// UTXO is a single unspent output pinned into a BTC sweep.
type UTXO struct {
	TxID       string
	Vout       uint32
	AmountSats int64
}

// BTCInput is the resolved state needed to build the BTC drain sweeps.
type BTCInput struct {
	ChainID int64
	To      btcutil.Address // hardcoded safe receiver
	FeeRate int64           // sat/vB
	UTXOs   []UTXO
	// MaxSats caps the total input value of the sweep so operators can rehearse the drain with
	// a small value. Zero means no cap (sweep everything). A UTXO cannot be spent partially and
	// the sweep has no change output, so the cap is applied by selecting a subset of UTXOs
	// (largest first, see selectCappedUTXOs) whose total stays within it — never by reducing
	// an output.
	MaxSats int64
}

// GenerateEVMTx builds a fully-resolved EVM drain tx from a resolved input.
func GenerateEVMTx(in EVMInput) (draintx.EVMTx, error) {
	amount, gasPrice, gasLimit, err := migration.ComputeEVMMigrationWithGasLimit(
		in.Balance,
		in.MedianGasPrice,
		migration.DrainEVMGasLimit,
	)
	if err != nil {
		return draintx.EVMTx{}, err
	}
	if gasPrice.IsZero() || gasLimit == 0 {
		return draintx.EVMTx{}, fmt.Errorf(
			"refusing to pin unbroadcastable evm tx: gas price %s, gas limit %d",
			gasPrice, gasLimit,
		)
	}
	// The cap only ever lowers the transfer; the fee is unchanged, so the untransferred
	// remainder simply stays at the TSS address for the real drain.
	if !in.MaxAmount.IsNil() && !in.MaxAmount.IsZero() && in.MaxAmount.LT(amount) {
		amount = in.MaxAmount
	}
	// A balance that is entirely fee (or, with a cap, a rounding artifact) would pin a
	// zero-value transfer that pays gas to move nothing.
	if amount.IsZero() {
		return draintx.EVMTx{}, fmt.Errorf("refusing to pin a zero-amount evm tx")
	}
	return draintx.EVMTx{
		ChainID:  in.ChainID,
		To:       in.To,
		Nonce:    in.Nonce,
		Amount:   amount.String(),
		GasPrice: gasPrice.String(),
		GasLimit: gasLimit,
	}, nil
}

// GenerateBTCTxs partitions the UTXO set into disjoint groups of at most
// btccommon.MaxNoOfInputsPerTx and builds one independent sweep per group. Each sweep spends
// its pinned inputs into a single output to the receiver, minus the miner fee.
//
// UTXOs are sorted by amount descending (tie-broken by TxID then Vout) so the largest inputs
// pack into economical txs and dust clusters into the trailing groups. Groups that cannot
// cover the fee, or whose output would fall below dust, are skipped rather than aborting the
// whole sweep — a returned empty slice means there is no economical BTC to sweep. The total
// ordering makes partitioning deterministic so every node builds byte-identical txs.
//
// When in.MaxSats is set the sweep is first narrowed to a subset of UTXOs within that cap
// (see selectCappedUTXOs); everything downstream is unchanged, so a capped run exercises the
// exact same partitioning, fee math and signing path as the real drain.
func GenerateBTCTxs(in BTCInput) ([]draintx.BTCTx, error) {
	if len(in.UTXOs) == 0 {
		return nil, nil
	}

	// copy before sorting so the caller's slice is untouched
	utxos := make([]UTXO, len(in.UTXOs))
	copy(utxos, in.UTXOs)
	sortUTXOsForSweep(utxos)
	if in.MaxSats > 0 {
		capped, err := selectCappedUTXOs(utxos, in.MaxSats, in.FeeRate, in.To)
		if err != nil {
			return nil, err
		}
		if len(capped) == 0 {
			return nil, nil
		}
		utxos = capped
	}

	to := in.To.EncodeAddress()
	txs := make([]draintx.BTCTx, 0, (len(utxos)+btccommon.MaxNoOfInputsPerTx-1)/btccommon.MaxNoOfInputsPerTx)

	for start := 0; start < len(utxos); start += btccommon.MaxNoOfInputsPerTx {
		end := start + btccommon.MaxNoOfInputsPerTx
		if end > len(utxos) {
			end = len(utxos)
		}
		group := utxos[start:end]

		var totalSats int64
		inputs := make([]draintx.BTCInput, len(group))
		for i, u := range group {
			totalSats += u.AmountSats
			inputs[i] = draintx.BTCInput{TxID: u.TxID, Vout: u.Vout, AmountSats: u.AmountSats}
		}

		outputSats, feeSats, err := migration.ComputeBTCSweep(totalSats, in.FeeRate, len(group), in.To)
		// A group that cannot cover the fee, whose output is below dust, or whose fee exceeds the
		// poller's cap is uneconomical (or would be rejected by validateBTCFee); skip it rather
		// than aborting the whole drain. The <=20-input partition keeps each emitted sweep within
		// OutboundBytesMax and, with this cap, inside the poller's fee bound.
		if err != nil || outputSats < constant.BTCWithdrawalDustAmount || feeSats > totalSats/MaxBTCFeeFraction {
			continue
		}

		txs = append(txs, draintx.BTCTx{
			ChainID:    in.ChainID,
			To:         to,
			OutputSats: outputSats,
			FeeSats:    feeSats,
			Inputs:     inputs,
		})
	}

	return txs, nil
}

// sortUTXOsForSweep orders utxos by amount descending, tie-broken by TxID then Vout. The total
// ordering is what makes both the capped selection and the group partitioning deterministic, so
// every node builds byte-identical txs.
func sortUTXOsForSweep(utxos []UTXO) {
	sort.Slice(utxos, func(i, j int) bool {
		if utxos[i].AmountSats != utxos[j].AmountSats {
			return utxos[i].AmountSats > utxos[j].AmountSats
		}
		if utxos[i].TxID != utxos[j].TxID {
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Vout < utxos[j].Vout
	})
}

// selectCappedUTXOs narrows utxos to a subset whose total value stays within maxSats, for a
// small-value rehearsal of the sweep. A UTXO is indivisible and the sweep has no change output,
// so this is the only way to bound a BTC drain's value. It expects utxos already ordered by
// sortUTXOsForSweep.
//
// Selection is therefore largest-first among the UTXOs that fit under the cap, and an input is
// taken only if the resulting group would still satisfy the poller's fee bound
// (fee <= total/MaxBTCFeeFraction). That second condition is what makes a capped run usable:
//
//   - Fitting under the cap is not enough. A UTXO too big for the cap is skipped and the walk
//     continues down the list, so without the fee test a small cap slides past every large UTXO
//     and fills the group with dust instead. The fee of 20 dust inputs dwarfs their value, the
//     whole group is then dropped as uneconomical, and the rehearsal silently covers no BTC at
//     all — the one leg with no testnet coverage.
//   - Topping up an already-viable group can destroy it. Each extra input adds a fixed ~68 vB of
//     fee, so a 300-sat input added to a viable 90k-sat sweep at 50 sat/vB raises the fee past
//     the bound and the group is dropped. An input must pay for its own marginal fee under the
//     same 1/MaxBTCFeeFraction rule the poller applies.
//
// The economics cannot be tested per input as it is added, because the fixed part of the fee is
// only amortised once the group is complete: two 15,000-sat UTXOs at 10 sat/vB are viable together
// (fee 2390 against a 3000 allowance) while neither clears the single-input threshold of 1710
// against 1500. Judging inputs one at a time rejects the pair and emits nothing. So selection runs
// in two phases:
//
//  1. Take largest-first everything that fits under the cap, up to one tx worth of inputs, on
//     value alone.
//  2. Drop the smallest input — the tail, since the list is descending — until the group satisfies
//     the poller's bound, or nothing is left.
//
// Trimming from the tail is what preserves the second property above: the marginal input that
// broke the bound is exactly the one removed first. And because phase 2 tests the completed group,
// a non-empty result is economical by construction rather than by luck.
//
// An empty result means no *prefix* of the cap-limited selection is viable. See MinViableSweepSats
// for the floor on a sweep's total, which is what the CLI reports.
func selectCappedUTXOs(utxos []UTXO, maxSats, feeRate int64, to btcutil.Address) ([]UTXO, error) {
	// phase 1: fill on value alone
	selected := make([]UTXO, 0, btccommon.MaxNoOfInputsPerTx)
	var total int64
	for _, u := range utxos {
		if len(selected) == btccommon.MaxNoOfInputsPerTx {
			break
		}
		if u.AmountSats <= 0 || total+u.AmountSats > maxSats {
			continue
		}
		total += u.AmountSats
		selected = append(selected, u)
	}

	// phase 2: trim the tail until the completed group clears the poller's bound
	for len(selected) > 0 {
		size, err := btccommon.EstimateOutboundSize(int64(len(selected)), []btcutil.Address{to})
		if err != nil {
			return nil, err
		}
		if size*feeRate <= total/MaxBTCFeeFraction {
			break
		}
		total -= selected[len(selected)-1].AmountSats
		selected = selected[:len(selected)-1]
	}
	return selected, nil
}

// MinViableSweepSats is the floor on a sweep's input *total* at this fee rate: enough to keep the
// fee within 1/MaxBTCFeeFraction of the total and to leave an output above dust. A --btc-max-sats
// below it can only ever produce an empty BTC section, so the CLI reports it rather than leaving
// the operator to infer it from a missing tx.
//
// It is priced for one input because the fee grows with the input count while the bound scales with
// the total, so a single input is the cheapest shape a viable sweep can take — which makes this a
// floor on the total for *any* input count, not a per-UTXO requirement. A group of UTXOs each
// individually below it can still clear it together, so never read this as "no single UTXO reaches
// it, therefore no cap can work".
func MinViableSweepSats(feeRate int64, to btcutil.Address) (int64, error) {
	size, err := btccommon.EstimateOutboundSize(1, []btcutil.Address{to})
	if err != nil {
		return 0, err
	}
	fee := size * feeRate
	minByFee := fee * MaxBTCFeeFraction
	minByDust := fee + constant.BTCWithdrawalDustAmount
	if minByFee > minByDust {
		return minByFee, nil
	}
	return minByDust, nil
}

// BuildPayload assembles the txs into a payload and signs it with priv. network is baked into the
// signed bytes so a payload can never be replayed against a client armed for a different network.
func BuildPayload(
	triggerHeight int64,
	seq uint64,
	final bool,
	network string,
	evmTxs []draintx.EVMTx,
	btcTxs []draintx.BTCTx,
	priv *ecdsa.PrivateKey,
) (draintx.Payload, error) {
	p := draintx.Payload{
		TriggerZetaHeight: triggerHeight,
		Seq:               seq,
		Final:             final,
		Network:           network,
		EVMTxs:            evmTxs,
		BTCTxs:            btcTxs,
	}
	if err := p.Sign(priv); err != nil {
		return draintx.Payload{}, err
	}
	return p, nil
}
