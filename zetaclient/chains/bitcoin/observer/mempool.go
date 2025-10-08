package observer

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/logs"
)

// number of blocks to await before considering a tx stuck in mempool
const (
	pendingTxFeeBumpWaitBlocks       = 3
	pendingTxFeeBumpWaitBlocksRegnet = 30
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

// ObserveMempool monitors pending outbound txs in the mempool.
func (ob *Observer) ObserveMempool(ctx context.Context) error {
	err := ob.refreshLastStuckOutbound(ctx)
	if err != nil {
		return errors.Wrap(err, "refresh last stuck outbound failed")
	}

	return nil
}

// refreshLastStuckOutbound refreshes the information about the last stuck tx in the Bitcoin mempool.
// Once 2/3+ of the observers reach consensus on last stuck outbound, RBF will start.
func (ob *Observer) refreshLastStuckOutbound(ctx context.Context) error {
	logger := ob.logger.Outbound

	pendingTxFinder := ob.getLastPendingOutbound
	if custom, ok := pendingTxFinderFromContext(ctx); ok {
		pendingTxFinder = custom
	}

	// step 1: get last TSS transaction
	lastTx, lastNonce, err := pendingTxFinder(ctx)
	if err != nil {
		logger.Info().Err(err).Msg("last pending outbound not found")
		return nil
	}
	lastTxHash := lastTx.MsgTx().TxID()

	logger.Info().
		Uint64(logs.FieldNonce, lastNonce).
		Str(logs.FieldTx, lastTxHash).
		Msg("checking last TSS outbound")

	// step 2: is last tx stuck in mempool?
	stuck, stuckFor, err := ob.bitcoinClient.IsTxStuckInMempool(ctx, lastTxHash, ob.feeBumpWaitBlocks)
	if err != nil {
		return errors.Wrapf(err, "cannot determine if tx %s nonce %d is stuck", lastTxHash, lastNonce)
	}

	// step 3: update last outbound stuck tx information
	//
	// The key ideas to determine if Bitcoin outbound is stuck/unstuck:
	//
	//  1. Outbound txs are a sequence of txs chained by nonce-mark UTXOs.
	//
	//  2. Outbound tx with nonce N+1 MUST spend the nonce-mark UTXO produced by parent tx with
	//     nonce N.
	//
	//  3. When the last descendant tx is stuck, none of its ancestor txs can go through, so the
	//     stuck flag is set.
	//
	//  4. Then RBF kicks in, it bumps the fee of the last descendant tx and aims to increase the
	//     average fee rate of the whole tx chain (as a package) to make it attractive to miners.
	//
	//  5. After RBF replacement, zetaclient clears the stuck flag immediately, hoping the new tx
	//     will be included within next 'PendingTxFeeBumpWaitBlocks' blocks.
	//
	//  6. The new tx may get stuck again (e.g. surging traffic) after 'PendingTxFeeBumpWaitBlocks'
	//     blocks, and the stuck flag will be set again to trigger another RBF, and so on.
	//
	//  7. All pending txs will be eventually cleared by fee bumping, and the stuck flag will be
	//     cleared.
	//
	// Note: reserved RBF bumping fee might be not enough to clear the stuck txs during extreme
	// traffic surges, two options:
	//  1. wait for the gas rate to drop.
	//  2. manually clear the stuck txs by using transaction accelerator services.
	var stuckOutbound *LastStuckOutbound
	if stuck {
		stuckOutbound = newLastStuckOutbound(lastNonce, lastTx, stuckFor)
	}

	ob.setLastStuckOutbound(stuckOutbound)

	return nil
}

// GetLastPendingOutbound gets the last pending outbound (with highest nonce) that sits in the Bitcoin mempool.
// Bitcoin outbound txs can be found from two sources:
//  1. txs that had been reported to tracker and then checked and included by this observer self.
//  2. txs that had been broadcasted by this observer self.
//
// Returns error if last pending outbound is not found
func (ob *Observer) getLastPendingOutbound(ctx context.Context) (tx *btcutil.Tx, nonce uint64, err error) {
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
	txResult := ob.GetIncludedTx(pendingNonce - 1)
	if txResult != nil {
		lastNonce = pendingNonce - 1
		lastHash = txResult.TxID
	}

	// source 2:
	// pick highest nonce tx from broadcasted txs
	p, err := ob.ZetaRepo().GetPendingNonces(ctx)
	if err != nil {
		return nil, 0, err
	}

	// #nosec G115 always in range
	for nonce := uint64(p.NonceLow); nonce < uint64(p.NonceHigh); nonce++ {
		if nonce > lastNonce {
			txID, found := ob.GetBroadcastedTx(nonce)
			if found {
				lastNonce = nonce
				lastHash = txID
			}
		}
	}

	// stop if last tx not found, and it is okay
	// this individual zetaclient lost track of the last tx for some reason (offline, db reset, etc.)
	if lastNonce == 0 {
		return nil, 0, errors.New("last tx not found")
	}

	// is tx in the mempool?
	if _, err = ob.bitcoinClient.GetMempoolEntry(ctx, lastHash); err != nil {
		return nil, 0, errors.Wrapf(err, "last tx %s is not in mempool", lastHash)
	}

	// ensure this tx is the REAL last transaction
	// cross-check the latest UTXO list, the nonce-mark utxo exists ONLY for last nonce
	if err := ob.FetchUTXOs(ctx); err != nil {
		return nil, 0, errors.Wrap(err, "unable to fetch UTXOs")
	}

	if _, err = ob.findNonceMarkUTXO(lastNonce, lastHash); err != nil {
		return nil, 0, errors.Wrapf(err, "findNonceMarkUTXO failed for last tx %s nonce %d", lastHash, lastNonce)
	}

	// query last transaction
	// 'GetRawTransaction' is preferred over 'GetTransaction' here for three reasons:
	//  1. it can fetch both stuck tx and non-stuck tx as far as they are valid txs.
	//  2. it never fetch invalid tx (e.g., old tx replaced by RBF), so we can exclude invalid ones.
	//  3. zetaclient needs the original tx body of a stuck tx to bump its fee and sign again.
	lastTx, err := ob.bitcoinClient.GetRawTransactionByStr(ctx, lastHash)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "GetRawTransactionByStr failed for last tx %s nonce %d", lastHash, lastNonce)
	}

	return lastTx, lastNonce, nil
}

// newLastStuckOutbound creates a new LastStuckOutbound struct.
func newLastStuckOutbound(nonce uint64, tx *btcutil.Tx, stuckFor time.Duration) *LastStuckOutbound {
	return &LastStuckOutbound{Nonce: nonce, Tx: tx, StuckFor: stuckFor}
}

// allows to override observer's pending tx finder.
// useful for testing because the default implementation is complex and requires thorough mocking.
type (
	pendingTxFinder    func(ctx context.Context) (*btcutil.Tx, uint64, error)
	pendingTxFinderKey struct{}
)

func pendingTxFinderFromContext(ctx context.Context) (pendingTxFinder, bool) {
	fn, ok := ctx.Value(pendingTxFinderKey{}).(pendingTxFinder)

	return fn, ok
}
