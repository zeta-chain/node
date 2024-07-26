package context

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observer "github.com/zeta-chain/zetacore/x/observer/types"
)

// ChainRegistry is a registry of supported chains
type ChainRegistry struct {
	chains map[int64]Chain

	// additionalChains is a list of additional static chain information to use when searching from
	// chain IDs. It's stored in the protocol to dynamically support new chains without doing an upgrade
	additionalChains []chains.Chain

	mu sync.Mutex
}

// Chain represents chain with its parameters
type Chain struct {
	id       int64
	chain    *chains.Chain
	params   *observer.ChainParams
	registry *ChainRegistry
}

var (
	ErrChainNotFound     = errors.New("chain not found")
	ErrChainNotSupported = errors.New("chain not supported")
)

// NewChainRegistry constructs a new ChainRegistry
func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{
		chains:           make(map[int64]Chain),
		additionalChains: []chains.Chain{},
		mu:               sync.Mutex{},
	}
}

func (cr *ChainRegistry) Get(chainID int64) (Chain, error) {
	chain, ok := cr.chains[chainID]
	if !ok {
		return Chain{}, ErrChainNotFound
	}

	return chain, nil
}

// Set sets a chain in the registry.
// A chain must be SUPPORTED and NOT ZetaChain itself; otherwise returns ErrChainNotSupported
func (cr *ChainRegistry) Set(chainID int64, chain *chains.Chain, params *observer.ChainParams) error {
	item, err := newChain(cr, chainID, chain, params)
	if err != nil {
		return err
	}

	item.registry = cr

	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.chains[item.id] = item

	return nil
}

func (cr *ChainRegistry) SetAdditionalChains(chains []chains.Chain) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.additionalChains = chains
}

// Delete deletes one or more chains from the registry
func (cr *ChainRegistry) Delete(chainIDs ...int64) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	for _, id := range chainIDs {
		delete(cr.chains, id)
	}
}

// Has checks if the chain is in the registry
func (cr *ChainRegistry) Has(chainID int64) bool {
	_, ok := cr.chains[chainID]
	return ok
}

// ChainIDs returns a list of chain IDs in the registry
func (cr *ChainRegistry) ChainIDs() []int64 {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	return maps.Keys(cr.chains)
}

func newChain(cr *ChainRegistry, chainID int64, chain *chains.Chain, params *observer.ChainParams) (Chain, error) {
	switch {
	case chainID < 1:
		return Chain{}, fmt.Errorf("invalid chain id %d", chainID)
	case chain == nil:
		return Chain{}, fmt.Errorf("chain is nil")
	case params == nil:
		return Chain{}, fmt.Errorf("chain params is nil")
	case chain.ChainId != chainID:
		return Chain{}, fmt.Errorf("chain id %d does not match chain.ChainId %d", chainID, chain.ChainId)
	case params.ChainId != chainID:
		return Chain{}, fmt.Errorf("chain id %d does not match params.ChainId %d", chainID, params.ChainId)
	case !params.IsSupported:
		return Chain{}, ErrChainNotSupported
	case chains.IsZetaChain(chainID, nil) || !chain.IsExternal:
		return Chain{}, errors.Wrap(ErrChainNotSupported, "ZetaChain itself cannot be in the registry")
	}

	return Chain{
		id:       chainID,
		chain:    chain,
		params:   params,
		registry: cr,
	}, nil
}

func (c Chain) Params() *observer.ChainParams {
	return c.params
}

func (c Chain) IsEVM() bool {
	return chains.IsEVMChain(c.id, c.registry.additionalChains)
}

func (c Chain) IsUTXO() bool {
	return chains.IsBitcoinChain(c.id, c.registry.additionalChains)
}

func (c Chain) IsSolana() bool {
	return chains.IsSolanaChain(c.id, c.registry.additionalChains)
}
