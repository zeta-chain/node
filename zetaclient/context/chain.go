package context

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/zeta-chain/node/pkg/chains"
	observer "github.com/zeta-chain/node/x/observer/types"
)

// ChainRegistry is a registry of supported chains
type ChainRegistry struct {
	chains map[int64]Chain

	// additionalChains is a list of additional static chain information to use when searching from
	// chain IDs. It's stored in the protocol to dynamically support new chains without doing an upgrade
	additionalChains []chains.Chain

	// relayerKeyPasswords maps network name to relayer key password
	relayerKeyPasswords map[string]string

	mu sync.Mutex
}

// Chain represents chain with its parameters
type Chain struct {
	chainInfo      *chains.Chain
	observerParams *observer.ChainParams

	// reference to the registry it necessary for some operations
	// like checking if the chain is EVM or not because it uses some "global" context state
	registry *ChainRegistry
}

var (
	ErrChainNotFound     = errors.New("chain not found")
	ErrChainNotSupported = errors.New("chain not supported")
)

// NewChainRegistry constructs a new ChainRegistry
func NewChainRegistry(relayerKeyPasswords map[string]string) *ChainRegistry {
	return &ChainRegistry{
		chains:              make(map[int64]Chain),
		additionalChains:    []chains.Chain{},
		relayerKeyPasswords: relayerKeyPasswords,
		mu:                  sync.Mutex{},
	}
}

// Get returns a chain by ID.
func (cr *ChainRegistry) Get(chainID int64) (Chain, error) {
	chain, ok := cr.chains[chainID]
	if !ok {
		return Chain{}, errors.Wrapf(ErrChainNotFound, "id=%d", chainID)
	}

	return chain, nil
}

// All returns all chains in the registry sorted by chain ID.
func (cr *ChainRegistry) All() []Chain {
	items := maps.Values(cr.chains)

	slices.SortFunc(items, func(a, b Chain) bool { return a.ID() < b.ID() })

	return items
}

// Set sets a chain in the registry.
// A chain must be SUPPORTED; otherwise returns ErrChainNotSupported
func (cr *ChainRegistry) Set(chainID int64, chain *chains.Chain, params *observer.ChainParams) error {
	item, err := newChain(cr, chainID, chain, params)
	if err != nil {
		return err
	}

	item.registry = cr

	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.chains[chainID] = item

	return nil
}

// SetAdditionalChains sets additional chains to the registry
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
	if err := validateNewChain(chainID, chain, params); err != nil {
		return Chain{}, errors.Wrap(err, "invalid input")
	}

	return Chain{
		chainInfo:      chain,
		observerParams: params,
		registry:       cr,
	}, nil
}

func (c Chain) ID() int64 {
	return c.chainInfo.ChainId
}

func (c Chain) Name() string {
	return c.chainInfo.Name
}

func (c Chain) Params() *observer.ChainParams {
	return c.observerParams
}

// RawChain returns the underlying Chain object. Better not to use this method
func (c Chain) RawChain() *chains.Chain {
	return c.chainInfo
}

func (c Chain) IsEVM() bool {
	return chains.IsEVMChain(c.ID(), c.registry.additionalChains)
}

func (c Chain) IsZeta() bool {
	return chains.IsZetaChain(c.ID(), c.registry.additionalChains)
}

func (c Chain) IsUTXO() bool {
	return chains.IsBitcoinChain(c.ID(), c.registry.additionalChains)
}

func (c Chain) IsSolana() bool {
	return chains.IsSolanaChain(c.ID(), c.registry.additionalChains)
}

// RelayerKeyPassword returns the relayer key password for the chain
func (c Chain) RelayerKeyPassword() string {
	network := c.RawChain().Network

	return c.registry.relayerKeyPasswords[network.String()]
}

func validateNewChain(chainID int64, chain *chains.Chain, params *observer.ChainParams) error {
	switch {
	case chainID < 1:
		return fmt.Errorf("invalid chain id %d", chainID)
	case chain == nil:
		return fmt.Errorf("chain is nil")
	case params == nil:
		return fmt.Errorf("chain params is nil")
	case chain.ChainId != chainID:
		return fmt.Errorf("chain id %d does not match chain.ChainId %d", chainID, chain.ChainId)
	case params.ChainId != chainID:
		return fmt.Errorf("chain id %d does not match params.ChainId %d", chainID, params.ChainId)
	case !params.IsSupported:
		return ErrChainNotSupported
	}

	return nil
}
