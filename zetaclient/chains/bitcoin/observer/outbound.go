package observer

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/types"
)

// GetTxID returns a unique id for outbound tx
func (ob *Observer) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.BTCAddress()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

// WatchOutTx watches Bitcoin chain for outgoing txs status
func (ob *Observer) WatchOutTx() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchOutTx", ob.GetChainParams().OutTxTicker)
	if err != nil {
		ob.logger.OutTx.Error().Err(err).Msg("error creating ticker ")
		return
	}
	defer ticker.Stop()

	ob.logger.OutTx.Info().Msgf("WatchInTx started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.OutTx.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !context.IsOutboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutTx: outbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			trackers, err := ob.zetacoreClient.GetAllOutTxTrackerByChain(ob.chain.ChainId, interfaces.Ascending)
			if err != nil {
				ob.logger.OutTx.Error().Err(err).Msgf("WatchOutTx: error GetAllOutTxTrackerByChain for chain %d", ob.chain.ChainId)
				continue
			}
			for _, tracker := range trackers {
				// get original cctx parameters
				outTxID := ob.GetTxID(tracker.Nonce)
				cctx, err := ob.zetacoreClient.GetCctxByNonce(ob.chain.ChainId, tracker.Nonce)
				if err != nil {
					ob.logger.OutTx.Info().Err(err).Msgf("WatchOutTx: can't find cctx for chain %d nonce %d", ob.chain.ChainId, tracker.Nonce)
					break
				}

				nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
				if tracker.Nonce != nonce { // Tanmay: it doesn't hurt to check
					ob.logger.OutTx.Error().Msgf("WatchOutTx: tracker nonce %d not match cctx nonce %d", tracker.Nonce, nonce)
					break
				}

				if len(tracker.HashList) > 1 {
					ob.logger.OutTx.Warn().Msgf("WatchOutTx: oops, outTxID %s got multiple (%d) outTx hashes", outTxID, len(tracker.HashList))
				}

				// iterate over all txHashes to find the truly included one.
				// we do it this (inefficient) way because we don't rely on the first one as it may be a false positive (for unknown reason).
				txCount := 0
				var txResult *btcjson.GetTransactionResult
				for _, txHash := range tracker.HashList {
					result, inMempool := ob.checkIncludedTx(cctx, txHash.TxHash)
					if result != nil && !inMempool { // included
						txCount++
						txResult = result
						ob.logger.OutTx.Info().Msgf("WatchOutTx: included outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, tracker.Nonce)
						if txCount > 1 {
							ob.logger.OutTx.Error().Msgf(
								"WatchOutTx: checkIncludedTx passed, txCount %d chain %d nonce %d result %v", txCount, ob.chain.ChainId, tracker.Nonce, result)
						}
					}
				}

				if txCount == 1 { // should be only one txHash included for each nonce
					ob.setIncludedTx(tracker.Nonce, txResult)
				} else if txCount > 1 {
					ob.removeIncludedTx(tracker.Nonce) // we can't tell which txHash is true, so we remove all (if any) to be safe
					ob.logger.OutTx.Error().Msgf("WatchOutTx: included multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, tracker.Nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutTxTicker, ob.logger.OutTx)
		case <-ob.stop:
			ob.logger.OutTx.Info().Msgf("WatchOutTx stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// IsOutboundProcessed returns isIncluded(or inMempool), isConfirmed, Error
func (ob *Observer) IsOutboundProcessed(cctx *crosschaintypes.CrossChainTx, logger zerolog.Logger) (bool, bool, error) {
	params := *cctx.GetCurrentOutTxParam()
	sendHash := cctx.Index
	nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce

	// get broadcasted outtx and tx result
	outTxID := ob.GetTxID(nonce)
	logger.Info().Msgf("IsOutboundProcessed %s", outTxID)

	ob.Mu.Lock()
	txnHash, broadcasted := ob.broadcastedTx[outTxID]
	res, included := ob.includedTxResults[outTxID]
	ob.Mu.Unlock()

	if !included {
		if !broadcasted {
			return false, false, nil
		}
		// If the broadcasted outTx is nonce 0, just wait for inclusion and don't schedule more keysign
		// Schedule more than one keysign for nonce 0 can lead to duplicate payments.
		// One purpose of nonce mark UTXO is to avoid duplicate payment based on the fact that Bitcoin
		// prevents double spending of same UTXO. However, for nonce 0, we don't have a prior nonce (e.g., -1)
		// for the signer to check against when making the payment. Signer treats nonce 0 as a special case in downstream code.
		if nonce == 0 {
			return true, false, nil
		}

		// Try including this outTx broadcasted by myself
		txResult, inMempool := ob.checkIncludedTx(cctx, txnHash)
		if txResult == nil { // check failed, try again next time
			return false, false, nil
		} else if inMempool { // still in mempool (should avoid unnecessary Tss keysign)
			ob.logger.OutTx.Info().Msgf("IsOutboundProcessed: outTx %s is still in mempool", outTxID)
			return true, false, nil
		}
		// included
		ob.setIncludedTx(nonce, txResult)

		// Get tx result again in case it is just included
		res = ob.getIncludedTx(nonce)
		if res == nil {
			return false, false, nil
		}
		ob.logger.OutTx.Info().Msgf("IsOutboundProcessed: setIncludedTx succeeded for outTx %s", outTxID)
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in observeOutTx()
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	logger.Debug().Msgf("Bitcoin outTx confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())
	zetaHash, ballot, err := ob.zetacoreClient.PostVoteOutbound(
		sendHash,
		res.TxID,
		// #nosec G701 always positive
		uint64(res.BlockIndex),
		0,   // gas used not used with Bitcoin
		nil, // gas price not used with Bitcoin
		0,   // gas limit not used with Bitcoin
		amountInSat,
		chains.ReceiveStatus_success,
		ob.chain,
		nonce,
		coin.CoinType_Gas,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("IsOutboundProcessed: error confirming bitcoin outTx %s, nonce %d ballot %s", res.TxID, nonce, ballot)
	} else if zetaHash != "" {
		logger.Info().Msgf("IsOutboundProcessed: confirmed Bitcoin outTx %s, zeta tx hash %s nonce %d ballot %s", res.TxID, zetaHash, nonce, ballot)
	}

	return true, true, nil
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
// Returns:
//   - a sublist (includes previous nonce-mark) of UTXOs or an error if the qualifying sublist cannot be found.
//   - the total value of the selected UTXOs.
//   - the number of consolidated UTXOs.
//   - the total value of the consolidated UTXOs.
func (ob *Observer) SelectUTXOs(
	amount float64,
	utxosToSpend uint16,
	nonce uint64,
	consolidateRank uint16,
	test bool,
) ([]btcjson.ListUnspentResult, float64, uint16, float64, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutTxidByNonce(nonce-1, test)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		idx, err = ob.findNonceMarkUTXO(nonce-1, preTxid)
		if err != nil {
			return nil, 0, 0, 0, err
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
		return nil, 0, 0, 0, fmt.Errorf("SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v", total, amount)
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

	return results, total, consolidatedUtxo, consolidatedValue, nil
}

// refreshPendingNonce tries increasing the artificial pending nonce of outTx (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetacore.
func (ob *Observer) refreshPendingNonce() {
	// get pending nonces from zetacore
	p, err := ob.zetacoreClient.GetPendingNoncesByChain(ob.chain.ChainId)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
	}

	// increase pending nonce if lagged behind
	ob.Mu.Lock()
	pendingNonce := ob.pendingNonce
	ob.Mu.Unlock()

	// #nosec G701 always non-negative
	nonceLow := uint64(p.NonceLow)
	if nonceLow > pendingNonce {
		// get the last included outTx hash
		txid, err := ob.getOutTxidByNonce(nonceLow-1, false)
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting last outTx txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.Chain.Info().Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

func (ob *Observer) getOutTxidByNonce(nonce uint64, test bool) (string, error) {

	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetacore based on majority vote.
	if res := ob.getIncludedTx(nonce); res != nil {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetacore
		send, err := ob.zetacoreClient.GetCctxByNonce(ob.chain.ChainId, nonce)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting cctx for nonce %d", nonce)
		}
		txid := send.GetCurrentOutTxParam().OutboundTxHash
		if txid == "" {
			return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
		}
		// make sure it's a real Bitcoin txid
		_, getTxResult, err := GetTxResultByHash(ob.rpcClient, txid)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting outTx result for nonce %d hash %s", nonce, txid)
		}
		if getTxResult.Confirmations <= 0 { // just a double check
			return "", fmt.Errorf("getOutTxidByNonce: outTx txid %s for nonce %d is not included", txid, nonce)
		}
		return txid, nil
	}
	return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
}

func (ob *Observer) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	amount := chains.NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := bitcoin.GetSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.OutTx.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid && utxo.Vout == 0 {
			ob.logger.OutTx.Info().Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %d satoshi", utxo.TxID, sats)
			return i, nil
		}
	}
	return -1, fmt.Errorf("findNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}

// checkIncludedTx checks if a txHash is included and returns (txResult, inMempool)
// Note: if txResult is nil, then inMempool flag should be ignored.
func (ob *Observer) checkIncludedTx(cctx *crosschaintypes.CrossChainTx, txHash string) (*btcjson.GetTransactionResult, bool) {
	outTxID := ob.GetTxID(cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
	hash, getTxResult, err := GetTxResultByHash(ob.rpcClient, txHash)
	if err != nil {
		ob.logger.OutTx.Error().Err(err).Msgf("checkIncludedTx: error GetTxResultByHash: %s", txHash)
		return nil, false
	}

	if txHash != getTxResult.TxID { // just in case, we'll use getTxResult.TxID later
		ob.logger.OutTx.Error().Msgf("checkIncludedTx: inconsistent txHash %s and getTxResult.TxID %s", txHash, getTxResult.TxID)
		return nil, false
	}

	if getTxResult.Confirmations >= 0 { // check included tx only
		err = ob.checkTssOutTxResult(cctx, hash, getTxResult)
		if err != nil {
			ob.logger.OutTx.Error().Err(err).Msgf("checkIncludedTx: error verify bitcoin outTx %s outTxID %s", txHash, outTxID)
			return nil, false
		}
		return getTxResult, false // included
	}
	return getTxResult, true // in mempool
}

// setIncludedTx saves included tx result in memory
func (ob *Observer) setIncludedTx(nonce uint64, getTxResult *btcjson.GetTransactionResult) {
	txHash := getTxResult.TxID
	outTxID := ob.GetTxID(nonce)

	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	res, found := ob.includedTxResults[outTxID]

	if !found { // not found.
		ob.includedTxHashes[txHash] = true
		ob.includedTxResults[outTxID] = getTxResult // include new outTx and enforce rigid 1-to-1 mapping: nonce <===> txHash
		if nonce >= ob.pendingNonce {               // try increasing pending nonce on every newly included outTx
			ob.pendingNonce = nonce + 1
		}
		ob.logger.OutTx.Info().Msgf("setIncludedTx: included new bitcoin outTx %s outTxID %s pending nonce %d", txHash, outTxID, ob.pendingNonce)
	} else if txHash == res.TxID { // found same hash.
		ob.includedTxResults[outTxID] = getTxResult // update tx result as confirmations may increase
		if getTxResult.Confirmations > res.Confirmations {
			ob.logger.OutTx.Info().Msgf("setIncludedTx: bitcoin outTx %s got confirmations %d", txHash, getTxResult.Confirmations)
		}
	} else { // found other hash.
		// be alert for duplicate payment!!! As we got a new hash paying same cctx (for whatever reason).
		delete(ob.includedTxResults, outTxID) // we can't tell which txHash is true, so we remove all to be safe
		ob.logger.OutTx.Error().Msgf("setIncludedTx: duplicate payment by bitcoin outTx %s outTxID %s, prior outTx %s", txHash, outTxID, res.TxID)
	}
}

// getIncludedTx gets the receipt and transaction from memory
func (ob *Observer) getIncludedTx(nonce uint64) *btcjson.GetTransactionResult {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.includedTxResults[ob.GetTxID(nonce)]
}

// removeIncludedTx removes included tx from memory
func (ob *Observer) removeIncludedTx(nonce uint64) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	txResult, found := ob.includedTxResults[ob.GetTxID(nonce)]
	if found {
		delete(ob.includedTxHashes, txResult.TxID)
		delete(ob.includedTxResults, ob.GetTxID(nonce))
	}
}

// Basic TSS outTX checks:
//   - should be able to query the raw tx
//   - check if all inputs are segwit && TSS inputs
//
// Returns: true if outTx passes basic checks.
func (ob *Observer) checkTssOutTxResult(cctx *crosschaintypes.CrossChainTx, hash *chainhash.Hash, res *btcjson.GetTransactionResult) error {
	params := cctx.GetCurrentOutTxParam()
	nonce := params.OutboundTxTssNonce
	rawResult, err := GetRawTxResult(ob.rpcClient, hash, res)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: error GetRawTxResultByHash %s", hash.String())
	}
	err = ob.checkTSSVin(rawResult.Vin, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vin in outTx %s nonce %d", hash, nonce)
	}

	// differentiate between normal and restricted cctx
	if compliance.IsCctxRestricted(cctx) {
		err = ob.checkTSSVoutCancelled(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vout in cancelled outTx %s nonce %d", hash, nonce)
		}
	} else {
		err = ob.checkTSSVout(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vout in outTx %s nonce %d", hash, nonce)
		}
	}
	return nil
}

// checkTSSVin checks vin is valid if:
//   - The first input is the nonce-mark
//   - All inputs are from TSS address
func (ob *Observer) checkTSSVin(vins []btcjson.Vin, nonce uint64) error {
	// vins: [nonce-mark, UTXO1, UTXO2, ...]
	if nonce > 0 && len(vins) <= 1 {
		return fmt.Errorf("checkTSSVin: len(vins) <= 1")
	}
	pubKeyTss := hex.EncodeToString(ob.Tss.PubKeyCompressedBytes())
	for i, vin := range vins {
		// The length of the Witness should be always 2 for SegWit inputs.
		if len(vin.Witness) != 2 {
			return fmt.Errorf("checkTSSVin: expected 2 witness items, got %d", len(vin.Witness))
		}
		if vin.Witness[1] != pubKeyTss {
			return fmt.Errorf("checkTSSVin: witness pubkey %s not match TSS pubkey %s", vin.Witness[1], pubKeyTss)
		}
		// 1st vin: nonce-mark MUST come from prior TSS outTx
		if nonce > 0 && i == 0 {
			preTxid, err := ob.getOutTxidByNonce(nonce-1, false)
			if err != nil {
				return fmt.Errorf("checkTSSVin: error findTxIDByNonce %d", nonce-1)
			}
			// nonce-mark MUST the 1st output that comes from prior TSS outTx
			if vin.Txid != preTxid || vin.Vout != 0 {
				return fmt.Errorf("checkTSSVin: invalid nonce-mark txid %s vout %d, expected txid %s vout 0", vin.Txid, vin.Vout, preTxid)
			}
		}
	}
	return nil
}

// checkTSSVout vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the correct payment to recipient
//   - The third output is the change to TSS (optional)
func (ob *Observer) checkTSSVout(params *crosschaintypes.OutboundTxParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, payment to recipient, change to TSS (optional)]
	if !(len(vouts) == 2 || len(vouts) == 3) {
		return fmt.Errorf("checkTSSVout: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.OutboundTxTssNonce
	tssAddress := ob.Tss.BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverExpected := tssAddress
		if vout.N == 1 {
			// the 2nd output is the payment to recipient
			receiverExpected = params.Receiver
		}
		receiverVout, amount, err := bitcoin.DecodeTSSVout(vout, receiverExpected, ob.chain)
		if err != nil {
			return err
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVout: nonce-mark address %s not match TSS address %s", receiverVout, tssAddress)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVout: nonce-mark amount %d not match nonce-mark amount %d", amount, chains.NonceMarkAmount(nonce))
			}
		case 1: // 2nd vout: payment to recipient
			if receiverVout != params.Receiver {
				return fmt.Errorf("checkTSSVout: output address %s not match params receiver %s", receiverVout, params.Receiver)
			}
			// #nosec G701 always positive
			if uint64(amount) != params.Amount.Uint64() {
				return fmt.Errorf("checkTSSVout: output amount %d not match params amount %d", amount, params.Amount)
			}
		case 2: // 3rd vout: change to TSS (optional)
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVout: change address %s not match TSS address %s", receiverVout, tssAddress)
			}
		}
	}
	return nil
}

// checkTSSVoutCancelled vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the change to TSS (optional)
func (ob *Observer) checkTSSVoutCancelled(params *crosschaintypes.OutboundTxParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, change to TSS (optional)]
	if !(len(vouts) == 1 || len(vouts) == 2) {
		return fmt.Errorf("checkTSSVoutCancelled: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.OutboundTxTssNonce
	tssAddress := ob.Tss.BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverVout, amount, err := bitcoin.DecodeTSSVout(vout, tssAddress, ob.chain)
		if err != nil {
			return errors.Wrap(err, "checkTSSVoutCancelled: error decoding P2WPKH vout")
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVoutCancelled: nonce-mark address %s not match TSS address %s", receiverVout, tssAddress)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVoutCancelled: nonce-mark amount %d not match nonce-mark amount %d", amount, chains.NonceMarkAmount(nonce))
			}
		case 1: // 2nd vout: change to TSS (optional)
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVoutCancelled: change address %s not match TSS address %s", receiverVout, tssAddress)
			}
		}
	}
	return nil
}
