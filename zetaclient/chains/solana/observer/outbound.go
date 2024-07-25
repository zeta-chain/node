package observer

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// GetTxID returns a unique id for Solana outbound
func (ob *Observer) GetTxID(nonce uint64) string {
	tssAddr := ob.TSS().EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.Chain().ChainId, tssAddr, nonce)
}

// WatchOutbound watches evm chain for outgoing txs status
// TODO(revamp): move ticker function to ticker file
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchOutbound(ctx context.Context) error {
	// get app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	// create outbound ticker based on chain params
	chainID := ob.Chain().ChainId
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchOutbound_%d", chainID),
		ob.GetChainParams().OutboundTicker,
	)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Msg("error creating ticker")
		return err
	}

	ob.Logger().Outbound.Info().Msgf("WatchOutbound started for chain %d", chainID)
	sampledLogger := ob.Logger().Outbound.Sample(&zerolog.BasicSampler{N: 10})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !app.IsOutboundObservationEnabled(ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutbound: outbound observation is disabled for chain %d", chainID)
				continue
			}
			trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
			if err != nil {
				continue
			}

			for _, tracker := range trackers {
				// go to next tracker if this one already has a finalized tx
				nonce := tracker.Nonce
				if ob.IsTxFinalized(tracker.Nonce) {
					continue
				}

				txCount := 0
				var txResult *rpc.GetTransactionResult
				for _, txHash := range tracker.HashList {
					if result, ok := ob.checkFinalizedTx(ctx, txHash.TxHash, nonce); ok {
						txCount++
						txResult = result
						ob.Logger().Outbound.Info().
							Msgf("WatchOutbound: confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
						if txCount > 1 {
							ob.Logger().Outbound.Error().Msgf(
								"WatchOutbound: checkFinalizedTx passed, txCount %d chain %d nonce %d txResult %v", txCount, chainID, nonce, txResult)
						}
					}
				}
				// should be only one finalized txHash for each nonce
				if txCount == 1 {
					ob.SetTxResult(nonce, txResult)
				} else if txCount > 1 {
					// should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.Logger().Outbound.Error().Msgf("WatchOutbound: finalized multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.Logger().Outbound)
		case <-ob.StopChannel():
			ob.Logger().Outbound.Info().Msgf("WatchOutbound: watcher stopped for chain %d", chainID)
			return nil
		}
	}
}

// IsOutboundProcessed checks outbound status and returns (isIncluded, isConfirmed, error)
// NOTE: There is a critical difference from EVM/Bitcoin chains regarding nonce and transaction status
// On EVM/Bitcoin chains, for each scheduled outbound nonce, there can be exactly 1 tx (successful or failed)
// included in the blockchain.  On Solana, this is no longer the case: there can be AT MOST 1 SUCCESSFUL tx
// corresponding to a scheduled outbound nonce. However, there can be multiple FAILED txs corresponding to
// the same nonce. Therefore, we must distinguish between failed tx due to 1) nonce rejected; 2) tx reverted.
// The first case, we should ignore it (in other chains, this tx should not be mined at all). The second case,
// we should report it to zetacore as a real failed tx.
// FIXME: implement the above logic that distinguishes between nonce-rejected and tx-reverted failed tx.
func (ob *Observer) IsOutboundProcessed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, bool, error) {
	// get outbound params from cctx
	params := cctx.GetCurrentOutboundParam()
	if params == nil {
		return false, false, fmt.Errorf("outbound param not found from cctx")
	}

	// skip if outbound is not finalized yet
	nonce := params.TssNonce
	txResult := ob.GetTxResult(nonce)
	if txResult == nil {
		return false, false, nil
	}

	// extract tx signature
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return false, false, errors.Wrapf(err, "GetTransaction error for nonce %d", nonce)
	}
	txSig := tx.Signatures[0]

	outboundAmount := params.Amount.BigInt()      // FIXME: parse this amount from txRes itself, not from cctx
	outboundStatus := chains.ReceiveStatus_failed // tx was failed/reverted: FIXME: see the note in function comment
	coinType := cctx.InboundParams.CoinType
	if txResult.Meta.Err == nil { // this indicates tx was successful
		outboundStatus = chains.ReceiveStatus_success
	}

	// post vote to zetacore
	ob.PostVoteOutbound(ctx, cctx.Index, txSig.String(), txResult, outboundAmount, outboundStatus, nonce, coinType)

	// tracker, err := ob.ZetacoreClient().GetOutboundTracker(ctx, ob.Chain(), nonce)
	// if err != nil {
	// 	return false, false, nil
	// }
	// for _, hash := range tracker.HashList {
	// 	sig, err := solana.SignatureFromBase58(hash.TxHash)
	// 	if err != nil {
	// 		ob.Logger().Outbound.Warn().Err(err).Msgf("solana.SignatureFromBase58 error: %s", hash.TxHash)
	// 		continue
	// 	}
	// 	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
	// 		Commitment: rpc.CommitmentFinalized, // must be finalized state so this tx will not be re-org'ed later
	// 	})
	// 	if err != nil {
	// 		ob.Logger().Outbound.Warn().Err(err).Msgf("solana.GetTransaction error: %s", hash.TxHash)
	// 		continue
	// 	}

	// 	outboundAmount := params.Amount.BigInt()      // FIXME: parse this amount from txRes itself, not from cctx
	// 	outboundStatus := chains.ReceiveStatus_failed // tx was failed/reverted: FIXME: see the note in function comment
	// 	coinType := cctx.InboundParams.CoinType
	// 	if txResult.Meta.Err == nil { // this indicates tx was successful
	// 		outboundStatus = chains.ReceiveStatus_success
	// 	}

	// 	// post vote to zetacore
	// 	ob.PostVoteOutbound(
	// 		ctx,
	// 		cctx.Index,
	// 		sig.String(),
	// 		txResult,
	// 		outboundAmount,
	// 		outboundStatus,
	// 		nonce,
	// 		coinType,
	// 	)

	// 	break // outbound tx confirmed on Solana and reported to zetacore; skip the rest of the list
	// }

	return true, true, nil
}

