package observer

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/bg"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zetaton "github.com/zeta-chain/node/zetaclient/chains/ton"
	"github.com/zeta-chain/node/zetaclient/common"
)

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	client  LiteClient
	gateway *toncontracts.Gateway

	outbounds *lru.Cache
}

var _ interfaces.ChainObserver = (*Observer)(nil)

const outboundsCacheSize = 1024

// LiteClient represents a TON client
// see https://github.com/ton-blockchain/ton/blob/master/tl/generate/scheme/tonlib_api.tl
//
//go:generate mockery --name LiteClient --filename ton_liteclient.go --case underscore --output ../../../testutils/mocks
type LiteClient interface {
	zetaton.ConfigGetter
	GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error)
	GetBlockHeader(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (tlb.BlockInfo, error)
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) ([]ton.Transaction, error)
	GetFirstTransaction(ctx context.Context, acc ton.AccountID) (*ton.Transaction, int, error)
	GetTransaction(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) (*ton.Transaction, error)
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

	outbounds, err := lru.New(outboundsCacheSize)
	if err != nil {
		return nil, err
	}

	bo.LoadLastTxScanned()

	return &Observer{
		Observer:  bo,
		client:    client,
		gateway:   gateway,
		outbounds: outbounds,
	}, nil
}

// Start starts the observer. This method is NOT blocking.
// Note that each `watch*` method has a ticker that will stop as soon as
// baseObserver.Stop() was called (ticker.WithStopChan)
func (ob *Observer) Start(ctx context.Context) {
	if ok := ob.Observer.Start(); !ok {
		ob.Logger().Chain.Info().Msg("observer is already started")
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

// watchGasPrice observes TON gas price and votes it to Zetacore.
func (ob *Observer) watchGasPrice(ctx context.Context) error {
	task := func(ctx context.Context, t *ticker.Ticker) error {
		if err := ob.postGasPrice(ctx); err != nil {
			ob.Logger().GasPrice.Err(err).Msg("WatchGasPrice: postGasPrice error")
		}

		newInterval := ticker.SecondsFromUint64(ob.ChainParams().GasPriceTicker)
		t.SetInterval(newInterval)

		return nil
	}

	ob.Logger().GasPrice.Info().Msg("WatchGasPrice started")

	return ticker.Run(
		ctx,
		ticker.SecondsFromUint64(ob.ChainParams().GasPriceTicker),
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().GasPrice, "WatchGasPrice"),
	)
}

// postGasPrice fetches on-chain gas config and reports it to Zetacore.
func (ob *Observer) postGasPrice(ctx context.Context) error {
	cfg, err := zetaton.FetchGasConfig(ctx, ob.client)
	if err != nil {
		return errors.Wrap(err, "failed to fetch gas config")
	}

	gasPrice, err := zetaton.ParseGasPrice(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to parse gas price")
	}

	blockID, err := ob.getLatestMasterchainBlock(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get latest masterchain block")
	}

	// There's no concept of priority fee in TON
	const priorityFee = 0

	_, errVote := ob.
		ZetacoreClient().
		PostVoteGasPrice(ctx, ob.Chain(), gasPrice, priorityFee, uint64(blockID.Seqno))

	return errVote
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
		ticker.WithLogger(ob.Logger().Chain, "WatchRPCStatus"),
	)
}

// checkRPCStatus checks TON RPC status and alerts if necessary.
func (ob *Observer) checkRPCStatus(ctx context.Context) error {
	blockID, err := ob.getLatestMasterchainBlock(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get latest masterchain block")
	}

	block, err := ob.client.GetBlockHeader(ctx, blockID, 0)
	if err != nil {
		return errors.Wrap(err, "failed to get masterchain block header")
	}

	if block.NotMaster {
		return errors.Errorf("block %q is not a master block", blockID.BlockID.String())
	}

	blockTime := time.Unix(int64(block.GenUtime), 0).UTC()

	// will be overridden by chain config
	const defaultAlertLatency = 30 * time.Second

	ob.AlertOnRPCLatency(blockTime, defaultAlertLatency)

	return nil
}

func (ob *Observer) getLatestMasterchainBlock(ctx context.Context) (ton.BlockIDExt, error) {
	mc, err := ob.client.GetMasterchainInfo(ctx)
	if err != nil {
		return ton.BlockIDExt{}, errors.Wrap(err, "failed to get masterchain info")
	}

	return mc.Last.ToBlockIdExt(), nil
}
