package ton

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

// TON represents TON observerSigner.
type TON struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
	proc      *outboundprocessor.Processor
}

// New TON constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *TON {
	return &TON{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
		proc:      outboundprocessor.NewProcessor(observer.Logger().Outbound),
	}
}

func (t *TON) Chain() chains.Chain {
	return t.observer.Chain()
}

func (t *TON) Start(ctx context.Context) error {
	if ok := t.observer.Observer.Start(); !ok {
		return errors.Errorf("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := t.observer.ZetacoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !t.observer.ChainParams().IsSupported
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(t.group()),
			scheduler.Name(name),
		}, opts...)

		t.scheduler.Register(ctx, exec, opts...)
	}

	register(t.observer.ObserveGateway, "observe_gateway", optInboundInterval, optInboundSkipper)
	register(t.observer.ObserveInboundTrackers, "observe_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(t.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(t.observer.CheckRPCStatus, "check_rpc_status")
	register(t.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(t.observer.ObserveOutbound, "observe_outbound", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(t.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

func (t *TON) Stop() {
	t.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	t.scheduler.StopGroup(t.group())
}

func (t *TON) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("ton:%d", t.observer.Chain().ChainId),
	)
}

func (t *TON) scheduleCCTX(ctx context.Context) error {
	if err := t.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	chainID := t.observer.Chain().ChainId

	zetaBlock, ok := scheduler.BlockFromContext(ctx)
	if !ok {
		return errors.New("unable to get zeta block from context")
	}

	// #nosec G115 always in range
	zetaHeight := uint64(zetaBlock.Block.Height)

	cctxList, _, err := t.observer.ZetacoreClient().ListPendingCCTX(ctx, chainID)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	for i := range cctxList {
		var (
			cctx       = cctxList[i]
			params     = cctx.GetCurrentOutboundParam()
			nonce      = params.TssNonce
			outboundID = outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)
		)

		if params.ReceiverChainId != chainID {
			t.outboundLogger(outboundID).Error().Msg("Schedule CCTX: chain id mismatch")
			continue
		}

		// vote outbound if it's already confirmed
		continueKeySign, err := t.observer.VoteOutboundIfConfirmed(ctx, cctx)
		switch {
		case err != nil:
			t.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX: VoteOutboundIfConfirmed failed")
			continue
		case !continueKeySign:
			t.outboundLogger(outboundID).Info().Msg("Schedule CCTX: outbound already processed")
			continue
		case t.proc.IsOutboundActive(outboundID):
			// outbound is already being processed
			continue
		}

		go t.signer.TryProcessOutbound(
			ctx,
			cctx,
			t.proc,
			outboundID,
			t.observer.ZetacoreClient(),
			zetaHeight,
		)
	}

	return nil
}

func (t *TON) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(t.observer.Chain().ChainId)
	if err != nil {
		return err
	}

	t.signer.SetGatewayAddress(chain.Params().GatewayAddress)
	t.observer.SetChainParams(*chain.Params())

	return nil
}

func (t *TON) outboundLogger(id string) *zerolog.Logger {
	l := t.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
