// Package context provides global app context for ZetaClient
package context

import (
	"sort"
	"sync"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AppContext represents application context.
type AppContext struct {
	config config.Config
	logger zerolog.Logger

	keygen             observertypes.Keygen
	chainsEnabled      []chains.Chain
	evmChainParams     map[int64]*observertypes.ChainParams
	bitcoinChainParams *observertypes.ChainParams
	currentTssPubkey   string
	crosschainFlags    observertypes.CrosschainFlags

	// additionalChains is a list of additional static chain information to use when searching from chain IDs
	// it is stored in the protocol to dynamically support new chains without doing an upgrade
	additionalChain []chains.Chain

	// blockHeaderEnabledChains is used to store the list of chains that have block header verification enabled
	// All chains in this list will have Enabled flag set to true
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain

	mu sync.RWMutex
}

// New creates and returns new AppContext
func New(cfg config.Config, logger zerolog.Logger) *AppContext {
	evmChainParams := make(map[int64]*observertypes.ChainParams)
	for _, e := range cfg.EVMChainConfigs {
		evmChainParams[e.Chain.ChainId] = &observertypes.ChainParams{}
	}

	var bitcoinChainParams *observertypes.ChainParams
	_, found := cfg.GetBTCConfig()
	if found {
		bitcoinChainParams = &observertypes.ChainParams{}
	}

	return &AppContext{
		config: cfg,
		logger: logger.With().Str("module", "appcontext").Logger(),

		chainsEnabled:            []chains.Chain{},
		evmChainParams:           evmChainParams,
		bitcoinChainParams:       bitcoinChainParams,
		crosschainFlags:          observertypes.CrosschainFlags{},
		blockHeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{},

		currentTssPubkey: "",
		keygen:           observertypes.Keygen{},
		mu:               sync.RWMutex{},
	}
}

// Config returns the config of the app
func (a *AppContext) Config() config.Config {
	return a.config
}

// GetEnabledBTCChains returns the enabled solana chains
func (a *AppContext) GetSolanaChainAndConfig() (chains.Chain, config.SolanaConfig, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// FIXME_SOLANA: config this
	chain := chains.SolanaLocalnet
	config, enabled := a.Config().GetSolanaConfig()
	return chain, config, enabled
}

// GetBTCChainAndConfig returns btc chain and config if enabled
func (a *AppContext) GetBTCChainAndConfig() (chains.Chain, config.BTCConfig, bool) {
	btcConfig, configEnabled := a.Config().GetBTCConfig()
	btcChain, _, paramsEnabled := a.GetBTCChainParams()

	if !configEnabled || !paramsEnabled {
		return chains.Chain{}, config.BTCConfig{}, false
	}

	return btcChain, btcConfig, true
}

// IsOutboundObservationEnabled returns true if the chain is supported and outbound flag is enabled
func (a *AppContext) IsOutboundObservationEnabled(chainParams observertypes.ChainParams) bool {
	flags := a.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsOutboundEnabled
}

// IsInboundObservationEnabled returns true if the chain is supported and inbound flag is enabled
func (a *AppContext) IsInboundObservationEnabled(chainParams observertypes.ChainParams) bool {
	flags := a.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsInboundEnabled
}

// GetKeygen returns the current keygen
func (a *AppContext) GetKeygen() observertypes.Keygen {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var copiedPubkeys []string
	if a.keygen.GranteePubkeys != nil {
		copiedPubkeys = make([]string, len(a.keygen.GranteePubkeys))
		copy(copiedPubkeys, a.keygen.GranteePubkeys)
	}

	return observertypes.Keygen{
		Status:         a.keygen.Status,
		GranteePubkeys: copiedPubkeys,
		BlockNumber:    a.keygen.BlockNumber,
	}
}

// GetCurrentTssPubKey returns the current tss pubkey
func (a *AppContext) GetCurrentTssPubKey() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.currentTssPubkey
}

// GetEnabledChains returns all enabled chains including zetachain
func (a *AppContext) GetEnabledChains() []chains.Chain {
	a.mu.RLock()
	defer a.mu.RUnlock()

	copiedChains := make([]chains.Chain, len(a.chainsEnabled))
	copy(copiedChains, a.chainsEnabled)

	return copiedChains
}

// GetEnabledExternalChains returns all enabled external chains
func (a *AppContext) GetEnabledExternalChains() []chains.Chain {
	a.mu.RLock()
	defer a.mu.RUnlock()

	externalChains := make([]chains.Chain, 0)
	for _, chain := range a.chainsEnabled {
		if chain.IsExternal {
			externalChains = append(externalChains, chain)
		}
	}
	return externalChains
}

