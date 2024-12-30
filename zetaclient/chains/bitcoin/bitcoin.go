package bitcoin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
)

type Bitcoin struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

func New(
	scheduler *scheduler.Scheduler,
	observer *observer.Observer,
	signer *signer.Signer,
) *Bitcoin {
	return &Bitcoin{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

func (b *Bitcoin) Chain() chains.Chain {
	return b.observer.Chain()
}

func (b *Bitcoin) Start(_ context.Context) error {
	if ok := b.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	// 	// watch bitcoin chain for incoming txs and post votes to zetacore
	//	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))
	//
	//	// watch bitcoin chain for outgoing txs status
	//	bg.Work(ctx, ob.WatchOutbound, bg.WithName("WatchOutbound"), bg.WithLogger(ob.Logger().Outbound))
	//
	//	// watch bitcoin chain for UTXOs owned by the TSS address
	//	bg.Work(ctx, ob.WatchUTXOs, bg.WithName("WatchUTXOs"), bg.WithLogger(ob.Logger().Outbound))
	//
	//	// watch bitcoin chain for gas rate and post to zetacore
	//	bg.Work(ctx, ob.WatchGasPrice, bg.WithName("WatchGasPrice"), bg.WithLogger(ob.Logger().GasPrice))
	//
	//	// watch zetacore for bitcoin inbound trackers
	//	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))
	//
	//	// watch the RPC status of the bitcoin chain
	//	bg.Work(ctx, ob.watchRPCStatus, bg.WithName("watchRPCStatus"), bg.WithLogger(ob.Logger().Chain))

	// todo start & schedule
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
