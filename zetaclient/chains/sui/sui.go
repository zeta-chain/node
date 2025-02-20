package sui

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
	"github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	"github.com/zeta-chain/node/zetaclient/chains/sui/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

// Sui observer-signer.
type Sui struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

// New Sui observer-signer constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *Sui {
	return &Sui{scheduler, observer, signer}
}

// Chain returns chain
func (s *Sui) Chain() chains.Chain {
	return s.observer.Chain()
}

// Start starts observer-signer for processing inbound & outbound cross-chain transactions.
func (s *Sui) Start(ctx context.Context) error {
	if ok := s.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := s.observer.ZetacoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	// todo
	//   - [ ] ObserveInbound
	//   - [ ] ProcessInboundTrackers
	//   - [ ] ProcessOutboundTrackers
	//   - [ ] ScheduleCCTX

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(s.group()),
			scheduler.Name(name),
		}, opts...)

		s.scheduler.Register(ctx, exec, opts...)
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().GasPriceTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !s.observer.ChainParams().IsSupported
	})

	register(s.observer.ObserveInbound, "observer_inbound", optInboundInterval, optInboundSkipper)
	register(s.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(s.observer.CheckRPCStatus, "check_rpc_status")
	register(s.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)

	// CCTX scheduler (every zetachain block)
	register(s.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops all relevant tasks.
func (s *Sui) Stop() {
	s.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	s.scheduler.StopGroup(s.group())
}

func (s *Sui) group() scheduler.Group {
	return scheduler.Group(fmt.Sprintf("sui:%d", s.Chain().ChainId))
}

// scheduleCCTX schedules outbound cross-chain transactions.
func (s *Sui) scheduleCCTX(ctx context.Context) error {
	if err := s.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta block from context")
	}

	time.Sleep(delay)

	// #nosec G115 always in range
	zetaHeight := uint64(zetaBlock.Block.Height)
	chain := s.observer.Chain()

	cctxList, _, err := s.observer.ZetacoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	for i := range cctxList {
		cctx := cctxList[i]

		if err := s.signer.ProcessCCTX(ctx, cctx, zetaHeight); err != nil {
			outboundID := base.OutboundIDFromCCTX(cctx)
			s.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX failed")
		}
	}

	return nil
}

func (s *Sui) updateChainParams(ctx context.Context) error {
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

	// note that address should be in format of `$packageID,$gatewayObjectID`
	if err := s.observer.Gateway().UpdateIDs(params.GatewayAddress); err != nil {
		return errors.Wrap(err, "unable to update gateway ids")
	}

	return nil
}

func (s *Sui) outboundLogger(id string) *zerolog.Logger {
	l := s.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
