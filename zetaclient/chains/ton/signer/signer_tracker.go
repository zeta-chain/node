package signer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// trackOutbound tracks sent external message and records it as outboundTracker.
// Explanation:
// Due to TON's nature, it's not possible to get tx hash before it's confirmed on-chain,
// So we need to poll from the latest account state (prevState) up to the most recent tx
// and search for desired tx hash. After it's found, we can record it as outboundTracker.
//
// Note that another zetaclient observers that scrolls Gateway's txs can publish this tracker concurrently.
func (s *Signer) trackOutbound(
	ctx context.Context,
	zetacore interfaces.ZetacoreClient,
	outbound outbound,
	prevState rpc.Account,
) error {
	metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Inc()
	defer metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Dec()

	const (
		timeout = 60 * time.Second
		tick    = time.Second
	)

	var (
		start   = time.Now()
		chainID = s.Chain().ChainId

		acc   = s.gateway.AccountID()
		lt    = prevState.LastTxLT
		hash  = ton.Bits256(prevState.LastTxHash)
		nonce = uint64(outbound.seqno)

		filter = outboundFilter(outbound)
	)

	for time.Since(start) <= timeout {
		txs, err := s.tonClient.GetTransactionsSince(ctx, acc, lt, hash)
		if err != nil {
			return errors.Wrapf(err, "unable to get transactions (lt %d, hash %s)", lt, hash.Hex())
		}

		results := s.gateway.ParseAndFilterMany(txs, filter)
		if len(results) == 0 {
			time.Sleep(tick)
			continue
		}

		tx := results[0].Transaction
		txHash := rpc.TransactionToHashString(results[0].Transaction)

		if !tx.IsSuccess() {
			// should not happen
			return errors.Errorf("transaction %q is not successful", txHash)
		}

		// Note that this method has a check for noop
		_, err = zetacore.PostOutboundTracker(ctx, chainID, nonce, txHash)
		if err != nil {
			return errors.Wrap(err, "unable to add outbound tracker")
		}

		return nil
	}

	return errors.Errorf("timeout exceeded (%s)", time.Since(start).String())
}

// creates a tx filter for this very outbound tx
func outboundFilter(ob outbound) func(tx *toncontracts.Transaction) (found bool) {
	return func(tx *toncontracts.Transaction) bool {
		auth, err := tx.OutboundAuth()

		return err == nil && auth.Seqno == ob.seqno && auth.Sig == ob.message.Signature()
	}
}
