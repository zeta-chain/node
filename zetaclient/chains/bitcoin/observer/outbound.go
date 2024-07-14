package observer

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// GetTxID returns a unique id for outbound tx
func (ob *Observer) GetTxID(nonce uint64) string {
	tssAddr := ob.TSS().BTCAddress()
	return fmt.Sprintf("%d-%s-%d", ob.Chain().ChainId, tssAddr, nonce)
}

// WatchOutbound watches Bitcoin chain for outgoing txs status
// TODO(revamp): move ticker functions to a specific file
// TODO(revamp): move into a separate package
func (ob *Observer) WatchOutbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	ticker, err := types.NewDynamicTicker("Bitcoin_WatchOutbound", ob.GetChainParams().OutboundTicker)
	if err != nil {
		return errors.Wrap(err, "unable to create dynamic ticker")
	}

	defer ticker.Stop()

	chainID := ob.Chain().ChainId
	ob.logger.Outbound.Info().Msgf("WatchOutbound started for chain %d", chainID)
	sampledLogger := ob.logger.Outbound.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !app.IsOutboundObservationEnabled(ob.GetChainParams()) {
				sampledLogger.Info().
					Msgf("WatchOutbound: outbound observation is disabled for chain %d", chainID)
				continue
			}
			trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
			if err != nil {
				ob.logger.Outbound.Error().
					Err(err).
					Msgf("WatchOutbound: error GetAllOutboundTrackerByChain for chain %d", chainID)
				continue
			}
			for _, tracker := range trackers {
				// get original cctx parameters
				outboundID := ob.GetTxID(tracker.Nonce)
				cctx, err := ob.ZetacoreClient().GetCctxByNonce(ctx, chainID, tracker.Nonce)
				if err != nil {
					ob.logger.Outbound.Info().
						Err(err).
						Msgf("WatchOutbound: can't find cctx for chain %d nonce %d", chainID, tracker.Nonce)
					break
				}

				nonce := cctx.GetCurrentOutboundParam().TssNonce
				if tracker.Nonce != nonce { // Tanmay: it doesn't hurt to check
					ob.logger.Outbound.Error().
						Msgf("WatchOutbound: tracker nonce %d not match cctx nonce %d", tracker.Nonce, nonce)
					break
				}

				if len(tracker.HashList) > 1 {
					ob.logger.Outbound.Warn().
						Msgf("WatchOutbound: oops, outboundID %s got multiple (%d) outbound hashes", outboundID, len(tracker.HashList))
				}

				// iterate over all txHashes to find the truly included one.
				// we do it this (inefficient) way because we don't rely on the first one as it may be a false positive (for unknown reason).
				txCount := 0
				var txResult *btcjson.GetTransactionResult
				for _, txHash := range tracker.HashList {
					result, inMempool := ob.checkIncludedTx(ctx, cctx, txHash.TxHash)
					if result != nil && !inMempool { // included
						txCount++
						txResult = result
						ob.logger.Outbound.Info().
							Msgf("WatchOutbound: included outbound %s for chain %d nonce %d", txHash.TxHash, chainID, tracker.Nonce)
						if txCount > 1 {
							ob.logger.Outbound.Error().Msgf(
								"WatchOutbound: checkIncludedTx passed, txCount %d chain %d nonce %d result %v", txCount, chainID, tracker.Nonce, result)
						}
					}
				}

				if txCount == 1 { // should be only one txHash included for each nonce
					ob.setIncludedTx(tracker.Nonce, txResult)
				} else if txCount > 1 {
					ob.removeIncludedTx(tracker.Nonce) // we can't tell which txHash is true, so we remove all (if any) to be safe
					ob.logger.Outbound.Error().Msgf("WatchOutbound: included multiple (%d) outbound for chain %d nonce %d", txCount, chainID, tracker.Nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.logger.Outbound)
		case <-ob.StopChannel():
			ob.logger.Outbound.Info().Msgf("WatchOutbound stopped for chain %d", chainID)
			return nil
		}
	}
}

