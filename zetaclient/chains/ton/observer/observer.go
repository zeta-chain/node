package observer

import (
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
)

const outboundsCacheSize = 1024

// Observer is a TON observer.
type Observer struct {
	*base.Observer

	tonRepo *repo.TONRepo

	gateway *toncontracts.Gateway // used to parse transactions

	outbounds *lru.Cache // indexed by nonce

	latestGasPrice atomic.Uint64
}

// New constructs a TON Observer.
func New(baseObserver *base.Observer,
	tonRepo *repo.TONRepo,
	gateway *toncontracts.Gateway,
) (*Observer, error) {
	if !baseObserver.Chain().IsTONChain() {
		return nil, errors.New("invalid chain (not TON)")
	}
	if tonRepo == nil {
		return nil, errors.New("invalid TON repository")
	}
	if gateway == nil {
		return nil, errors.New("invalid gateway")
	}

	outbounds, err := lru.New(outboundsCacheSize)
	if err != nil {
		return nil, err
	}

	baseObserver.LoadLastTxScanned()

	return &Observer{
		Observer:  baseObserver,
		tonRepo:   tonRepo,
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

// setLatestGasPrice atomically sets the latest gas price.
func (ob *Observer) setLatestGasPrice(price uint64) {
	ob.latestGasPrice.Store(price)
}
