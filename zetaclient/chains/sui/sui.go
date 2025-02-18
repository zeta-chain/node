package sui

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
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
func (s *Sui) scheduleCCTX(_ context.Context) error {
	// todo
	return nil
}