// IsOutboundProcessed returns isIncluded(or inMempool), isConfirmed, Error
// TODO(revamp): rename as it vote the outbound and doesn't only check if outbound is processed
func (ob *Observer) IsOutboundProcessed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	logger zerolog.Logger,
) (bool, bool, error) {
	const (
		// not used with Bitcoin
		outboundGasUsed  = 0
		outboundGasPrice = 0
		outboundGasLimit = 0

		gasLimit      = zetacore.PostVoteOutboundGasLimit
		gasRetryLimit = 0
	)

	params := *cctx.GetCurrentOutboundParam()
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	// get broadcasted outbound and tx result
	outboundID := ob.GetTxID(nonce)
	logger.Info().Msgf("IsOutboundProcessed %s", outboundID)

	ob.Mu().Lock()
	txnHash, broadcasted := ob.broadcastedTx[outboundID]
	res, included := ob.includedTxResults[outboundID]
	ob.Mu().Unlock()

	if !included {
		if !broadcasted {
			return false, false, nil
		}
		// If the broadcasted outbound is nonce 0, just wait for inclusion and don't schedule more keysign
		// Schedule more than one keysign for nonce 0 can lead to duplicate payments.
		// One purpose of nonce mark UTXO is to avoid duplicate payment based on the fact that Bitcoin
		// prevents double spending of same UTXO. However, for nonce 0, we don't have a prior nonce (e.g., -1)
		// for the signer to check against when making the payment. Signer treats nonce 0 as a special case in downstream code.
		if nonce == 0 {
			return true, false, nil
		}

		// Try including this outbound broadcasted by myself
		txResult, inMempool := ob.checkIncludedTx(ctx, cctx, txnHash)
		if txResult == nil { // check failed, try again next time
			return false, false, nil
		} else if inMempool { // still in mempool (should avoid unnecessary Tss keysign)
			ob.logger.Outbound.Info().Msgf("IsOutboundProcessed: outbound %s is still in mempool", outboundID)
			return true, false, nil
		}
		// included
		ob.setIncludedTx(nonce, txResult)

		// Get tx result again in case it is just included
		res = ob.getIncludedTx(nonce)
		if res == nil {
			return false, false, nil
		}
		ob.logger.Outbound.Info().Msgf("IsOutboundProcessed: setIncludedTx succeeded for outbound %s", outboundID)
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in observeOutbound()
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	// Get outbound block height
	blockHeight, err := rpc.GetBlockHeightByHash(ob.btcClient, res.BlockHash)
	if err != nil {
		return true, false, errors.Wrapf(
			err,
			"IsOutboundProcessed: error getting block height by hash %s",
			res.BlockHash,
		)
	}

	logger.Debug().Msgf("Bitcoin outbound confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())

	signer := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		signer.String(),
		cctx.Index,
		res.TxID,

		// #nosec G115 always positive
		uint64(blockHeight),

		// not used with Bitcoin
		outboundGasUsed,
		math.NewInt(outboundGasPrice),
		outboundGasLimit,

		math.NewUintFromBigInt(amountInSat),
		chains.ReceiveStatus_success,
		ob.Chain().ChainId,
		nonce,
		coin.CoinType_Gas,
	)

	zetaHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, gasRetryLimit, msg)

	logFields := map[string]any{
		"outbound.external_tx_hash": res.TxID,
		"outbound.nonce":            nonce,
		"outbound.zeta_tx_hash":     zetaHash,
		"outbound.ballot":           ballot,
	}

	if err != nil {
		logger.Error().Err(err).Fields(logFields).Msg("IsOutboundProcessed: error confirming bitcoin outbound")
	} else if zetaHash != "" {
		logger.Info().Fields(logFields).Msgf("IsOutboundProcessed: confirmed Bitcoin outbound")
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
//
// TODO(revamp): move to utxo file
func (ob *Observer) SelectUTXOs(
	ctx context.Context,
	amount float64,
	utxosToSpend uint16,
	nonce uint64,
	consolidateRank uint16,
	test bool,
) ([]btcjson.ListUnspentResult, float64, uint16, float64, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.Mu().Lock()
		defer ob.Mu().Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutboundIDByNonce(ctx, nonce-1, test)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		ob.Mu().Lock()
		defer ob.Mu().Unlock()
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
		return nil, 0, 0, 0, fmt.Errorf(
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

	return results, total, consolidatedUtxo, consolidatedValue, nil
}

// refreshPendingNonce tries increasing the artificial pending nonce of outbound (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetacore.
func (ob *Observer) refreshPendingNonce(ctx context.Context) {
	// get pending nonces from zetacore
	p, err := ob.ZetacoreClient().GetPendingNoncesByChain(ctx, ob.Chain().ChainId)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
	}

	// increase pending nonce if lagged behind
	ob.Mu().Lock()
	pendingNonce := ob.pendingNonce
	ob.Mu().Unlock()

	// #nosec G115 always non-negative
	nonceLow := uint64(p.NonceLow)
	if nonceLow > pendingNonce {
		// get the last included outbound hash
		txid, err := ob.getOutboundIDByNonce(ctx, nonceLow-1, false)
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting last outbound txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.Mu().Lock()
		defer ob.Mu().Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.Chain.Info().
			Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

// getOutboundIDByNonce gets the outbound ID from the nonce of the outbound transaction
// test is true for unit test only
func (ob *Observer) getOutboundIDByNonce(ctx context.Context, nonce uint64, test bool) (string, error) {
	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetacore based on majority vote.
	if res := ob.getIncludedTx(nonce); res != nil {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetacore
		send, err := ob.ZetacoreClient().GetCctxByNonce(ctx, ob.Chain().ChainId, nonce)
		if err != nil {
			return "", errors.Wrapf(err, "getOutboundIDByNonce: error getting cctx for nonce %d", nonce)
		}
		txid := send.GetCurrentOutboundParam().Hash
		if txid == "" {
			return "", fmt.Errorf("getOutboundIDByNonce: cannot find outbound txid for nonce %d", nonce)
		}
		// make sure it's a real Bitcoin txid
		_, getTxResult, err := rpc.GetTxResultByHash(ob.btcClient, txid)
		if err != nil {
			return "", errors.Wrapf(
				err,
				"getOutboundIDByNonce: error getting outbound result for nonce %d hash %s",
				nonce,
				txid,
			)
		}
		if getTxResult.Confirmations <= 0 { // just a double check
			return "", fmt.Errorf("getOutboundIDByNonce: outbound txid %s for nonce %d is not included", txid, nonce)
		}
		return txid, nil
	}
	return "", fmt.Errorf("getOutboundIDByNonce: cannot find outbound txid for nonce %d", nonce)
}

// findNonceMarkUTXO finds the nonce-mark UTXO in the list of UTXOs.
func (ob *Observer) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	tssAddress := ob.TSS().BTCAddressWitnessPubkeyHash().EncodeAddress()
	amount := chains.NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := bitcoin.GetSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.Outbound.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid && utxo.Vout == 0 {
			ob.logger.Outbound.Info().
				Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %d satoshi", utxo.TxID, sats)
			return i, nil
		}
	}
	return -1, fmt.Errorf("findNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}

// checkIncludedTx checks if a txHash is included and returns (txResult, inMempool)
// Note: if txResult is nil, then inMempool flag should be ignored.
func (ob *Observer) checkIncludedTx(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	txHash string,
) (*btcjson.GetTransactionResult, bool) {
	outboundID := ob.GetTxID(cctx.GetCurrentOutboundParam().TssNonce)
	hash, getTxResult, err := rpc.GetTxResultByHash(ob.btcClient, txHash)
	if err != nil {
		ob.logger.Outbound.Error().Err(err).Msgf("checkIncludedTx: error GetTxResultByHash: %s", txHash)
		return nil, false
	}

	if txHash != getTxResult.TxID { // just in case, we'll use getTxResult.TxID later
		ob.logger.Outbound.Error().
			Msgf("checkIncludedTx: inconsistent txHash %s and getTxResult.TxID %s", txHash, getTxResult.TxID)
		return nil, false
	}

	if getTxResult.Confirmations >= 0 { // check included tx only
		err = ob.checkTssOutboundResult(ctx, cctx, hash, getTxResult)
		if err != nil {
			ob.logger.Outbound.Error().
				Err(err).
				Msgf("checkIncludedTx: error verify bitcoin outbound %s outboundID %s", txHash, outboundID)
			return nil, false
		}
		return getTxResult, false // included
	}
	return getTxResult, true // in mempool
}

// setIncludedTx saves included tx result in memory
func (ob *Observer) setIncludedTx(nonce uint64, getTxResult *btcjson.GetTransactionResult) {
	txHash := getTxResult.TxID
	outboundID := ob.GetTxID(nonce)

	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	res, found := ob.includedTxResults[outboundID]

	if !found { // not found.
		ob.includedTxHashes[txHash] = true
		ob.includedTxResults[outboundID] = getTxResult // include new outbound and enforce rigid 1-to-1 mapping: nonce <===> txHash
		if nonce >= ob.pendingNonce {                  // try increasing pending nonce on every newly included outbound
			ob.pendingNonce = nonce + 1
		}
		ob.logger.Outbound.Info().
			Msgf("setIncludedTx: included new bitcoin outbound %s outboundID %s pending nonce %d", txHash, outboundID, ob.pendingNonce)
	} else if txHash == res.TxID { // found same hash.
		ob.includedTxResults[outboundID] = getTxResult // update tx result as confirmations may increase
		if getTxResult.Confirmations > res.Confirmations {
			ob.logger.Outbound.Info().Msgf("setIncludedTx: bitcoin outbound %s got confirmations %d", txHash, getTxResult.Confirmations)
		}
	} else { // found other hash.
		// be alert for duplicate payment!!! As we got a new hash paying same cctx (for whatever reason).
		delete(ob.includedTxResults, outboundID) // we can't tell which txHash is true, so we remove all to be safe
		ob.logger.Outbound.Error().Msgf("setIncludedTx: duplicate payment by bitcoin outbound %s outboundID %s, prior outbound %s", txHash, outboundID, res.TxID)
	}
}

// getIncludedTx gets the receipt and transaction from memory
func (ob *Observer) getIncludedTx(nonce uint64) *btcjson.GetTransactionResult {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.includedTxResults[ob.GetTxID(nonce)]
}

// removeIncludedTx removes included tx from memory
func (ob *Observer) removeIncludedTx(nonce uint64) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	txResult, found := ob.includedTxResults[ob.GetTxID(nonce)]
	if found {
		delete(ob.includedTxHashes, txResult.TxID)
		delete(ob.includedTxResults, ob.GetTxID(nonce))
	}
}

// Basic TSS outbound checks:
//   - should be able to query the raw tx
//   - check if all inputs are segwit && TSS inputs
//
// Returns: true if outbound passes basic checks.
func (ob *Observer) checkTssOutboundResult(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	hash *chainhash.Hash,
	res *btcjson.GetTransactionResult,
) error {
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce
	rawResult, err := rpc.GetRawTxResult(ob.btcClient, hash, res)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutboundResult: error GetRawTxResultByHash %s", hash.String())
	}
	err = ob.checkTSSVin(ctx, rawResult.Vin, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutboundResult: invalid TSS Vin in outbound %s nonce %d", hash, nonce)
	}

	// differentiate between normal and restricted cctx
	if compliance.IsCctxRestricted(cctx) {
		err = ob.checkTSSVoutCancelled(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(
				err,
				"checkTssOutboundResult: invalid TSS Vout in cancelled outbound %s nonce %d",
				hash,
				nonce,
			)
		}
	} else {
		err = ob.checkTSSVout(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutboundResult: invalid TSS Vout in outbound %s nonce %d", hash, nonce)
		}
	}
	return nil
}

// checkTSSVin checks vin is valid if:
//   - The first input is the nonce-mark
//   - All inputs are from TSS address
func (ob *Observer) checkTSSVin(ctx context.Context, vins []btcjson.Vin, nonce uint64) error {
	// vins: [nonce-mark, UTXO1, UTXO2, ...]
	if nonce > 0 && len(vins) <= 1 {
		return fmt.Errorf("checkTSSVin: len(vins) <= 1")
	}
	pubKeyTss := hex.EncodeToString(ob.TSS().PubKeyCompressedBytes())
	for i, vin := range vins {
		// The length of the Witness should be always 2 for SegWit inputs.
		if len(vin.Witness) != 2 {
			return fmt.Errorf("checkTSSVin: expected 2 witness items, got %d", len(vin.Witness))
		}
		if vin.Witness[1] != pubKeyTss {
			return fmt.Errorf("checkTSSVin: witness pubkey %s not match TSS pubkey %s", vin.Witness[1], pubKeyTss)
		}
		// 1st vin: nonce-mark MUST come from prior TSS outbound
		if nonce > 0 && i == 0 {
			preTxid, err := ob.getOutboundIDByNonce(ctx, nonce-1, false)
			if err != nil {
				return fmt.Errorf("checkTSSVin: error findTxIDByNonce %d", nonce-1)
			}
			// nonce-mark MUST the 1st output that comes from prior TSS outbound
			if vin.Txid != preTxid || vin.Vout != 0 {
				return fmt.Errorf(
					"checkTSSVin: invalid nonce-mark txid %s vout %d, expected txid %s vout 0",
					vin.Txid,
					vin.Vout,
					preTxid,
				)
			}
		}
	}
	return nil
}

// checkTSSVout vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the correct payment to recipient
//   - The third output is the change to TSS (optional)
func (ob *Observer) checkTSSVout(params *crosschaintypes.OutboundParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, payment to recipient, change to TSS (optional)]
	if !(len(vouts) == 2 || len(vouts) == 3) {
		return fmt.Errorf("checkTSSVout: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.TssNonce
	tssAddress := ob.TSS().BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverExpected := tssAddress
		if vout.N == 1 {
			// the 2nd output is the payment to recipient
			receiverExpected = params.Receiver
		}
		receiverVout, amount, err := bitcoin.DecodeTSSVout(vout, receiverExpected, ob.Chain())
		if err != nil {
			return err
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf(
					"checkTSSVout: nonce-mark address %s not match TSS address %s",
					receiverVout,
					tssAddress,
				)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf(
					"checkTSSVout: nonce-mark amount %d not match nonce-mark amount %d",
					amount,
					chains.NonceMarkAmount(nonce),
				)
			}
		case 1: // 2nd vout: payment to recipient
			if receiverVout != params.Receiver {
				return fmt.Errorf(
					"checkTSSVout: output address %s not match params receiver %s",
					receiverVout,
					params.Receiver,
				)
			}
			// #nosec G115 always positive
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
func (ob *Observer) checkTSSVoutCancelled(params *crosschaintypes.OutboundParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, change to TSS (optional)]
	if !(len(vouts) == 1 || len(vouts) == 2) {
		return fmt.Errorf("checkTSSVoutCancelled: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.TssNonce
	tssAddress := ob.TSS().BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverVout, amount, err := bitcoin.DecodeTSSVout(vout, tssAddress, ob.Chain())
		if err != nil {
			return errors.Wrap(err, "checkTSSVoutCancelled: error decoding P2WPKH vout")
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf(
					"checkTSSVoutCancelled: nonce-mark address %s not match TSS address %s",
					receiverVout,
					tssAddress,
				)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf(
					"checkTSSVoutCancelled: nonce-mark amount %d not match nonce-mark amount %d",
					amount,
					chains.NonceMarkAmount(nonce),
				)
			}
		case 1: // 2nd vout: change to TSS (optional)
			if receiverVout != tssAddress {
				return fmt.Errorf(
					"checkTSSVoutCancelled: change address %s not match TSS address %s",
					receiverVout,
					tssAddress,
				)
			}
		}
	}
	return nil
}
