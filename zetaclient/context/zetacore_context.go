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

// ZetacoreContext contains zetacore context params
// these are initialized and updated at runtime at every height
type ZetacoreContext struct {
	mu                 *sync.RWMutex
	keygen             observertypes.Keygen
	chainsEnabled      []chains.Chain
	chainParamMap      map[int64]*observertypes.ChainParams
	evmChainParams     map[int64]*observertypes.ChainParams
	bitcoinChainParams *observertypes.ChainParams
	currentTssPubkey   string
	crosschainFlags    observertypes.CrosschainFlags

	// blockHeaderEnabledChains is used to store the list of chains that have block header verification enabled
	// All chains in this list will have Enabled flag set to true
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain
}

// NewZetacoreContext creates and returns new ZetacoreContext
// it is initializing chain params from provided config
func NewZetacoreContext(cfg config.Config) *ZetacoreContext {
	evmChainParams := make(map[int64]*observertypes.ChainParams)
	chainParamsMap := make(map[int64]*observertypes.ChainParams)
	for _, e := range cfg.EVMChainConfigs {
		evmChainParams[e.Chain.ChainId] = &observertypes.ChainParams{}
		chainParamsMap[e.Chain.ChainId] = &observertypes.ChainParams{}
	}

	var bitcoinChainParams *observertypes.ChainParams
	_, found := cfg.GetBTCConfig()
	if found {
		bitcoinChainParams = &observertypes.ChainParams{}
	}

	return &ZetacoreContext{
		mu:                       new(sync.RWMutex),
		chainsEnabled:            []chains.Chain{},
		evmChainParams:           evmChainParams,
		bitcoinChainParams:       bitcoinChainParams,
		crosschainFlags:          observertypes.CrosschainFlags{},
		blockHeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{},
	}
}

// CreateZetacoreContext creates and returns new ZetacoreContext
func CreateZetacoreContext(
	keygen *observertypes.Keygen,
	chainsEnabled []chains.Chain,
	chainParamMap map[int64]*observertypes.ChainParams,
	evmChainParams map[int64]*observertypes.ChainParams,
	bitcoinChainParams *observertypes.ChainParams,
	tssPubKey string,
	crosschainFlags observertypes.CrosschainFlags,
	blockHeaderEnabledChains []lightclienttypes.HeaderSupportedChain,
) *ZetacoreContext {
	return &ZetacoreContext{
		mu:                       new(sync.RWMutex),
		keygen:                   *keygen,
		chainsEnabled:            chainsEnabled,
		chainParamMap:            chainParamMap,
		evmChainParams:           evmChainParams,
		bitcoinChainParams:       bitcoinChainParams,
		currentTssPubkey:         tssPubKey,
		crosschainFlags:          crosschainFlags,
		blockHeaderEnabledChains: blockHeaderEnabledChains,
	}
}

func (c *ZetacoreContext) GetKeygen() observertypes.Keygen {
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

func (c *ZetacoreContext) GetCurrentTssPubkey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentTssPubkey
}

// GetEnabledExternalChains returns all enabled external chains (excluding zetachain)
func (c *ZetacoreContext) GetEnabledExternalChains() []chains.Chain {
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

func (c *ZetacoreContext) GetEVMChainParams(chainID int64) (*observertypes.ChainParams, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	evmChainParams, found := c.evmChainParams[chainID]
	return evmChainParams, found
}

func (c *ZetacoreContext) GetAllEVMChainParams() map[int64]*observertypes.ChainParams {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy evm chain params
	copied := make(map[int64]*observertypes.ChainParams, len(c.evmChainParams))
	for chainID, evmConfig := range c.evmChainParams {
		copied[chainID] = &observertypes.ChainParams{}
		*copied[chainID] = *evmConfig
	}
	return copied
}

// GetBTCChainParams returns (chain, chain params, found) for bitcoin chain
func (c *ZetacoreContext) GetBTCChainParams() (chains.Chain, *observertypes.ChainParams, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.bitcoinChainParams == nil { // bitcoin is not enabled
		return chains.Chain{}, nil, false
	}

	chain := chains.GetChainFromChainID(c.bitcoinChainParams.ChainId)
	if chain == nil {
		return chains.Chain{}, nil, false
	}

	return *chain, c.bitcoinChainParams, true
}

func (c *ZetacoreContext) GetCrossChainFlags() observertypes.CrosschainFlags {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.crosschainFlags
}

// GetAllHeaderEnabledChains returns all verification flags
func (c *ZetacoreContext) GetAllHeaderEnabledChains() []lightclienttypes.HeaderSupportedChain {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.blockHeaderEnabledChains
}

// GetBlockHeaderEnabledChains checks if block header verification is enabled for a specific chain
func (c *ZetacoreContext) GetBlockHeaderEnabledChains(chainID int64) (lightclienttypes.HeaderSupportedChain, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, flags := range c.blockHeaderEnabledChains {
		if flags.ChainId == chainID {
			return flags, true
		}
	}
	return lightclienttypes.HeaderSupportedChain{}, false
}

// Update updates zetacore context and params for all chains
// this must be the ONLY function that writes to zetacore context
func (c *ZetacoreContext) Update(
	keygen *observertypes.Keygen,
	newChains []chains.Chain,
	evmChainParams map[int64]*observertypes.ChainParams,
	btcChainParams *observertypes.ChainParams,
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

	if keygen != nil {
		c.keygen = *keygen
	}

	c.chainsEnabled = newChains
	c.crosschainFlags = crosschainFlags
	c.blockHeaderEnabledChains = blockHeaderEnabledChains

	// update chain params for bitcoin if it has config in file
	if c.bitcoinChainParams != nil && btcChainParams != nil {
		c.bitcoinChainParams = btcChainParams
	}

	// update core params for evm chains we have configs in file
	for _, params := range evmChainParams {
		_, found := c.evmChainParams[params.ChainId]
		if !found {
			continue
		}
		c.evmChainParams[params.ChainId] = params
	}

	if tssPubKey != "" {
		c.currentTssPubkey = tssPubKey
	}
}

// IsOutboundObservationEnabled returns true if the chain is supported and outbound flag is enabled
func IsOutboundObservationEnabled(c *ZetacoreContext, chainParams observertypes.ChainParams) bool {
	flags := c.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsOutboundEnabled
}

// IsInboundObservationEnabled returns true if the chain is supported and inbound flag is enabled
func IsInboundObservationEnabled(c *ZetacoreContext, chainParams observertypes.ChainParams) bool {
	flags := c.GetCrossChainFlags()
	return chainParams.IsSupported && flags.IsInboundEnabled
}
