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

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// minTxConfirmations is the minimum confirmations for a Bitcoin tx to be considered valid by the observer
	// Note: please change this value to 1 to be able to run the Bitcoin E2E RBF test
	minTxConfirmations = 0
)

// WatchOutbound watches Bitcoin chain for outgoing txs status
// TODO(revamp): move ticker functions to a specific file
// TODO(revamp): move into a separate package
func (ob *Observer) WatchOutbound(ctx context.Context) error {
	// get app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	// create outbound ticker
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchOutbound", ob.ChainParams().OutboundTicker)
	if err != nil {
		return errors.Wrap(err, "unable to create dynamic ticker")
	}
	defer ticker.Stop()

	ob.logger.Outbound.Info().Msg("WatchOutbound: started")
	sampledLogger := ob.logger.Outbound.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !app.IsOutboundObservationEnabled() {
				sampledLogger.Info().Msg("WatchOutbound: outbound observation is disabled")
				continue
			}

			// process outbound trackers
			err := ob.ProcessOutboundTrackers(ctx)
			if err != nil {
				ob.Logger().Outbound.Error().Err(err).Msg("WatchOutbound: ProcessOutboundTrackers failed")
			}

			ticker.UpdateInterval(ob.ChainParams().OutboundTicker, ob.logger.Outbound)
		case <-ob.StopChannel():
			ob.logger.Outbound.Info().Msg("WatchOutbound: stopped")
			return nil
		}
	}
}

// ProcessOutboundTrackers processes outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain failed")
	}

	// logger fields
	lf := map[string]any{
		logs.FieldMethod: "ProcessOutboundTrackers",
	}

	// process outbound trackers
	for _, tracker := range trackers {
		// set logger fields
		lf[logs.FieldNonce] = tracker.Nonce

		// get the CCTX
		cctx, err := ob.ZetacoreClient().GetCctxByNonce(ctx, chainID, tracker.Nonce)
		if err != nil {
			ob.logger.Outbound.Err(err).Fields(lf).Msg("cannot find cctx")
			break
		}
		if len(tracker.HashList) > 1 {
			ob.logger.Outbound.Warn().Msgf("oops, got multiple (%d) outbound hashes", len(tracker.HashList))
		}

		// Iterate over all txHashes to find the truly included outbound.
		// At any time, there is guarantee that only one single txHash will be considered valid and included for each nonce.
		// The reasons are:
		//   1. CCTX with nonce 'N = 0' is the past and well-controlled.
		//   2. Given any CCTX with nonce 'N > 0', its outbound MUST spend the previous nonce-mark UTXO (nonce N-1) to be considered valid.
		//   3. Bitcoin prevents double spending of the same UTXO except for RBF.
		//   4. When RBF happens, the original tx will be removed from Bitcoin core, and only the new tx will be valid.
		for _, txHash := range tracker.HashList {
			_, included := ob.TryIncludeOutbound(ctx, cctx, txHash.TxHash)
			if included {
				break
			}
		}
	}

	return nil
}

