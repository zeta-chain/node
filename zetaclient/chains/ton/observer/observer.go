package observer

import (
	"context"

	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
)

type Observer struct {
	base.Observer

	client    *liteapi.Client
	gatewayID ton.AccountID
}

var _ interfaces.ChainObserver = (*Observer)(nil)

func New(bo *base.Observer, client *liteapi.Client, gatewayID ton.AccountID) (*Observer, error) {
	bo.LoadLastTxScanned()

	return &Observer{
		Observer:  *bo,
		gatewayID: gatewayID,
		client:    client,
	}, nil
}

func (ob *Observer) Start(ctx context.Context) {
	if ok := ob.Observer.Start(); !ok {
		ob.Logger().Chain.Info().Msgf("observer is already started for chain %d", ob.Chain().ChainId)
		return
	}

	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	// Note that each `watch*` method has a ticker that will stop as soon as
	// baseObserver.Stop() was called (ticker.WithStopChan)

	// watch for incoming txs and post votes to zetacore
	bg.Work(ctx, ob.watchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))

	// todo
	//  watchOutbound
	//  watchGasPrice
	//  watchInboundTracker
	//  watchOutbound
	//  watchRPCStatus
}

func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *types.CrossChainTx) (bool, error) {
	// todo
	return false, nil
}
