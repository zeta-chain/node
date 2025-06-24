package observer

import (
	"context"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	rpc     RPC
	gateway *toncontracts.Gateway

	outbounds *lru.Cache

	latestGasPrice atomic.Uint64
}

const outboundsCacheSize = 1024

type RPC interface {
	GetConfigParam(ctx context.Context, index uint32) (*boc.Cell, error)
	GetBlockHeader(ctx context.Context, blockID rpc.BlockIDExt) (rpc.BlockHeader, error)
	GetMasterchainInfo(ctx context.Context) (rpc.MasterchainInfo, error)
	HealthCheck(ctx context.Context) (time.Time, error)
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) ([]ton.Transaction, error)
	GetTransactions(
		ctx context.Context,
		count uint32,
		accountID ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)
	GetTransaction(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) (ton.Transaction, error)
}

// New constructor for TON Observer.
func New(bo *base.Observer, rpc RPC, gateway *toncontracts.Gateway) (*Observer, error) {
	switch {
	case !bo.Chain().IsTONChain():
		return nil, errors.New("base observer chain is not TON")
	case rpc == nil:
		return nil, errors.New("rpc is nil")
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
		rpc:       rpc,
		gateway:   gateway,
		outbounds: outbounds,
	}, nil
}

// PostGasPrice fetches on-chain gas config and reports it to Zetacore.
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	cfg, err := rpc.FetchGasConfigRPC(ctx, ob.rpc)
	if err != nil {
		return errors.Wrap(err, "failed to fetch gas config")
	}

	gasPrice, err := rpc.ParseGasPrice(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to parse gas price")
	}

	info, err := ob.rpc.GetMasterchainInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get masterchain info")
	}

	blockNum := uint64(info.Last.Seqno)

	// There's no concept of priority fee in TON
	const priorityFee = 0

	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), gasPrice, priorityFee, blockNum)
	if err != nil {
		return errors.Wrap(err, "failed to post gas price")
	}

	ob.setLatestGasPrice(gasPrice)

	return nil
}

// CheckRPCStatus checks TON RPC status and alerts if necessary.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.rpc.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}

func (ob *Observer) getLatestGasPrice() (uint64, error) {
	price := ob.latestGasPrice.Load()

	if price > 0 {
		return price, nil
	}

	return 0, errors.New("latest gas price is not set")
}

func (ob *Observer) setLatestGasPrice(price uint64) {
	ob.latestGasPrice.Store(price)
}