// TryIncludeOutbound tries to include an outbound for the given cctx and txHash.
//
// Due to 10-min block time, zetaclient observes outbounds both in mempool and in blocks.
// An outbound is considered included if it satisfies one of the following two cases:
//  1. a valid tx pending in mempool with confirmation == 0
//  2. a valid tx included in a block with confirmation > 0
//
// Returns: (txResult, included)
//
// Note: A 'included' tx may still be considered stuck if it sits in the mempool for too long.
func (ob *Observer) TryIncludeOutbound(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	txHash string,
) (*btcjson.GetTransactionResult, bool) {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	// check tx inclusion and save tx result
	txResult, included := ob.checkTxInclusion(ctx, cctx, txHash)
	if included {
		ob.SetIncludedTx(nonce, txResult)
	}

	return txResult, included
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *crosschaintypes.CrossChainTx) (bool, error) {
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
	outboundID := ob.OutboundID(nonce)
	ob.Logger().Outbound.Info().Msgf("VoteOutboundIfConfirmed %s", outboundID)

	ob.Mu().Lock()
	txnHash, broadcasted := ob.broadcastedTx[outboundID]
	res, included := ob.includedTxResults[outboundID]
	ob.Mu().Unlock()

	// Short-circuit in following two cases:
	//   1. Outbound neither broadcasted nor included. It requires a keysign.
	//   2. Outbound was broadcasted for nonce 0. It's an edge case (happened before) to avoid duplicate payments.
	if !included {
		if !broadcasted {
			return true, nil
		}
		// If the broadcasted outbound is nonce 0, just wait for inclusion and don't schedule more keysign
		// Schedule more than one keysign for nonce 0 can lead to duplicate payments.
		// One purpose of nonce mark UTXO is to avoid duplicate payment based on the fact that Bitcoin
		// prevents double spending of same UTXO. However, for nonce 0, we don't have a prior nonce (e.g., -1)
		// for the signer to check against when making the payment. Signer treats nonce 0 as a special case in downstream code.
		if nonce == 0 {
			ob.logger.Outbound.Info().Msgf("VoteOutboundIfConfirmed: outbound %s is nonce 0", outboundID)
			return false, nil
		}

		// Try including this outbound broadcasted by myself to supplement outbound trackers.
		// Note: each Bitcoin outbound usually gets included right after broadcasting.
		res, included = ob.TryIncludeOutbound(ctx, cctx, txnHash)
		if !included {
			return true, nil
		}
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in checkTxInclusion().
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		ob.logger.Outbound.Debug().
			Int64("currentConfirmations", res.Confirmations).
			Int64("requiredConfirmations", ob.ConfirmationsThreshold(amountInSat)).
			Msg("VoteOutboundIfConfirmed: outbound not confirmed yet")

		return false, nil
	}

	// Get outbound block height
	blockHeight, err := rpc.GetBlockHeightByHash(ob.btcClient, res.BlockHash)
	if err != nil {
		return false, errors.Wrapf(
			err,
			"VoteOutboundIfConfirmed: error getting block height by hash %s",
			res.BlockHash,
		)
	}

	ob.Logger().
		Outbound.Debug().
		Msgf("Bitcoin outbound confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())

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
		ob.Logger().
			Outbound.Error().
			Err(err).
			Fields(logFields).
			Msg("VoteOutboundIfConfirmed: error confirming bitcoin outbound")
	} else if zetaHash != "" {
		ob.Logger().Outbound.Info().Fields(logFields).Msgf("VoteOutboundIfConfirmed: confirmed Bitcoin outbound")
	}

	return false, nil
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
	// #nosec G115 always non-negative
	nonceLow := uint64(p.NonceLow)
	if nonceLow > ob.GetPendingNonce() {
		// get the last included outbound hash
		txid, err := ob.getOutboundHashByNonce(ctx, nonceLow-1, false)
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting last outbound txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.SetPendingNonce(nonceLow)
		ob.logger.Chain.Info().
			Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", nonceLow, txid)
	}
}

