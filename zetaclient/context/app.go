// Package context provides global app context for ZetaClient
package context

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
)

// AppContext represents application (zetaclient) context.
type AppContext struct {
	// config is the config of the app
	config config.Config

	// logger is the logger of the app
	logger zerolog.Logger

	// chainRegistry is a registry of supported chains
	chainRegistry *ChainRegistry

	// currentTssPubKey is the current tss pubKey
	currentTssPubKey string

	// crosschainFlags is the current crosschain flags state
	crosschainFlags observertypes.CrosschainFlags

	// keygen is the current tss keygen state
	keygen observertypes.Keygen

	mu sync.RWMutex
}

// New creates and returns new empty AppContext
func New(cfg config.Config, relayerKeyPasswords map[string]string, logger zerolog.Logger) *AppContext {
	return &AppContext{
		config: cfg,
		logger: logger.With().Str("module", "appcontext").Logger(),

		chainRegistry: NewChainRegistry(relayerKeyPasswords),

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

// ListChains returns the list of existing chains in the registry.
func (a *AppContext) ListChains() []Chain {
	return a.chainRegistry.All()
}

// FilterChains returns the list of chains that satisfy the filter
func (a *AppContext) FilterChains(filter func(Chain) bool) []Chain {
	var (
		all = a.ListChains()
		out = make([]Chain, 0, len(all))
	)

	for _, chain := range all {
		if filter(chain) {
			out = append(out, chain)
		}
	}

	return out
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
	case tssPubKey == "" && a.currentTssPubKey != "":
		// note that if we're doing a fresh start, we ALLOW an empty tssPubKey
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

// updateChainRegistry updates the chain registry with fresh chains and chain params.
// Note that there's an edge-case for ZetaChain itself because we WANT to have it in chains list,
// but it doesn't have chain params.
func (a *AppContext) updateChainRegistry(
	freshChains []chains.Chain,
	additionalChains []chains.Chain,
	freshChainParams map[int64]*observertypes.ChainParams,
) error {
	var zetaChainID int64

	// 1. build map[chainId]Chain
	freshChainsByID := make(map[int64]chains.Chain, len(freshChains)+len(additionalChains))
	for _, c := range freshChains {
		freshChainsByID[c.ChainId] = c

		if isZeta(c.ChainId) && zetaChainID == 0 {
			zetaChainID = c.ChainId
		}
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

	// 2. Compare existing chains with fresh ones
	if len(existingChainIDs) > 0 && !elementsMatch(existingChainIDs, freshChainIDs) {
		a.logger.Warn().
			Ints64("chains.current", existingChainIDs).
			Ints64("chains.new", freshChainIDs).
			Msg("Chain list changed at the runtime!")
	}

	// Log warn if somehow chain doesn't chainParam
	for _, chainID := range freshChainIDs {
		if _, ok := freshChainParams[chainID]; !ok && !isZeta(chainID) {
			a.logger.Warn().
				Int64("chain.id", chainID).
				Msg("Chain doesn't have according ChainParams present. Skipping.")
		}
	}

	// 3. If we have zeta chain, we want to force "fake" chainParams for it
	if zetaChainID != 0 {
		freshChainParams[zetaChainID] = zetaObserverChainParams(zetaChainID)
	}

	// 3. Update chain registry
	// okay, let's update the chains.
	// Set() ensures that chain, chainID, and params are consistent and chain is not zeta + chain is supported
	for chainID, params := range freshChainParams {
		chain, ok := freshChainsByID[chainID]
		if !ok {
			return fmt.Errorf("unable to locate fresh chain %d based on chain params", chainID)
		}

		if !isZeta(chainID) {
			if err := observertypes.ValidateChainParams(params); err != nil {
				return errors.Wrapf(err, "invalid chain params for chain %d", chainID)
			}
		}

		if err := a.chainRegistry.Set(chainID, &chain, params); err != nil {
			return errors.Wrap(err, "unable to set chain in the registry")
		}
	}

	a.chainRegistry.SetAdditionalChains(additionalChains)

	toBeDeleted, _ := lo.Difference(existingChainIDs, freshChainIDs)
	if len(toBeDeleted) > 0 {
		a.logger.Warn().
			Ints64("chains.deleted", toBeDeleted).
			Msg("Deleting chains that are no longer relevant")

		a.chainRegistry.Delete(toBeDeleted...)
	}

	return nil
}

func isZeta(chainID int64) bool {
	return chains.IsZetaChain(chainID, nil)
}

// zetaObserverChainParams returns "fake" chain params because
// actually chainParams is a concept of observer
func zetaObserverChainParams(chainID int64) *observertypes.ChainParams {
	return &observertypes.ChainParams{ChainId: chainID, IsSupported: true}
}

// elementsMatch returns true if two slices are equal.
// SORTS the slices before comparison.
func elementsMatch[T constraints.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
