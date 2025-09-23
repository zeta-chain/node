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
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	tonRepo *repo.TONRepo

	gateway *toncontracts.Gateway

	outbounds *lru.Cache

	latestGasPrice atomic.Uint64
}

const outboundsCacheSize = 1024

type TONClient interface {
	GetConfigParam(_ context.Context, index uint32) (*boc.Cell, error)

	GetBlockHeader(_ context.Context, blockID rpc.BlockIDExt) (rpc.BlockHeader, error)

	GetMasterchainInfo(context.Context) (rpc.MasterchainInfo, error)

	HealthCheck(context.Context) (time.Time, error)

	GetTransaction(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) (ton.Transaction, error)

	GetTransactions(_ context.Context,
		count uint32,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)

	GetTransactionsSince(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)
}

// New constructor for TON Observer.
func New(baseObserver *base.Observer,
	tonClient TONClient,
	gateway *toncontracts.Gateway,
) (*Observer, error) {
	if !baseObserver.Chain().IsTONChain() {
		return nil, errors.New("base observer chain is not TON")
	}
	if tonClient == nil {
		return nil, errors.New("ton client is nil")
	}
	if gateway == nil {
		return nil, errors.New("gateway is nil")
	}

	outbounds, err := lru.New(outboundsCacheSize)
	if err != nil {
		return nil, err
	}

	baseObserver.LoadLastTxScanned()

	return &Observer{
		Observer:  baseObserver,
		tonRepo:   repo.NewTONRepo(tonClient, gateway, baseObserver.Chain()),
		gateway:   gateway,
		outbounds: outbounds,
	}, nil
}

// CheckRPCStatus checks TON RPC status and alerts if necessary.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.tonRepo.Client.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check TON Client health")
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
