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
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

type Bitcoin struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
	proc      *outboundprocessor.Processor
}

func New(
	scheduler *scheduler.Scheduler,
	observer *observer.Observer,
	signer *signer.Signer,
) *Bitcoin {
	// TODO move this to base signer
	// https://github.com/zeta-chain/node/issues/3330
	proc := outboundprocessor.NewProcessor(observer.Logger().Outbound)

	return &Bitcoin{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
		proc:      proc,
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

	// TODO: should we share & fan-out the same chan across all chains?
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
	register(b.observer.ObserveInboundTrackers, "observe_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(b.observer.FetchUTXOs, "fetch_utxos", optUTXOInterval, optGenericSkipper)
	register(b.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(b.observer.CheckRPCStatus, "check_rpc_status")
	register(b.observer.ObserveOutbound, "observe_outbound", optOutboundInterval, optOutboundSkipper)

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
	var (
		lookahead = b.observer.ChainParams().OutboundScheduleLookahead
		chainID   = b.observer.Chain().ChainId
	)

	if err := b.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	zetaBlock, ok := scheduler.BlockFromContext(ctx)
	if !ok {
		return errors.New("unable to get zeta block from context")
	}

	// #nosec G115 always in range
	zetaHeight := uint64(zetaBlock.Block.Height)

	cctxList, _, err := b.observer.ZetacoreClient().ListPendingCCTX(ctx, chainID)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		var (
			params     = cctx.GetCurrentOutboundParam()
			nonce      = params.TssNonce
			outboundID = outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)
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
		case !b.proc.IsOutboundActive(outboundID):
			// outbound is already being processed
			continue
		}

		b.proc.StartTryProcess(outboundID)

		go b.signer.TryProcessOutbound(
			ctx,
			cctx,
			b.proc,
			outboundID,
			b.observer,
			b.observer.ZetacoreClient(),
			zetaHeight,
		)
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
