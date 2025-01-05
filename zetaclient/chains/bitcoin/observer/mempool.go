package observer

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// LastStuckOutbound contains the last stuck outbound tx information.
type LastStuckOutbound struct {
	// Nonce is the nonce of the outbound.
	Nonce uint64

	// Tx is the original transaction.
	Tx *btcutil.Tx

	// StuckFor is the duration for which the tx has been stuck.
	StuckFor time.Duration
}

// NewLastStuckOutbound creates a new LastStuckOutbound struct.
func NewLastStuckOutbound(nonce uint64, tx *btcutil.Tx, stuckFor time.Duration) *LastStuckOutbound {
	return &LastStuckOutbound{
		Nonce:    nonce,
		Tx:       tx,
		StuckFor: stuckFor,
	}
}

// WatchMempoolTxs monitors pending outbound txs in the Bitcoin mempool.
func (ob *Observer) WatchMempoolTxs(ctx context.Context) error {
	task := func(ctx context.Context, _ *ticker.Ticker) error {
		if err := ob.refreshLastStuckOutbound(ctx); err != nil {
			ob.Logger().Chain.Err(err).Msg("refreshLastStuckOutbound error")
		}
		return nil
	}

	return ticker.Run(
		ctx,
		common.MempoolStuckTxCheckInterval,
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().Chain, "WatchMempoolTxs"),
	)
}

// refreshLastStuckOutbound refreshes the information about the last stuck tx in the Bitcoin mempool.
func (ob *Observer) refreshLastStuckOutbound(ctx context.Context) error {
	// log fields
	lf := map[string]any{
		logs.FieldMethod: "refreshLastStuckOutbound",
	}

	// step 1: get last TSS transaction
	lastTx, lastNonce, err := ob.GetLastOutbound(ctx)
	if err != nil {
		return errors.Wrap(err, "GetLastOutbound failed")
	}
	txHash := lastTx.MsgTx().TxID()
	lf[logs.FieldNonce] = lastNonce
	lf[logs.FieldTx] = txHash
	ob.logger.Outbound.Info().Fields(lf).Msg("checking last TSS outbound")

	// step 2: is last tx stuck in mempool?
	stuck, stuckFor, err := rpc.IsTxStuckInMempool(ob.btcClient, txHash, rpc.PendingTxFeeBumpWaitBlocks)
	if err != nil {
		return errors.Wrapf(err, "cannot determine if tx %s nonce %d is stuck", txHash, lastNonce)
	}

	// step 3: update last outbound stuck tx information
	//
	// the key ideas to determine if Bitcoin outbound is stuck/unstuck:
	// 	1. outbound txs are a sequence of txs chained by nonce-mark UTXOs.
	//  2. outbound tx with nonce N+1 MUST spend the nonce-mark UTXO produced by parent tx with nonce N.
	//  3. when the last descendant tx is stuck, none of its ancestor txs can go through, so the stuck flag is set.
	//  4. then RBF kicks in, it bumps the fee of the last descendant tx and aims to increase the average fee
	//     rate of the whole tx chain (as a package) to make it attractive to miners.
	//  5. after RBF replacement, zetaclient clears the stuck flag immediately, hoping the new tx will be included
	//     within next 'PendingTxFeeBumpWaitBlocks' blocks.
	//  6. the new tx may get stuck again (e.g. surging traffic) after 'PendingTxFeeBumpWaitBlocks' blocks, and
	//     the stuck flag will be set again to trigger another RBF, and so on.
	//  7. all pending txs will be eventually cleared by fee bumping, and the stuck flag will be cleared.
	//
	// Note: reserved RBF bumping fee might be not enough to clear the stuck txs during extreme traffic surges, two options:
	//  1. wait for the gas rate to drop.
	//  2. manually clear the stuck txs by using offline accelerator services.
	if stuck {
		ob.SetLastStuckOutbound(NewLastStuckOutbound(lastNonce, lastTx, stuckFor))
	} else {
		ob.SetLastStuckOutbound(nil)
	}

	return nil
}

// GetLastOutbound gets the last outbound (with highest nonce) that had been sent to Bitcoin network.
// Bitcoin outbound txs can be found from two sources:
//  1. txs that had been reported to tracker and then checked and included by this observer self.
//  2. txs that had been broadcasted by this observer self.
//
// Once 2/3+ of the observers reach consensus on last outbound, RBF will start.
func (ob *Observer) GetLastOutbound(ctx context.Context) (*btcutil.Tx, uint64, error) {
	var (
		lastNonce uint64
		lastHash  string
	)

	// wait for pending nonce to refresh
	pendingNonce := ob.GetPendingNonce()
	if ob.GetPendingNonce() == 0 {
		return nil, 0, errors.New("pending nonce is zero")
	}

	// source 1:
	// pick highest nonce tx from included txs
	lastNonce = pendingNonce - 1
	txResult := ob.GetIncludedTx(lastNonce)
	if txResult == nil {
		// should NEVER happen by design
		return nil, 0, errors.New("last included tx not found")
	}
	lastHash = txResult.TxID

	// source 2:
	// pick highest nonce tx from broadcasted txs
	p, err := ob.ZetacoreClient().GetPendingNoncesByChain(ctx, ob.Chain().ChainId)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetPendingNoncesByChain failed")
	}
	for nonce := uint64(p.NonceLow); nonce < uint64(p.NonceHigh); nonce++ {
		if nonce > lastNonce {
			txID, found := ob.GetBroadcastedTx(nonce)
			if found {
				lastNonce = nonce
				lastHash = txID
			}
		}
	}

	// ensure this nonce is the REAL last transaction
	// cross-check the latest UTXO list, the nonce-mark utxo exists ONLY for last nonce
	if ob.FetchUTXOs(ctx) != nil {
		return nil, 0, errors.New("FetchUTXOs failed")
	}
	if _, err = ob.findNonceMarkUTXO(lastNonce, lastHash); err != nil {
		return nil, 0, errors.Wrapf(err, "findNonceMarkUTXO failed for last tx %s nonce %d", lastHash, lastNonce)
	}

	// query last transaction
	// 'GetRawTransaction' is preferred over 'GetTransaction' here for three reasons:
	//  1. it can fetch both stuck tx and non-stuck tx as far as they are valid txs.
	//  2. it never fetch invalid tx (e.g., old tx replaced by RBF), so we can exclude invalid ones.
	//  3. zetaclient needs the original tx body of a stuck tx to bump its fee and sign again.
	lastTx, err := rpc.GetRawTxByHash(ob.btcClient, lastHash)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "GetRawTxByHash failed for last tx %s nonce %d", lastHash, lastNonce)
	}

	return lastTx, lastNonce, nil
}
