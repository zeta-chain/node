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
	chainsEnabled    []chains.Chain
	chainParamMap    map[int64]*observertypes.ChainParams
	btcNetParams     *chaincfg.Params
	currentTssPubkey string
	crosschainFlags  observertypes.CrosschainFlags

	// blockHeaderEnabledChains is used to store the list of chains that have block header verification enabled
	// All chains in this list will have Enabled flag set to true
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain

	// mu is to protect the app context from concurrent access
	mu *sync.RWMutex
}

// NewAppContext creates and returns new AppContext
// it is initializing chain params from provided config
func NewAppContext(cfg config.Config) *AppContext {
	return &AppContext{
		config:                   cfg,
		chainsEnabled:            []chains.Chain{},
		chainParamMap:            make(map[int64]*observertypes.ChainParams),
		crosschainFlags:          observertypes.CrosschainFlags{},
		blockHeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{},
		mu:                       new(sync.RWMutex),
	}
}

// CreateAppContext creates a new AppContext
func CreateAppContext(
	keygen observertypes.Keygen,
	chainsEnabled []chains.Chain,
	chainParamMap map[int64]*observertypes.ChainParams,
	btcNetParams *chaincfg.Params,
	tssPubKey string,
	crosschainFlags observertypes.CrosschainFlags,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
) *AppContext {
	return &AppContext{
		mu:                       new(sync.RWMutex),
		keygen:                   keygen,
		chainsEnabled:            chainsEnabled,
		chainParamMap:            chainParamMap,
		btcNetParams:             btcNetParams,
		currentTssPubkey:         tssPubKey,
		crosschainFlags:          crosschainFlags,
		blockHeaderEnabledChains: blockHeaderEnabledChains,
	}
}

// Config returns the app context config
func (c *AppContext) Config() config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

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

func (c *AppContext) GetCurrentTssPubkey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentTssPubkey
}

// GetEnabledExternalChains returns all enabled external chains (excluding zetachain)
func (c *AppContext) GetEnabledExternalChains() []chains.Chain {
	c.mu.RLock()
	defer c.mu.RUnlock()

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

// GetAllHeaderEnabledChains returns all verification flags
func (c *AppContext) GetAllHeaderEnabledChains() []lightclienttypes.HeaderSupportedChain {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.blockHeaderEnabledChains
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

// Update updates the inner config and app context
func (c *AppContext) UpdateContext(config config.Config, newContext *AppContext, logger zerolog.Logger) {
	c.Update(
		config,
		newContext.GetKeygen(),
		newContext.GetEnabledExternalChains(),
		newContext.GetEnabledExternalChainParams(),
		newContext.GetBTCNetParams(),
		newContext.GetCurrentTssPubkey(),
		newContext.GetCrossChainFlags(),
		newContext.GetAllHeaderEnabledChains(),
		false,
		logger,
	)
}

// Update updates app context and params for all chains
// this must be the ONLY function that writes to app context
func (c *AppContext) Update(
	config config.Config,
	keygen observertypes.Keygen,
	newChains []chains.Chain,
	newChainParams map[int64]*observertypes.ChainParams,
	btcNetParams *chaincfg.Params,
	tssPubKey string,
	crosschainFlags observertypes.CrosschainFlags,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
	init bool,
	logger zerolog.Logger,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ignore whatever order zetacore organizes chain list in state
	sort.SliceStable(newChains, func(i, j int) bool {
		return newChains[i].ChainId < newChains[j].ChainId
	})

	if len(newChains) == 0 {
		logger.Warn().Msg("UpdateChainParams: No chains enabled in ZeroCore")
	}

	// Add some warnings if chain list changes at runtime
	if !init {
		if len(c.chainsEnabled) != len(newChains) {
			logger.Warn().Msgf(
				"UpdateChainParams: ChainsEnabled changed at runtime!! current: %v, new: %v",
				c.chainsEnabled,
				newChains,
			)
		} else {
			for i, chain := range newChains {
				if chain != c.chainsEnabled[i] {
					logger.Warn().Msgf(
						"UpdateChainParams: ChainsEnabled changed at runtime!! current: %v, new: %v",
						c.chainsEnabled,
						newChains,
					)
				}
			}
		}
	}

	// btcNetParams points one of [mainnet, testnet, regnet]
	// btcNetParams initialize only once and should never change
	if c.btcNetParams == nil {
		c.btcNetParams = btcNetParams
	}

	c.config = config
	c.keygen = keygen
	c.chainsEnabled = newChains
	c.chainParamMap = newChainParams
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
