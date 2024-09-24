package observer

import (
	"context"
	"fmt"
	"slices"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/zetacore"
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

func (ob *Observer) observeInbound(ctx context.Context, _ *zctx.AppContext) error {
	if err := ob.ensureLastScannedTX(ctx); err != nil {
		return errors.Wrap(err, "unable to ensure last scanned tx")
	}

	lt, hashBits, err := liteapi.TransactionHashFromString(ob.LastTxScanned())
	if err != nil {
		return errors.Wrapf(err, "unable to parse last scanned tx %q", ob.LastTxScanned())
	}

	txs, err := ob.client.GetTransactionsUntil(ctx, ob.gateway.AccountID(), lt, hashBits)
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

	for i := range txs {
		tx := txs[i]

		parsedTX, skip, err := ob.gateway.ParseAndFilter(tx, toncontracts.FilterInbounds)
		if err != nil {
			return errors.Wrap(err, "unable to parse and filter tx")
		}

		if skip {
			ob.setLastScannedTX(&tx)
			continue
		}

		if _, err := ob.voteInbound(ctx, parsedTX); err != nil {
			return errors.Wrapf(err, "unable to vote inbound (hash %s)", parsedTX.Hash().Hex())
		}

		ob.setLastScannedTX(&parsedTX.Transaction)
	}

	return nil
}

func (ob *Observer) voteInbound(ctx context.Context, tx *toncontracts.Transaction) (string, error) {
	// noop
	if tx.Operation == toncontracts.OpDonate {
		ob.Logger().Inbound.Info().
			Uint64("tx.lt", tx.Lt).
			Str("tx.hash", tx.Hash().Hex()).
			Msg("Thank you rich folk for your donation!")

		return "", nil
	}

	// todo add compliance check
	//   https://github.com/zeta-chain/node/issues/2916

	blockHeader, err := ob.client.GetBlockHeader(ctx, tx.BlockID, 0)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get block header %s", tx.BlockID.String())
	}

	sender, amount, memo, err := extractInboundData(tx)
	if err != nil {
		return "", err
	}

	seqno := blockHeader.MinRefMcSeqno

	return ob.voteDeposit(ctx, tx, sender, amount, memo, seqno)
}

func extractInboundData(tx *toncontracts.Transaction) (string, math.Uint, []byte, error) {
	if tx.Operation == toncontracts.OpDeposit {
		d, err := tx.Deposit()
		if err != nil {
			return "", math.NewUint(0), nil, err
		}

		return d.Sender.ToRaw(), d.Amount, d.Memo(), nil
	}

	if tx.Operation == toncontracts.OpDepositAndCall {
		d, err := tx.DepositAndCall()
		if err != nil {
			return "", math.NewUint(0), nil, err
		}

		return d.Sender.ToRaw(), d.Amount, d.Memo(), nil
	}

	return "", math.NewUint(0), nil, fmt.Errorf("unknown operation %d", tx.Operation)
}

func (ob *Observer) voteDeposit(
	ctx context.Context,
	tx *toncontracts.Transaction,
	sender string,
	amount math.Uint,
	memo []byte,
	seqno uint32,
) (string, error) {
	const (
		eventIndex    = 0 // not a smart contract call
		coinType      = coin.CoinType_Gas
		asset         = "" // empty for gas coin
		gasLimit      = 0
		retryGasLimit = zetacore.PostVoteInboundExecutionGasLimit
	)

	var (
		operatorAddress = ob.ZetacoreClient().GetKeys().GetOperatorAddress()
		inboundHash     = liteapi.TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))
	)

	msg := zetacore.GetInboundVoteMessage(
		sender,
		ob.Chain().ChainId,
		sender,
		sender,
		ob.ZetacoreClient().Chain().ChainId,
		amount,
		string(memo),
		inboundHash,
		uint64(seqno),
		gasLimit,
		coinType,
		asset,
		operatorAddress.String(),
		eventIndex,
	)

	return ob.PostVoteInbound(ctx, msg, retryGasLimit)
}

func (ob *Observer) ensureLastScannedTX(ctx context.Context) error {
	// noop
	if ob.LastTxScanned() != "" {
		return nil
	}

	tx, _, err := ob.client.GetFirstTransaction(ctx, ob.gateway.AccountID())
	if err != nil {
		return err
	}

	ob.setLastScannedTX(tx)

	return nil
}

func (ob *Observer) setLastScannedTX(tx *ton.Transaction) {
	txHash := liteapi.TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))

	ob.WithLastTxScanned(txHash)

	if err := ob.WriteLastTxScannedToDB(txHash); err != nil {
		ob.Logger().Inbound.Error().
			Err(err).
			Uint64("tx.lt", tx.Lt).
			Str("tx.hash", tx.Hash().Hex()).
			Msgf("ObserveInbound: unable to WriteLastTxScannedToDB")

		return
	}

	ob.Logger().Inbound.Info().
		Uint64("tx.lt", tx.Lt).
		Str("tx.hash", tx.Hash().Hex()).
		Msgf("ObserveInbound: WriteLastTxScannedToDB")
}