// GetEVMChainParams returns chain params for a specific EVM chain
func (a *AppContext) GetEVMChainParams(chainID int64) (*observertypes.ChainParams, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	evmChainParams, found := a.evmChainParams[chainID]
	return evmChainParams, found
}

// GetAllEVMChainParams returns all chain params for EVM chains
func (a *AppContext) GetAllEVMChainParams() map[int64]*observertypes.ChainParams {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// deep copy evm chain params
	copied := make(map[int64]*observertypes.ChainParams, len(a.evmChainParams))
	for chainID, evmConfig := range a.evmChainParams {
		copied[chainID] = &observertypes.ChainParams{}
		*copied[chainID] = *evmConfig
	}
	return copied
}

// GetBTCChainParams returns (chain, chain params, found) for bitcoin chain
func (a *AppContext) GetBTCChainParams() (chains.Chain, *observertypes.ChainParams, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.bitcoinChainParams == nil { // bitcoin is not enabled
		return chains.Chain{}, nil, false
	}

	chain, found := chains.GetChainFromChainID(a.bitcoinChainParams.ChainId, a.additionalChain)
	if !found {
		return chains.Chain{}, nil, false
	}

	return chain, a.bitcoinChainParams, true
}

// GetCrossChainFlags returns crosschain flags
func (a *AppContext) GetCrossChainFlags() observertypes.CrosschainFlags {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.crosschainFlags
}

// GetAdditionalChains returns additional chains
func (a *AppContext) GetAdditionalChains() []chains.Chain {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.additionalChain
}

// GetAllHeaderEnabledChains returns all verification flags
func (a *AppContext) GetAllHeaderEnabledChains() []lightclienttypes.HeaderSupportedChain {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.blockHeaderEnabledChains
}

// GetBlockHeaderEnabledChains checks if block header verification is enabled for a specific chain
func (a *AppContext) GetBlockHeaderEnabledChains(chainID int64) (lightclienttypes.HeaderSupportedChain, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, flags := range a.blockHeaderEnabledChains {
		if flags.ChainId == chainID {
			return flags, true
		}
	}

	return lightclienttypes.HeaderSupportedChain{}, false
}

// Update updates zetacore context and params for all chains
// this must be the ONLY function that writes to zetacore context
func (a *AppContext) Update(
	keygen *observertypes.Keygen,
	newChains []chains.Chain,
	evmChainParams map[int64]*observertypes.ChainParams,
	btcChainParams *observertypes.ChainParams,
	tssPubKey string,
	crosschainFlags observertypes.CrosschainFlags,
	additionalChains []chains.Chain,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
	init bool,
) {
	// Ignore whatever order zetacore organizes chain list in state
	sort.SliceStable(newChains, func(i, j int) bool {
		return newChains[i].ChainId < newChains[j].ChainId
	})

	if len(newChains) == 0 {
		a.logger.Warn().Msg("UpdateChainParams: No chains enabled in ZeroCore")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Add some warnings if chain list changes at runtime
	if !init && !chainsEqual(a.chainsEnabled, newChains) {
		a.logger.Warn().
			Interface("chains.current", a.chainsEnabled).
			Interface("chains.new", newChains).
			Msg("UpdateChainParams: ChainsEnabled changed at runtime!")
	}

	if keygen != nil {
		a.keygen = *keygen
	}

	a.chainsEnabled = newChains
	a.crosschainFlags = crosschainFlags
	a.additionalChain = additionalChains
	a.blockHeaderEnabledChains = blockHeaderEnabledChains

	// update chain params for bitcoin if it has config in file
	if a.bitcoinChainParams != nil && btcChainParams != nil {
		a.bitcoinChainParams = btcChainParams
	}

	// update core params for evm chains we have configs in file
	for _, params := range evmChainParams {
		_, found := a.evmChainParams[params.ChainId]
		if !found {
			continue
		}
		a.evmChainParams[params.ChainId] = params
	}

	if tssPubKey != "" {
		a.currentTssPubkey = tssPubKey
	}
}

func chainsEqual(a []chains.Chain, b []chains.Chain) bool {
	if len(a) != len(b) {
		return false
	}

	for i, left := range a {
		right := b[i]

		if left.ChainId != right.ChainId || left.ChainName != right.ChainName {
			return false
		}
	}

	return true
}
