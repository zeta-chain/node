package ton

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
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

	newBlockChan, err := t.observer.ZetaRepo().WatchNewBlocks(ctx)
	if err != nil {
		return err
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

	optInboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipInbound(t.observer.Observer, app) })
	optOutboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipOutbound(t.observer.Observer, app) })
	optGasPriceSkipper := scheduler.Skipper(func() bool { return base.CheckSkipGasPrice(t.observer.Observer, app) })

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(t.group()),
			scheduler.Name(name),
		}, opts...)

		t.scheduler.Register(ctx, exec, opts...)
	}

	register(t.observer.CheckRPCStatus, "check_rpc_status")
	register(t.observer.ObserveGasPrice, "observe_gas_price", optGasInterval, optGasPriceSkipper)
	register(t.observer.ObserveInbounds, "observe_inbounds", optInboundInterval, optInboundSkipper)
	register(t.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(t.observer.ProcessInternalTrackers, "process_internal_trackers", optInboundInterval, optInboundSkipper)
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
	zetaRepo := t.observer.ZetaRepo()

	// skip stale block event in channel if any
	blockHeight, stale, err := t.signer.CheckBlockEvent(ctx, zetaRepo)
	if err != nil {
		return errors.Wrap(err, "unable to check stale block event")
	} else if stale {
		return nil
	}

	if err := t.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "failed to update chain parameters")
	}

	cctxs, err := zetaRepo.GetPendingCCTXs(ctx)
	if err != nil {
		return err
	}

	// no-op
	if len(cctxs) == 0 {
		return nil
	}

	var (
		// #nosec G115 always in range
		zetaHeight = uint64(blockHeight)
		chainID    = t.observer.Chain().ChainId

		// #nosec G115 positive
		interval  = uint64(t.observer.ChainParams().OutboundScheduleInterval)
		lookahead = t.observer.ChainParams().OutboundScheduleLookahead
		// #nosec G115 always in range
		maxNonceOffset = uint64(float64(lookahead) * constant.MaxNonceOffsetFactor)

		firstNonce = cctxs[0].GetCurrentOutboundParam().TssNonce
		maxNonce   = firstNonce + maxNonceOffset
	)

	for i, cctx := range cctxs {
		var (
			outboundID     = base.OutboundIDFromCCTX(cctx)
			outboundParams = cctx.GetCurrentOutboundParam()
			nonce          = outboundParams.TssNonce
			logger         = t.outboundLogger(outboundID)
		)

		switch {
		case int64(i) == lookahead:
			// stop if lookahead is reached
			return nil
		case outboundParams.ReceiverChainId != chainID:
			// should not happen
			logger.Error().Msg("chain id mismatch")
			continue
		case nonce > maxNonce:
			return fmt.Errorf("nonce %d is too high (%s). Earliest nonce %d", nonce, outboundID, firstNonce)
		case t.signer.IsOutboundActive(outboundID):
			// cctx is already being processed & broadcasted by signer
			continue
		}

		// vote outbound if it's already confirmed
		continueKeysign, err := t.observer.VoteOutboundIfConfirmed(ctx, cctx)
		if err != nil {
			logger.Error().Err(err).Msg("call to VoteOutboundIfConfirmed failed")
			continue
		}
		if !continueKeysign {
			logger.Info().Msg("outbound already processed")
			continue
		}

		// schedule keysign if the interval has arrived
		if nonce%interval == zetaHeight%interval {
			go t.signer.TryProcessOutbound(ctx, cctx, t.observer.ZetaRepo(), zetaHeight)
		}
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
	l := t.observer.Logger().Outbound.With().Str(logs.FieldOutboundID, id).Logger()

	return &l
}
