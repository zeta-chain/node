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
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// TON represents TON observer-signer components that is responsible
// for processing and scheduling inbound and outbound TON transactions.
type TON struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

// New TON constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *TON {
	return &TON{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

// Chain returns the chain struct
func (t *TON) Chain() chains.Chain {
	return t.observer.Chain()
}

// Start starts the observer-signer and schedules various regular background tasks e.g. inbound observation.
func (t *TON) Start(ctx context.Context) error {
	if ok := t.observer.Observer.Start(); !ok {
		return errors.Errorf("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get app from context")
	}

	newBlockChan, err := t.observer.ZetacoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create new block subscriber")
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

	register(t.observer.CheckRPCStatus, "check_rpc_status")
	register(t.observer.ObserveGasPrice, "observe_gas_price", optGasInterval, optGenericSkipper)
	register(t.observer.ObserveInbounds, "observe_inbounds", optInboundInterval, optInboundSkipper)
	register(t.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(t.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(t.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops the observer-signer.
func (t *TON) Stop() {
	t.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	t.scheduler.StopGroup(t.group())
}

func (t *TON) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("ton:%d", t.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules cross-chain tx processing.
// It loads pending cctx from zetacore, then tries to sign and broadcast them.
func (t *TON) scheduleCCTX(ctx context.Context) error {
	if err := t.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "failed to update chain parameters")
	}

	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get zeta block from context")
	}

	time.Sleep(delay)

	// #nosec G115 always in range
	zetaHeight := uint64(zetaBlock.Block.Height)
	chain := t.observer.Chain()

	cctxList, _, err := t.observer.ZetacoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "failed to list pending CCTXs")
	}

	for i := range cctxList {
		cctx := cctxList[i]
		outboundID := base.OutboundIDFromCCTX(cctx)

		err := t.processCCTX(ctx, outboundID, cctx, zetaHeight)
		if err != nil {
			t.outboundLogger(outboundID).Error().Err(err).Msg("failed to schedule CCTX")
		}
	}

	return nil
}

func (t *TON) processCCTX(ctx context.Context,
	outboundID string,
	cctx *types.CrossChainTx,
	zetaHeight uint64,
) error {
	switch {
	case t.signer.IsOutboundActive(outboundID):
		return nil //no-op
	case cctx.GetCurrentOutboundParam().ReceiverChainId != t.observer.Chain().ChainId:
		return errors.New("chain id mismatch")
	}

	// vote outbound if it's already confirmed
	continueKeySign, err := t.observer.VoteOutboundIfConfirmed(ctx, cctx)
	if err != nil {
		return errors.Wrap(err, "failed to VoteOutboundIfConfirmed")
	}
	if !continueKeySign {
		t.outboundLogger(outboundID).Info().Msg("schedule CCTX: outbound already processed")
		return nil
	}

	go t.signer.TryProcessOutbound(
		ctx,
		cctx,
		t.observer.ZetacoreClient(),
		zetaHeight,
	)

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
	l := t.observer.Logger().Outbound.With().Str(logs.FieldOutboundID, id).Logger()

	return &l
}
