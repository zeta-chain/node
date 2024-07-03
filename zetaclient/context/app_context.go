// Package context provides global app context for ZetaClient
package context

import (
	"sort"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AppContext contains zetaclient application context
// these are initialized and updated at runtime periodically
type AppContext struct {
	config           config.Config
	keygen           observertypes.Keygen
	currentTssPubkey string
	chainsEnabled    []chains.Chain
	chainParamMap    map[int64]*observertypes.ChainParams
	btcNetParams     *chaincfg.Params
	crosschainFlags  observertypes.CrosschainFlags

	// additionalChains is a list of additional static chain information to use when searching from chain IDs
	// it is stored in the protocol to dynamically support new chains without doing an upgrade
	additionalChain []chains.Chain

	// blockHeaderEnabledChains is used to store the list of chains that have block header verification enabled
	// All chains in this list will have Enabled flag set to true
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain

	// mu is to protect the app context from concurrent access
	mu sync.RWMutex
}

// New creates empty app context with given config
func New(cfg config.Config) *AppContext {
	return &AppContext{
		config:                   cfg,
		chainsEnabled:            []chains.Chain{},
		chainParamMap:            make(map[int64]*observertypes.ChainParams),
		crosschainFlags:          observertypes.CrosschainFlags{},
		additionalChain:          []chains.Chain{},
		blockHeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{},
		mu:                       sync.RWMutex{},
	}
}

// SetConfig sets a new config to the app context
func (a *AppContext) SetConfig(cfg config.Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.config = cfg
}

// Config returns the app context config
func (a *AppContext) Config() config.Config {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.config.Clone()
}

// GetKeygen returns the current keygen information
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

// GetCurrentTssPubkey returns the current TSS public key
func (a *AppContext) GetCurrentTssPubkey() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.currentTssPubkey
}

// GetEnabledExternalChains returns all enabled external chains (excluding zetachain)
func (a *AppContext) GetEnabledExternalChains() []chains.Chain {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// deep copy chains
	externalChains := make([]chains.Chain, 0)
	for _, chain := range a.chainsEnabled {
		if chain.IsExternal {
			externalChains = append(externalChains, chain)
		}
	}
	return externalChains
}

// GetEnabledBTCChains returns the enabled bitcoin chains
func (a *AppContext) GetEnabledBTCChains() []chains.Chain {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// deep copy btc chains
	btcChains := make([]chains.Chain, 0)
	for _, chain := range a.chainsEnabled {
		if chain.Consensus == chains.Consensus_bitcoin {
			btcChains = append(btcChains, chain)
		}
	}
	return btcChains
}

// GetEnabledExternalChainParams returns all enabled chain params
func (a *AppContext) GetEnabledExternalChainParams() map[int64]*observertypes.ChainParams {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// deep copy chain params
	copied := make(map[int64]*observertypes.ChainParams, len(a.chainParamMap))
	for chainID, chainParams := range a.chainParamMap {
		copied[chainID] = &observertypes.ChainParams{}
		*copied[chainID] = *chainParams
	}
	return copied
}

// GetExternalChainParams returns chain params for a specific chain ID
func (a *AppContext) GetExternalChainParams(chainID int64) (*observertypes.ChainParams, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	chainParams, found := a.chainParamMap[chainID]
	return chainParams, found
}

// GetBTCNetParams returns bitcoin network params
func (a *AppContext) GetBTCNetParams() *chaincfg.Params {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.btcNetParams
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

	// deep copy additional chains
	additionalChains := make([]chains.Chain, len(a.additionalChain))
	copy(additionalChains, a.additionalChain)

	return a.additionalChain
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

// Update updates app context and params for all chains
// this must be the ONLY function that writes to app context
func (a *AppContext) Update(
	keygen observertypes.Keygen,
	tssPubKey string,
	chainsEnabled []chains.Chain,
	chainParamMap map[int64]*observertypes.ChainParams,
	btcNetParams *chaincfg.Params,
	crosschainFlags observertypes.CrosschainFlags,
	additionalChains []chains.Chain,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
	logger zerolog.Logger,
) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Ignore whatever order zetacore organizes chain list in state
	sort.SliceStable(chainsEnabled, func(i, j int) bool {
		return chainsEnabled[i].ChainId < chainsEnabled[j].ChainId
	})

	if len(chainsEnabled) == 0 {
		logger.Warn().Msg("UpdateChainParams: no external chain enabled in the zetacore")
	}

	// Add log print if the number of enabled chains changes at runtime
	if len(a.chainsEnabled) != len(chainsEnabled) {
		logger.Info().Msgf(
			"UpdateChainParams: number of enabled chains changed at runtime!! before: %d, after: %d",
			len(a.chainsEnabled),
			len(chainsEnabled),
		)
	}

	// btcNetParams points one of [mainnet, testnet, regnet]
	// btcNetParams initialize only once and should never change
	if a.btcNetParams == nil {
		a.btcNetParams = btcNetParams
	}

	a.keygen = keygen
	a.chainsEnabled = chainsEnabled
	a.chainParamMap = chainParamMap
	a.currentTssPubkey = tssPubKey
	a.crosschainFlags = crosschainFlags
	a.additionalChain = additionalChains
	a.blockHeaderEnabledChains = blockHeaderEnabledChains
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