// getOutboundHashByNonce gets the outbound hash for given nonce.
// test is true for unit test only
func (ob *Observer) getOutboundHashByNonce(ctx context.Context, nonce uint64, test bool) (string, error) {
	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetacore based on majority vote.
	if res := ob.GetIncludedTx(nonce); res != nil {
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

// checkTxInclusion checks if a txHash is included and returns (txResult, included)
func (ob *Observer) checkTxInclusion(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	txHash string,
) (*btcjson.GetTransactionResult, bool) {
	// logger fields
	lf := map[string]any{
		logs.FieldMethod: "checkTxInclusion",
		logs.FieldNonce:  cctx.GetCurrentOutboundParam().TssNonce,
		logs.FieldTx:     txHash,
	}

	// fetch tx result
	hash, txResult, err := rpc.GetTxResultByHash(ob.btcClient, txHash)
	if err != nil {
		ob.logger.Outbound.Warn().Err(err).Fields(lf).Msg("GetTxResultByHash failed")
		return nil, false
	}

	// check minimum confirmations
	if txResult.Confirmations < minTxConfirmations {
		ob.logger.Outbound.Warn().Fields(lf).Msgf("invalid confirmations %d", txResult.Confirmations)
		return nil, false
	}

	// validate tx result
	err = ob.checkTssOutboundResult(ctx, cctx, hash, txResult)
	if err != nil {
		ob.logger.Outbound.Error().Err(err).Fields(lf).Msg("checkTssOutboundResult failed")
		return nil, false
	}

	// tx is valid and included
	return txResult, true
}

// SetIncludedTx saves included tx result in memory.
//   - the outbounds are chained (by nonce) txs sequentially included.
//   - tx results may still be set in arbitrary order as the method is called across goroutines, and it doesn't matter.
func (ob *Observer) SetIncludedTx(nonce uint64, getTxResult *btcjson.GetTransactionResult) {
	var (
		txHash     = getTxResult.TxID
		outboundID = ob.OutboundID(nonce)
		lf         = map[string]any{
			logs.FieldMethod:     "setIncludedTx",
			logs.FieldNonce:      nonce,
			logs.FieldTx:         txHash,
			logs.FieldOutboundID: outboundID,
		}
	)

	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	res, found := ob.includedTxResults[outboundID]
	if !found {
		// for new hash:
		//   - include new outbound and enforce rigid 1-to-1 mapping: nonce <===> txHash
		//   - try increasing pending nonce on every newly included outbound
		ob.tssOutboundHashes[txHash] = true
		ob.includedTxResults[outboundID] = getTxResult
		if nonce >= ob.pendingNonce {
			ob.pendingNonce = nonce + 1
		}
		ob.logger.Outbound.Info().Fields(lf).Msgf("included new bitcoin outbound, pending nonce %d", ob.pendingNonce)
	} else if txHash == res.TxID {
		// for existing hash:
		//   - update tx result because confirmations may increase
		ob.includedTxResults[outboundID] = getTxResult
		if getTxResult.Confirmations > res.Confirmations {
			ob.logger.Outbound.Info().Msgf("bitcoin outbound got %d confirmations", getTxResult.Confirmations)
		}
	} else {
		// got multiple hashes for same nonce. RBF happened.
		ob.logger.Outbound.Info().Fields(lf).Msgf("replaced bitcoin outbound %s", res.TxID)

		// remove prior txHash and txResult
		delete(ob.tssOutboundHashes, res.TxID)
		delete(ob.includedTxResults, outboundID)

		// add new txHash and txResult
		ob.tssOutboundHashes[txHash] = true
		ob.includedTxResults[outboundID] = getTxResult
	}
}

// GetIncludedTx gets the receipt and transaction from memory
func (ob *Observer) GetIncludedTx(nonce uint64) *btcjson.GetTransactionResult {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.includedTxResults[ob.OutboundID(nonce)]
}

// Basic TSS outbound checks:
//   - confirmations >= 0
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
		return errors.Wrapf(err, "checkTssOutboundResult: error GetRawTxResult %s", hash.String())
	}
	err = ob.checkTSSVin(ctx, rawResult.Vin, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutboundResult: invalid TSS Vin in outbound %s nonce %d", hash, nonce)
	}

	// differentiate between normal and cancelled cctx
	if compliance.IsCctxRestricted(cctx) || params.Amount.Uint64() < constant.BTCWithdrawalDustAmount {
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
	pubKeyTss := hex.EncodeToString(ob.TSS().PubKey().Bytes(true))
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
			preTxid, err := ob.getOutboundHashByNonce(ctx, nonce-1, false)
			if err != nil {
				return fmt.Errorf("checkTSSVin: error getOutboundHashByNonce %d", nonce-1)
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
	tssAddress := ob.TSSAddressString()
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
	tssAddress := ob.TSSAddressString()
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