// PostVoteOutbound posts vote to zetacore for the finalized outbound
func (ob *Observer) PostVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	outboundHash string,
	txResult *rpc.GetTransactionResult,
	valueReceived *big.Int,
	status chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
) {
	chainID := ob.Chain().ChainId
	computeUnitsConsumed := uint64(0)
	cuPrice := big.NewInt(0)
	if txResult.Meta.ComputeUnitsConsumed == nil {
		ob.Logger().Outbound.Warn().Msgf("solana.GetTransaction: compute units consumed is nil")
	} else {
		computeUnitsConsumed = *txResult.Meta.ComputeUnitsConsumed
		if computeUnitsConsumed <= 0 {
			ob.Logger().Outbound.Warn().Msgf("solana.GetTransaction: compute units consumed is %d", computeUnitsConsumed)
			computeUnitsConsumed = 5001 // default to 5000, for a single signature tx; make it 5001 to distinguish
		}
		cuPrice.SetUint64(txResult.Meta.Fee / computeUnitsConsumed)
	}

	creator := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		creator.String(),
		cctxIndex,
		outboundHash,
		txResult.Slot, // TODO: check this; is slot equivalent to block height?
		computeUnitsConsumed,
		math.NewIntFromBigInt(cuPrice),
		200_000,                               // this is default compute unit budget;
		math.NewUintFromBigInt(valueReceived), // FIXME: parse this amount from txRes itself, not from cctx
		status,
		chainID,
		nonce, // FIXME: parse this from the txRes/tx ?
		coinType,
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit
	var retryGasLimit uint64
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	logFields := map[string]any{
		"chain":    chainID,
		"nonce":    nonce,
		"outbound": outboundHash,
	}
	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Fields(logFields).Msg("PostVoteOutbound: error posting outbound vote")
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logFields["vote"] = zetaTxHash
		logFields["ballot"] = ballot
		ob.Logger().Outbound.Info().Fields(logFields).Msg("PostVoteOutbound: posted outbound vote successfully")
	}
}

// checkFinalizedTx checks if a txHash is finalized for given nonce
// returns (tx result, true) if finalized or (nil, false) otherwise
func (ob *Observer) checkFinalizedTx(ctx context.Context, txHash string, _ uint64) (*rpc.GetTransactionResult, bool) {
	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msgf("checkFinalizedTx: err SignatureFromBase58 for tx hash %s", txHash)
		return nil, false
	}

	// get tx result with finalized commitment to avoid re-org
	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msgf("checkFinalizedTx: err GetTransaction for sig %s", sig)
		return nil, false
	}

	// TODO:
	// - decode tx message and 'withdraw' instruction
	// - check tx sender and nonce
	// - for successful tx, simply return txResult
	// - for failed tx, find out if it's due to nonce rejection, signature rejection, or reverted tx
	//     - ignore nonce rejection, report signature rejection
	//     - report reverted tx back to zetacore

	return txResult, true
}
