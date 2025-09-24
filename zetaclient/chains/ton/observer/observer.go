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
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

const outboundsCacheSize = 1024

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	tonRepo  *repo.TONRepo
	zetaRepo *zrepo.ZetaRepo

	gateway *toncontracts.Gateway

	outbounds *lru.Cache // indexed by nonce

	latestGasPrice atomic.Uint64
}

// New constructs a TON Observer.
func New(baseObserver *base.Observer,
	tonClient TONClient,
	gateway *toncontracts.Gateway,
) (*Observer, error) {
	if !baseObserver.Chain().IsTONChain() {
		return nil, errors.New("invalid chain (not TON)")
	}
	if tonClient == nil {
		return nil, errors.New("invalid TON client")
	}
	if gateway == nil {
		return nil, errors.New("invalid gateway")
	}

	outbounds, err := lru.New(outboundsCacheSize)
	if err != nil {
		return nil, err
	}

	baseObserver.LoadLastTxScanned()

	chain := baseObserver.Chain()
	return &Observer{
		Observer:  baseObserver,
		tonRepo:   repo.NewTONRepo(tonClient, gateway, chain),
		zetaRepo:  zrepo.New(baseObserver.ZetacoreClient(), chain),
		gateway:   gateway,
		outbounds: outbounds,
	}, nil
}

// getOutbound returns an outbound from the in-memory cache given a nonce.
// It returns nil if the cache does not contain an outbound associated with that nonce.
func (ob *Observer) getOutbound(nonce uint64) *Outbound {
	v, ok := ob.outbounds.Get(nonce)
	if !ok {
		return nil
	}
	outbound := v.(Outbound)
	return &outbound
}

// addOutbound adds an outbound to the in-memory cache indexed by nonce.
func (ob *Observer) addOutbound(outbound Outbound) {
	ob.outbounds.Add(outbound.nonce, outbound)
}

// getLatestGasPrice atomically retrieves the latest gas price stored in the memory.
func (ob *Observer) getLatestGasPrice() (uint64, error) {
	price := ob.latestGasPrice.Load()
	if price > 0 {
		return price, nil
	}

	return 0, errors.New("latest gas price is not set")
}

// getLatestGasPrice atomically sets the latest gas price.
func (ob *Observer) setLatestGasPrice(price uint64) {
	ob.latestGasPrice.Store(price)
}

// ------------------------------------------------------------------------------------------------
// TODO
// ------------------------------------------------------------------------------------------------

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

// CheckRPCStatus checks the status of the TON client and reports the block latency.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.tonRepo.CheckHealth(ctx)
	if err != nil {
		return err
	}

	metrics.ReportBlockLatency(ob.Chain().Name, *blockTime)
	return nil
}
