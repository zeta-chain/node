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
	config           *config.Config
	keygen           observertypes.Keygen
	currentTssPubkey string
	chainsEnabled    []chains.Chain
	chainParamMap    map[int64]*observertypes.ChainParams
	btcNetParams     *chaincfg.Params
	crosschainFlags  observertypes.CrosschainFlags

	// blockHeaderEnabledChains is used to store the list of chains that have block header verification enabled
	// All chains in this list will have Enabled flag set to true
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain

	// mu is to protect the app context from concurrent access
	mu sync.RWMutex
}

// NewAppContext creates empty app context with given config
func NewAppContext(cfg *config.Config) *AppContext {
	return &AppContext{
		config:                   cfg,
		chainsEnabled:            []chains.Chain{},
		chainParamMap:            make(map[int64]*observertypes.ChainParams),
		crosschainFlags:          observertypes.CrosschainFlags{},
		blockHeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{},
		mu:                       sync.RWMutex{},
	}
}

// SetConfig sets a new config to the app context
func (c *AppContext) SetConfig(cfg *config.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = cfg
}

// Config returns the app context config
func (c *AppContext) Config() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// GetKeygen returns the current keygen information
func (c *AppContext) GetKeygen() observertypes.Keygen {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var copiedPubkeys []string
	if c.keygen.GranteePubkeys != nil {
		copiedPubkeys = make([]string, len(c.keygen.GranteePubkeys))
		copy(copiedPubkeys, c.keygen.GranteePubkeys)
	}

	return observertypes.Keygen{
		Status:         c.keygen.Status,
		GranteePubkeys: copiedPubkeys,
		BlockNumber:    c.keygen.BlockNumber,
	}
}

// GetCurrentTssPubkey returns the current TSS public key
func (c *AppContext) GetCurrentTssPubkey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentTssPubkey
}

// GetEnabledExternalChains returns all enabled external chains (excluding zetachain)
func (c *AppContext) GetEnabledExternalChains() []chains.Chain {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy chains
	externalChains := make([]chains.Chain, 0)
	for _, chain := range c.chainsEnabled {
		if chain.IsExternal {
			externalChains = append(externalChains, chain)
		}
	}
	return externalChains
}

// GetEnabledBTCChains returns the enabled bitcoin chains
func (c *AppContext) GetEnabledBTCChains() []chains.Chain {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy btc chains
	btcChains := make([]chains.Chain, 0)
	for _, chain := range c.chainsEnabled {
		if chain.Consensus == chains.Consensus_bitcoin {
			btcChains = append(btcChains, chain)
		}
	}
	return btcChains
}

// GetEnabledExternalChainParams returns all enabled chain params
func (c *AppContext) GetEnabledExternalChainParams() map[int64]*observertypes.ChainParams {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy chain params
	copied := make(map[int64]*observertypes.ChainParams, len(c.chainParamMap))
	for chainID, chainParams := range c.chainParamMap {
		copied[chainID] = &observertypes.ChainParams{}
		*copied[chainID] = *chainParams
	}
	return copied
}

// GetExternalChainParams returns chain params for a specific chain ID
func (c *AppContext) GetExternalChainParams(chainID int64) (*observertypes.ChainParams, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	chainParams, found := c.chainParamMap[chainID]
	return chainParams, found
}

// GetBTCNetParams returns bitcoin network params
func (c *AppContext) GetBTCNetParams() *chaincfg.Params {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.btcNetParams
}

// GetCrossChainFlags returns crosschain flags
func (c *AppContext) GetCrossChainFlags() observertypes.CrosschainFlags {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.crosschainFlags
}

// GetBlockHeaderEnabledChains checks if block header verification is enabled for a specific chain
func (c *AppContext) GetBlockHeaderEnabledChains(chainID int64) (lightclienttypes.HeaderSupportedChain, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, flags := range c.blockHeaderEnabledChains {
		if flags.ChainId == chainID {
			return flags, true
		}
	}
	return lightclienttypes.HeaderSupportedChain{}, false
}

// Update updates app context and params for all chains
// this must be the ONLY function that writes to app context
func (c *AppContext) Update(
	keygen observertypes.Keygen,
	tssPubKey string,
	chainsEnabled []chains.Chain,
	chainParamMap map[int64]*observertypes.ChainParams,
	btcNetParams *chaincfg.Params,
	crosschainFlags observertypes.CrosschainFlags,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
	logger zerolog.Logger,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ignore whatever order zetacore organizes chain list in state
	sort.SliceStable(chainsEnabled, func(i, j int) bool {
		return chainsEnabled[i].ChainId < chainsEnabled[j].ChainId
	})

	if len(chainsEnabled) == 0 {
		logger.Warn().Msg("UpdateChainParams: no external chain enabled in the zetacore")
	}

	// Add log print if the number of enabled chains changes at runtime
	if len(c.chainsEnabled) != len(chainsEnabled) {
		logger.Info().Msgf(
			"UpdateChainParams: number of enabled chains changed at runtime!! before: %d, after: %d",
			len(c.chainsEnabled),
			len(chainsEnabled),
		)
	}

	// btcNetParams points one of [mainnet, testnet, regnet]
	// btcNetParams initialize only once and should never change
	if c.btcNetParams == nil {
		c.btcNetParams = btcNetParams
	}

	c.keygen = keygen
	c.chainsEnabled = chainsEnabled
	c.chainParamMap = chainParamMap
	c.currentTssPubkey = tssPubKey
	c.crosschainFlags = crosschainFlags
	c.blockHeaderEnabledChains = blockHeaderEnabledChains
}

// IsOutboundObservationEnabled returns true if the chain is supported and outbound flag is enabled
func IsOutboundObservationEnabled(c *AppContext, chainParams observertypes.ChainParams) bool {
	flags := c.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsOutboundEnabled
}

// IsInboundObservationEnabled returns true if the chain is supported and inbound flag is enabled
func IsInboundObservationEnabled(c *AppContext, chainParams observertypes.ChainParams) bool {
	flags := c.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsInboundEnabled
}
