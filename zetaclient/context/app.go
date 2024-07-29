// Package context provides global app context for ZetaClient
package context

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AppContext represents application (zetaclient) context.
type AppContext struct {
	config config.Config
	logger zerolog.Logger

	chainRegistry *ChainRegistry

	currentTssPubKey string
	crosschainFlags  observertypes.CrosschainFlags
	keygen           observertypes.Keygen

	mu sync.RWMutex
}

// New creates and returns new empty AppContext
func New(cfg config.Config, logger zerolog.Logger) *AppContext {
	return &AppContext{
		config: cfg,
		logger: logger.With().Str("module", "appcontext").Logger(),

		chainRegistry: NewChainRegistry(),

		crosschainFlags:  observertypes.CrosschainFlags{},
		currentTssPubKey: "",
		keygen:           observertypes.Keygen{},

		mu: sync.RWMutex{},
	}
}

// Config returns the config of the app
func (a *AppContext) Config() config.Config {
	return a.config
}

// GetChain returns the chain by ID.
func (a *AppContext) GetChain(chainID int64) (Chain, error) {
	return a.chainRegistry.Get(chainID)
}

// ListChainIDs returns the list of existing chain ids in the registry.
func (a *AppContext) ListChainIDs() []int64 {
	return a.chainRegistry.ChainIDs()
}

func (a *AppContext) ListChains() []Chain {
	return a.chainRegistry.All()
}

// FirstChain returns the first chain that satisfies the filter
func (a *AppContext) FirstChain(filter func(Chain) bool) (Chain, error) {
	ids := a.ListChainIDs()

	for _, id := range ids {
		chain, err := a.GetChain(id)
		if err != nil {
			return Chain{}, errors.Wrapf(err, "unable to get chain %d", id)
		}

		if filter(chain) {
			return chain, nil
		}
	}

	return Chain{}, errors.Wrap(ErrChainNotFound, "no chain satisfies the filter")
}

// IsOutboundObservationEnabled returns true if outbound flag is enabled
func (a *AppContext) IsOutboundObservationEnabled() bool {
	return a.GetCrossChainFlags().IsOutboundEnabled
}

// IsInboundObservationEnabled returns true if inbound flag is enabled
func (a *AppContext) IsInboundObservationEnabled() bool {
	return a.GetCrossChainFlags().IsInboundEnabled
}

// GetKeygen returns the current keygen
func (a *AppContext) GetKeygen() observertypes.Keygen {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var copiedPubKeys []string
	if a.keygen.GranteePubkeys != nil {
		copiedPubKeys = make([]string, len(a.keygen.GranteePubkeys))
		copy(copiedPubKeys, a.keygen.GranteePubkeys)
	}

	return observertypes.Keygen{
		Status:         a.keygen.Status,
		GranteePubkeys: copiedPubKeys,
		BlockNumber:    a.keygen.BlockNumber,
	}
}

// GetCurrentTssPubKey returns the current tss pubKey.
func (a *AppContext) GetCurrentTssPubKey() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.currentTssPubKey
}

// GetCrossChainFlags returns crosschain flags
func (a *AppContext) GetCrossChainFlags() observertypes.CrosschainFlags {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.crosschainFlags
}

// Update updates AppContext and params for all chains
// this must be the ONLY function that writes to AppContext
func (a *AppContext) Update(
	keygen observertypes.Keygen,
	freshChains, additionalChains []chains.Chain,
	freshChainParams map[int64]*observertypes.ChainParams,
	tssPubKey string,
	crosschainFlags observertypes.CrosschainFlags,
) error {
	// some sanity checks
	switch {
	case len(freshChains) == 0:
		return fmt.Errorf("no chains present")
	case len(freshChainParams) == 0:
		return fmt.Errorf("no chain params present")
	case tssPubKey == "":
		return fmt.Errorf("tss pubkey is empty")
	case len(additionalChains) > 0:
		for _, c := range additionalChains {
			if !c.IsExternal {
				return fmt.Errorf("additional chain %d is not external", c.ChainId)
			}
		}
	}

	err := a.updateChainRegistry(freshChains, additionalChains, freshChainParams)
	if err != nil {
		return errors.Wrap(err, "unable to update chain registry")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.crosschainFlags = crosschainFlags
	a.keygen = keygen
	a.currentTssPubKey = tssPubKey

	return nil
}

func (a *AppContext) updateChainRegistry(
	freshChains, additionalChains []chains.Chain,
	freshChainParams map[int64]*observertypes.ChainParams,
) error {
	freshChainsByID := make(map[int64]chains.Chain, len(freshChains)+len(additionalChains))
	for _, c := range freshChains {
		freshChainsByID[c.ChainId] = c
	}

	for _, c := range additionalChains {
		// shouldn't happen, but just in case
		if _, found := freshChainsByID[c.ChainId]; found {
			continue
		}

		freshChainsByID[c.ChainId] = c
	}

	var (
		freshChainIDs    = maps.Keys(freshChainsByID)
		existingChainIDs = a.chainRegistry.ChainIDs()
	)

	if len(existingChainIDs) > 0 && !elementsMatch(existingChainIDs, freshChainIDs) {
		a.logger.Warn().
			Ints64("chains.current", existingChainIDs).
			Ints64("chains.new", freshChainIDs).
			Msg("Chain list changed at the runtime!")
	}

	// Log warn if somehow chain doesn't chainParam
	for _, chainID := range freshChainIDs {
		if _, ok := freshChainParams[chainID]; !ok {
			a.logger.Warn().
				Int64("chain.id", chainID).
				Msg("Chain doesn't have according ChainParams present. Skipping.")
		}
	}

	// okay, let's update the chains.
	// Set() ensures that chain, chainID, and params are consistent and chain is not zeta + chain is supported
	for chainID, params := range freshChainParams {
		if err := observertypes.ValidateChainParams(params); err != nil {
			return errors.Wrapf(err, "invalid chain params for chain %d", chainID)
		}

		chain, ok := freshChainsByID[chainID]
		if !ok {
			return fmt.Errorf("unable to locate fresh chain %d based on chain params", chainID)
		}

		if err := a.chainRegistry.Set(chainID, &chain, params); err != nil {
			return errors.Wrap(err, "unable to set chain in the registry")
		}
	}

	a.chainRegistry.SetAdditionalChains(additionalChains)

	toBeDeleted := diff(existingChainIDs, freshChainIDs)
	if len(toBeDeleted) > 0 {
		a.logger.Warn().
			Ints64("chains.deleted", toBeDeleted).
			Msg("Deleting chains that are no longer relevant")

		a.chainRegistry.Delete(toBeDeleted...)
	}

	return nil
}

func elementsMatch(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// diff returns the elements in `a` that are not in `b`
func diff(a, b []int64) []int64 {
	var (
		cache  = map[int64]struct{}{}
		result []int64
	)

	for _, v := range b {
		cache[v] = struct{}{}
	}

	for _, v := range a {
		if _, ok := cache[v]; !ok {
			result = append(result, v)
		}
	}

	return result
}
