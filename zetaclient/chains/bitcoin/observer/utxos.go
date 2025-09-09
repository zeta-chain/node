package observer

import (
	"context"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// SelectedUTXOs is a struct containing the selected UTXOs' details.
type SelectedUTXOs struct {
	// A sublist of UTXOs selected for the outbound transaction.
	UTXOs []btcjson.ListUnspentResult

	// The total value of the selected UTXOs.
	Value float64

	// The number of consolidated UTXOs.
	ConsolidatedUTXOs uint16

	// The total value of the consolidated UTXOs.
	ConsolidatedValue float64
}

// FetchUTXOs fetches TSS-owned UTXOs from the Bitcoin node
func (ob *Observer) FetchUTXOs(ctx context.Context) error {
	// skip query if node is not ready
	if !ob.isNodeEnabled() {
		return nil
	}

	// this is useful when a zetaclient's pending nonce lagged behind for whatever reason.
	ob.refreshPendingNonce(ctx)

	// list all unspent UTXOs (160ms)
	tssAddr, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	if err != nil {
		return err
	}
	utxos, err := ob.rpc.ListUnspentMinMaxAddresses(ctx, 0, 9999999, []btcutil.Address{tssAddr})
	if err != nil {
		return err
	}

	// rigid sort to make utxo list deterministic
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Amount == utxos[j].Amount {
			if utxos[i].TxID == utxos[j].TxID {
				return utxos[i].Vout < utxos[j].Vout
			}
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Amount < utxos[j].Amount
	})

	// filter UTXOs good to spend for next TSS transaction
	utxosFiltered := make([]btcjson.ListUnspentResult, 0)
	for _, utxo := range utxos {
		// UTXOs big enough to cover the cost of spending themselves
		if utxo.Amount < common.DefaultDepositorFee {
			continue
		}
		// we don't want to spend other people's unconfirmed UTXOs as they may not be safe to spend
		if utxo.Confirmations == 0 {
			if !ob.IsTSSTransaction(utxo.TxID) {
				continue
			}
		}
		utxosFiltered = append(utxosFiltered, utxo)
	}

	ob.TelemetryServer().SetNumberOfUTXOs(ob.Chain(), len(utxosFiltered))
	ob.Mu().Lock()
	ob.utxos = utxosFiltered
	ob.Mu().Unlock()
	return nil
}

// SelectUTXOs selects a sublist of utxos to be used as inputs.
//
// Parameters:
//   - amount: The desired minimum total value of the selected UTXOs.
//   - utxos2Spend: The maximum number of UTXOs to spend.
//   - nonce: The nonce of the outbound transaction.
//   - consolidateRank: The rank below which UTXOs will be consolidated.
//   - test: true for unit test only.
//
// Returns: a sublist (includes previous nonce-mark) of UTXOs or an error if the qualifying sublist cannot be found.
func (ob *Observer) SelectUTXOs(
	ctx context.Context,
	amount float64,
	utxosToSpend uint16,
	nonce uint64,
	consolidateRank uint16,
) (SelectedUTXOs, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.Mu().Lock()
		defer ob.Mu().Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutboundHashByNonce(ctx, nonce-1)
		if err != nil {
			return SelectedUTXOs{}, err
		}
		ob.Mu().Lock()
		defer ob.Mu().Unlock()
		idx, err = ob.findNonceMarkUTXO(nonce-1, preTxid)
		if err != nil {
			return SelectedUTXOs{}, err
		}
	}

	// select smallest possible UTXOs to make payment
	total := 0.0
	left, right := 0, 0
	for total < amount && right < len(ob.utxos) {
		if utxosToSpend > 0 { // expand sublist
			total += ob.utxos[right].Amount
			right++
			utxosToSpend--
		} else { // pop the smallest utxo and append the current one
			total -= ob.utxos[left].Amount
			total += ob.utxos[right].Amount
			left++
			right++
		}
	}
	results := make([]btcjson.ListUnspentResult, right-left)
	copy(results, ob.utxos[left:right])

	// include nonce-mark as the 1st input
	if idx >= 0 { // for nonce > 0
		if idx < left || idx >= right {
			total += ob.utxos[idx].Amount
			results = append([]btcjson.ListUnspentResult{ob.utxos[idx]}, results...)
		} else { // move nonce-mark to left
			for i := idx - left; i > 0; i-- {
				results[i], results[i-1] = results[i-1], results[i]
			}
		}
	}
	if total < amount {
		return SelectedUTXOs{}, fmt.Errorf(
			"SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v",
			total,
			amount,
		)
	}

	// consolidate biggest possible UTXOs to maximize consolidated value
	// consolidation happens only when there are more than (or equal to) consolidateRank (10) UTXOs
	utxoRank, consolidatedUtxo, consolidatedValue := uint16(0), uint16(0), 0.0
	for i := len(ob.utxos) - 1; i >= 0 && utxosToSpend > 0; i-- { // iterate over UTXOs big-to-small
		if i != idx && (i < left || i >= right) { // exclude nonce-mark and already selected UTXOs
			utxoRank++
			if utxoRank >= consolidateRank { // consolication starts from the 10-ranked UTXO based on value
				utxosToSpend--
				consolidatedUtxo++
				total += ob.utxos[i].Amount
				consolidatedValue += ob.utxos[i].Amount
				results = append(results, ob.utxos[i])
			}
		}
	}

	return SelectedUTXOs{
		UTXOs:             results,
		Value:             total,
		ConsolidatedUTXOs: consolidatedUtxo,
		ConsolidatedValue: consolidatedValue,
	}, nil
}

// findNonceMarkUTXO finds the nonce-mark UTXO in the list of UTXOs.
func (ob *Observer) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	logger := ob.logger.Outbound.With().Str(logs.FieldMethod, "findNonceMarkUTXO").Logger()

	tssAddress := ob.TSSAddressString()
	amount := chains.NonceMarkAmount(nonce)

	for i, utxo := range ob.utxos {
		sats, err := common.GetSatoshis(utxo.Amount)
		if err != nil {
			logger.Error().
				Err(err).
				Any("utxo", utxo).
				Msg("error getting satoshis for utxo")
			return i, err
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid && utxo.Vout == 0 {
			logger.Info().
				Str(logs.FieldBtcTxid, utxo.TxID).
				Int64("amount", sats).
				Msg("found nonce-mark utxo")
			return i, nil
		}
	}

	return -1, fmt.Errorf("FindNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}
