package bitcoin

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/common"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

type Bitcoin struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *Bitcoin {
	return &Bitcoin{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

func (b *Bitcoin) Chain() chains.Chain {
	return b.observer.Chain()
}

func (b *Bitcoin) Start(ctx context.Context) error {
	if ok := b.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := b.observer.ZetacoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(b.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(b.observer.ChainParams().GasPriceTicker)
	})

	optUTXOInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(b.observer.ChainParams().WatchUtxoTicker)
	})

	optMempoolInterval := scheduler.Interval(common.BTCMempoolStuckTxCheckInterval)

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(b.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !b.observer.ChainParams().IsSupported
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(b.group()),
			scheduler.Name(name),
		}, opts...)

		b.scheduler.Register(ctx, exec, opts...)
	}

	// Observers
	register(b.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(b.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(b.observer.FetchUTXOs, "fetch_utxos", optUTXOInterval, optGenericSkipper)
	register(b.observer.ObserveMempool, "observe_mempool", optMempoolInterval, optGenericSkipper)
	register(b.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(b.observer.CheckRPCStatus, "check_rpc_status")
	register(b.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(b.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

func (b *Bitcoin) Stop() {
	b.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	b.scheduler.StopGroup(b.group())
}

func (b *Bitcoin) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("btc:%d", b.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules pending cross-chain transactions on NEW zeta blocks
// 1. schedule at most one keysign per ticker
// 2. schedule keysign only when nonce-mark UTXO is available
// 3. stop keysign when lookahead is reached
func (b *Bitcoin) scheduleCCTX(ctx context.Context) error {
	if err := b.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta block from context")
	}

	time.Sleep(delay)

	// #nosec G115 always in range
	zetaHeight := uint64(zetaBlock.Block.Height)
	chain := b.observer.Chain()
	chainID := chain.ChainId
	lookahead := b.observer.ChainParams().OutboundScheduleLookahead
	// #nosec G115 positive
	scheduleInterval := uint64(b.observer.ChainParams().OutboundScheduleInterval)

	cctxList, _, err := b.observer.ZetacoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		var (
			params     = cctx.GetCurrentOutboundParam()
			nonce      = params.TssNonce
			outboundID = base.OutboundIDFromCCTX(cctx)
		)

		if params.ReceiverChainId != chainID {
			b.outboundLogger(outboundID).Error().Msg("Schedule CCTX: chain id mismatch")

			continue
		}

		// try confirming the outbound
		continueKeysign, err := b.observer.VoteOutboundIfConfirmed(ctx, cctx)

		switch {
		case err != nil:
			b.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX: VoteOutboundIfConfirmed failed")
			continue
		case !continueKeysign:
			b.outboundLogger(outboundID).Info().Msg("Schedule CCTX: outbound already processed")
			continue
		case nonce > b.observer.GetPendingNonce():
			// stop if the nonce being processed is higher than the pending nonce
			return nil
		case int64(idx) >= lookahead:
			// stop if lookahead is reached 2 bitcoin confirmations span is 20 minutes on average.
			// We look ahead up to 100 pending cctx to target TPM of 5.
			b.outboundLogger(outboundID).Warn().
				Uint64("outbound.earliest_pending_nonce", cctxList[0].GetCurrentOutboundParam().TssNonce).
				Msg("Schedule CCTX: lookahead reached")
			return nil
		case b.signer.IsOutboundActive(outboundID):
			// outbound is already being processed
			continue
		}

		// schedule TSS keysign if retry interval has arrived
		if nonce%scheduleInterval == zetaHeight%scheduleInterval {
			go b.signer.TryProcessOutbound(
				ctx,
				cctx,
				b.observer,
				zetaHeight,
			)
		}
	}

	return nil
}

func (b *Bitcoin) updateChainParams(ctx context.Context) error {
	// no changes for signer

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(b.observer.Chain().ChainId)
	if err != nil {
		return err
	}

	b.observer.SetChainParams(*chain.Params())

	return nil
}

func (b *Bitcoin) outboundLogger(id string) *zerolog.Logger {
	l := b.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
