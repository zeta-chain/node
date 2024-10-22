package observer

import (
	"context"
	"errors"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/bg"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	client  LiteClient
	gateway *toncontracts.Gateway
}

// LiteClient represents a TON client
//
//go:generate mockery --name LiteClient --filename ton_liteclient.go --case underscore --output ../../../testutils/mocks
type LiteClient interface {
	GetBlockHeader(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (tlb.BlockInfo, error)
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, bits ton.Bits256) ([]ton.Transaction, error)
	GetFirstTransaction(ctx context.Context, id ton.AccountID) (*ton.Transaction, int, error)
}

var _ interfaces.ChainObserver = (*Observer)(nil)

// New constructor for TON Observer.
func New(bo *base.Observer, client LiteClient, gateway *toncontracts.Gateway) (*Observer, error) {
	switch {
	case !bo.Chain().IsTONChain():
		return nil, errors.New("base observer chain is not TON")
	case client == nil:
		return nil, errors.New("liteapi client is nil")
	case gateway == nil:
		return nil, errors.New("gateway is nil")
	}

	bo.LoadLastTxScanned()

	return &Observer{
		Observer: bo,
		client:   client,
		gateway:  gateway,
	}, nil
}

// Start starts the observer. This method is NOT blocking.
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

	// TODO: watchInboundTracker
	// https://github.com/zeta-chain/node/issues/2935

	// TODO: outbounds/withdrawals: (watchOutbound, watchGasPrice, watchRPCStatus)
	// https://github.com/zeta-chain/node/issues/2807
}

func (ob *Observer) VoteOutboundIfConfirmed(_ context.Context, _ *types.CrossChainTx) (bool, error) {
	return false, errors.New("not implemented")
}
