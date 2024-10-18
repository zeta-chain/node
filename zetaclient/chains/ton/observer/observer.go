package observer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/bg"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/common"
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
	GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error)
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
// Note that each `watch*` method has a ticker that will stop as soon as
// baseObserver.Stop() was called (ticker.WithStopChan)
func (ob *Observer) Start(ctx context.Context) {
	if ok := ob.Observer.Start(); !ok {
		ob.Logger().Chain.Info().Msgf("observer is already started")
		return
	}

	ob.Logger().Chain.Info().Msg("observer is starting")

	start := func(job func(ctx context.Context) error, name string, log zerolog.Logger) {
		bg.Work(ctx, job, bg.WithName(name), bg.WithLogger(log))
	}

	// TODO: watchInboundTracker
	// https://github.com/zeta-chain/node/issues/2935

	start(ob.watchInbound, "WatchInbound", ob.Logger().Inbound)
	start(ob.watchOutbound, "WatchOutbound", ob.Logger().Outbound)
	start(ob.watchGasPrice, "WatchGasPrice", ob.Logger().GasPrice)
	start(ob.watchRPCStatus, "WatchRPCStatus", ob.Logger().Chain)
}

func (ob *Observer) VoteOutboundIfConfirmed(_ context.Context, _ *types.CrossChainTx) (bool, error) {
	return false, errors.New("not implemented")
}

// watchGasPrice observes TON gas price and votes it to Zetacore.
func (ob *Observer) watchGasPrice(_ context.Context) error {
	// todo implement me
	return nil
}

// watchRPCStatus observes TON RPC status.
func (ob *Observer) watchRPCStatus(ctx context.Context) error {
	task := func(ctx context.Context, _ *ticker.Ticker) error {
		if err := ob.checkRPCStatus(ctx); err != nil {
			ob.Logger().Chain.Err(err).Msg("checkRPCStatus error")
		}

		return nil
	}

	return ticker.Run(
		ctx,
		common.RPCStatusCheckInterval,
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().Inbound, "WatchRPCStatus"),
	)
}

// checkRPCStatus checks TON RPC status and alerts if necessary.
func (ob *Observer) checkRPCStatus(ctx context.Context) error {
	mc, err := ob.client.GetMasterchainInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get masterchain info")
	}

	blockID := mc.Last

	block, err := ob.client.GetBlockHeader(ctx, blockID.ToBlockIdExt(), 0)
	if err != nil {
		return errors.Wrap(err, "failed to get masterchain block header")
	}

	if block.NotMaster {
		return errors.Errorf("block %q is not a master block", blockID.ToBlockIdExt().BlockID.String())
	}

	blockTime := time.Unix(int64(block.GenUtime), 0).UTC()

	// will be overridden by chain config
	const defaultAlertLatency = 30 * time.Second

	ob.AlertOnRPCLatency(blockTime, defaultAlertLatency)

	return nil
}
