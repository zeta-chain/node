package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// Solana represents Solana observer-signer.
type Solana struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

// New Solana constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *Solana {
	return &Solana{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

// Chain returns chain
func (s *Solana) Chain() chains.Chain {
	return s.observer.Chain()
}

// Start starts observer-signer for
// processing inbound & outbound cross-chain transactions.
func (s *Solana) Start(ctx context.Context) error {
	if ok := s.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := s.observer.ZetaRepo().WatchNewBlocks(ctx)
	if err != nil {
		return err
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipInbound(s.observer.Observer, app) })
	optOutboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipOutbound(s.observer.Observer, app) })
	optGasPriceSkipper := scheduler.Skipper(func() bool { return base.CheckSkipGasPrice(s.observer.Observer, app) })

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(s.group()),
			scheduler.Name(name),
		}, opts...)

		s.scheduler.Register(ctx, exec, opts...)
	}

	register(s.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(s.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(s.observer.ProcessInternalTrackers, "process_internal_trackers", optInboundInterval, optInboundSkipper)
	register(s.observer.PostGasPrice, "post_gas_price", optGasInterval, optGasPriceSkipper)
	register(s.observer.CheckRPCStatus, "check_rpc_status")
	register(s.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX scheduler (every zetachain block)
	register(s.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops all relevant tasks.
func (s *Solana) Stop() {
	s.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	s.scheduler.StopGroup(s.group())
}

func (s *Solana) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("sol:%d", s.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules solana outbound keysign
func (s *Solana) scheduleCCTX(ctx context.Context) error {
	if err := s.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta block from context")
	}

	time.Sleep(delay)

	var (
		chain   = s.observer.Chain()
		chainID = chain.ChainId

		// #nosec G115 positive
		zetaHeight = uint64(zetaBlock.Block.Height)

		// #nosec G115 positive
		interval           = uint64(s.observer.ChainParams().OutboundScheduleInterval)
		scheduleLookahead  = s.observer.ChainParams().OutboundScheduleLookahead
		scheduleLookback   = uint64(float64(scheduleLookahead) * constant.OutboundLookbackFactor)
		needsProcessingCtr = 0
	)

	cctxList, err := s.observer.ZetaRepo().GetPendingCCTXs(ctx)
	if err != nil {
		return err
	}

	// schedule keysign for each pending cctx
	for i, cctx := range cctxList {
		var (
			params        = cctx.GetCurrentOutboundParam()
			inboundParams = cctx.GetInboundParams()
			nonce         = params.TssNonce
			outboundID    = base.OutboundIDFromCCTX(cctx)
		)

		logger := s.observer.Logger().Outbound.With().Str(logs.FieldOutboundID, outboundID).Logger()

		switch {
		case int64(i) == scheduleLookahead:
			// stop if lookahead is reached
			return nil
		case params.ReceiverChainId != chainID:
			logger.Error().Msg("chain id mismatch")
			continue
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+scheduleLookback:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// schedule newly created cctx right away, no need to wait for next interval
		// 1. schedule the very first cctx (there can be multiple) created in the last Zeta block.
		// 2. schedule new cctx only when there is no other older cctx to process
		isCCTXNewlyCreated := inboundParams.ObservedExternalHeight == zetaHeight
		shouldProcessCCTXImmedately := isCCTXNewlyCreated && needsProcessingCtr == 0

		// even if the outbound is currently active, we should increment this counter
		// to avoid immediate processing logic
		needsProcessingCtr++

		if s.signer.IsOutboundActive(outboundID) {
			continue
		}

		// vote outbound if it's already confirmed
		continueKeysign, err := s.observer.VoteOutboundIfConfirmed(ctx, cctx)
		switch {
		case err != nil:
			logger.Error().Err(err).Msg("schedule CCTX: error calling VoteOutboundIfConfirmed")
			continue
		case !continueKeysign:
			logger.Info().Msg("schedule CCTX: outbound already processed")
			continue
		}

		shouldScheduleProcess := nonce%interval == zetaHeight%interval

		// schedule a TSS keysign
		if shouldProcessCCTXImmedately || shouldScheduleProcess {
			go s.signer.TryProcessOutbound(
				ctx,
				cctx,
				s.observer.ZetaRepo(),
				zetaHeight,
			)
		}
	}

	return nil
}

func (s *Solana) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(s.observer.Chain().ChainId)
	if err != nil {
		return err
	}

	params := chain.Params()

	s.observer.SetChainParams(*params)
	s.signer.SetGatewayAddress(params.GatewayAddress)

	return nil
}
