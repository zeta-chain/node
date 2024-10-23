package signer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
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
	w *toncontracts.Withdrawal,
	prevState tlb.ShardAccount,
) error {
	const (
		timeout = 60 * time.Second
		tick    = time.Second
	)

	var (
		start   = time.Now()
		chainID = s.Chain().ChainId

		acc   = s.gateway.AccountID()
		lt    = prevState.LastTransLt
		hash  = ton.Bits256(prevState.LastTransHash)
		nonce = uint64(w.Seqno)

		filter = withdrawalFilter(w)
	)

	for time.Since(start) <= timeout {
		txs, err := s.client.GetTransactionsSince(ctx, acc, lt, hash)
		if err != nil {
			return errors.Wrapf(err, "unable to get transactions (lt %d, hash %s)", lt, hash.Hex())
		}

		results := s.gateway.ParseAndFilterMany(txs, filter)
		if len(results) == 0 {
			time.Sleep(tick)
			continue
		}

		tx := results[0]
		txHash := liteapi.TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))

		// Note that this method has a check for noop
		_, err = zetacore.AddOutboundTracker(ctx, chainID, nonce, txHash, nil, "", 0)
		if err != nil {
			return errors.Wrap(err, "unable to add outbound tracker")
		}

		return nil
	}

	return errors.Errorf("timeout exceeded (%s)", time.Since(start).String())
}

// creates a tx filter for this very withdrawal
func withdrawalFilter(w *toncontracts.Withdrawal) func(tx *toncontracts.Transaction) bool {
	return func(tx *toncontracts.Transaction) bool {
		if !tx.IsOutbound() || tx.Operation != toncontracts.OpWithdraw {
			return false
		}

		wd, err := tx.Withdrawal()
		if err != nil {
			return false
		}

		return wd.Seqno == w.Seqno && wd.Sig == w.Sig
	}
}
