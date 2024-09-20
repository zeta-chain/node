package observer

import (
	"context"
	"fmt"
	"slices"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

const (
	// MaxTransactionsPerTick is the maximum number of transactions to process on a ticker
	MaxTransactionsPerTick = 100
)

func (ob *Observer) watchInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	var (
		chainID         = ob.Chain().ChainId
		initialInterval = ticker.SecondsFromUint64(ob.ChainParams().InboundTicker)
		sampledLogger   = ob.Logger().Inbound.Sample(&zerolog.BasicSampler{N: 10})
	)

	ob.Logger().Inbound.Info().Msgf("WatchInbound started for chain %d", chainID)

	task := func(ctx context.Context, t *ticker.Ticker) error {
		if !app.IsInboundObservationEnabled() {
			sampledLogger.Info().Msgf("WatchInbound: inbound observation is disabled for chain %d", chainID)
			return nil
		}

		if err := ob.observeInbound(ctx, app); err != nil {
			ob.Logger().Inbound.Err(err).Msg("WatchInbound: observeInbound error")
		}

		newInterval := ticker.SecondsFromUint64(ob.ChainParams().InboundTicker)
		t.SetInterval(newInterval)

		return nil
	}

	return ticker.Run(
		ctx,
		initialInterval,
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().Inbound, "WatchInbound"),
	)
}

// Flow:
// - [x] Ensure last scanned transaction is set
// - [x] Get all transaction between [lastScannedTx; now]
// - [ ] Filter only valid and inbound transactions
// - [ ] For each transaction (ordered by *ASC*)
//   - [ ] Construct crosschain cosmos message
//   - [ ] Vote
//   - [ ] Save last scanned tx
func (ob *Observer) observeInbound(ctx context.Context, _ *zctx.AppContext) error {
	if err := ob.ensureLastScannedTX(ctx); err != nil {
		return errors.Wrap(err, "unable to ensure last scanned tx")
	}

	lt, hashBits, err := liteapi.TransactionHashFromString(ob.LastTxScanned())
	if err != nil {
		return errors.Wrapf(err, "unable to parse last scanned tx %q", ob.LastTxScanned())
	}

	txs, err := ob.client.GetTransactionsUntil(ctx, ob.gatewayID, lt, hashBits)
	if err != nil {
		return errors.Wrap(err, "unable to get transactions")
	}

	// Process from oldest to latest (ASC)
	slices.Reverse(txs)

	switch {
	case len(txs) == 0:
		// noop
		return nil
	case len(txs) > MaxTransactionsPerTick:
		ob.Logger().Inbound.Info().
			Msgf("ObserveInbound: got %d transactions. Taking first %d", len(txs), MaxTransactionsPerTick)

		txs = txs[:MaxTransactionsPerTick]
	default:
		ob.Logger().Inbound.Info().Msgf("ObserveInbound: got %d transactions", len(txs))
	}

	// todo deploy sample GW to testnet
	// todo send some TON and test

	// todo FilterInboundEvent

	for _, tx := range txs {
		fmt.Println("TON TX", tx)
	}

	/*
		// loop signature from oldest to latest to filter inbound events
		for i := len(signatures) - 1; i >= 0; i-- {
			sig := signatures[i]
			sigString := sig.Signature.String()

			// process successfully signature only
			if sig.Err == nil {
				txResult, err := ob.solClient.GetTransaction(ctx, sig.Signature, &rpc.GetTransactionOpts{})
				if err != nil {
					// we have to re-scan this signature on next ticker
					return errors.Wrapf(err, "error GetTransaction for chain %d sig %s", chainID, sigString)
				}

				// filter inbound events and vote
				err = ob.FilterInboundEventsAndVote(ctx, txResult)
				if err != nil {
					// we have to re-scan this signature on next ticker
					return errors.Wrapf(err, "error FilterInboundEventAndVote for chain %d sig %s", chainID, sigString)
				}
			}

			// signature scanned; save last scanned signature to both memory and db, ignore db error
			if err := ob.SaveLastTxScanned(sigString, sig.Slot); err != nil {
				ob.Logger().
					Inbound.Error().
					Err(err).
					Msgf("ObserveInbound: error saving last sig %s for chain %d", sigString, chainID)
			}
			ob.Logger().
				Inbound.Info().
				Msgf("ObserveInbound: last scanned sig is %s for chain %d in slot %d", sigString, chainID, sig.Slot)
	*/

	return nil
}

func (ob *Observer) ensureLastScannedTX(ctx context.Context) error {
	// noop
	if ob.LastTxScanned() != "" {
		return nil
	}

	tx, _, err := ob.client.GetFirstTransaction(ctx, ob.gatewayID)
	if err != nil {
		return err
	}

	txHash := liteapi.TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))

	ob.WithLastTxScanned(txHash)

	return ob.WriteLastTxScannedToDB(txHash)
}
