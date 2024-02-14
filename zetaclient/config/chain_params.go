package config

import (
	"fmt"
	"sort"
	"sync"

	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// TODO: update name and package
type Params struct {
	paramsLock         *sync.RWMutex
	Keygen             *observertypes.Keygen
	ChainsEnabled      []common.Chain
	EVMChainParams     map[int64]*observertypes.ChainParams
	BitcoinChainParams *observertypes.ChainParams
	CurrentTssPubkey   string
}

func NewParams() *Params {
	return &Params{
		paramsLock:         new(sync.RWMutex),
		ChainsEnabled:      []common.Chain{},
		EVMChainParams:     make(map[int64]*observertypes.ChainParams),
		BitcoinChainParams: &observertypes.ChainParams{},
	}
}

func (c *Params) GetKeygen() observertypes.Keygen {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()
	copiedPubkeys := make([]string, len(c.Keygen.GranteePubkeys))
	copy(copiedPubkeys, c.Keygen.GranteePubkeys)

	return observertypes.Keygen{
		Status:         c.Keygen.Status,
		GranteePubkeys: copiedPubkeys,
		BlockNumber:    c.Keygen.BlockNumber,
	}
}

func (c *Params) GetEnabledChains() []common.Chain {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()
	copiedChains := make([]common.Chain, len(c.ChainsEnabled))
	copy(copiedChains, c.ChainsEnabled)
	return copiedChains
}

func (c *Params) GetEVMChainParams(chainID int64) (*observertypes.ChainParams, bool) {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()
	evmChainParams, found := c.EVMChainParams[chainID]
	return evmChainParams, found
}

func (c *Params) GetAllEVMChainParams() map[int64]*observertypes.ChainParams {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()

	// deep copy evm chain params
	copied := make(map[int64]*observertypes.ChainParams, len(c.EVMChainParams))
	for chainID, evmConfig := range c.EVMChainParams {
		copied[chainID] = &observertypes.ChainParams{}
		*copied[chainID] = *evmConfig
	}
	return copied
}

func (c *Params) GetBTCChainParams() (common.Chain, *observertypes.ChainParams, bool) {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()

	if c.BitcoinChainParams == nil { // bitcoin is not enabled
		return common.Chain{}, &observertypes.ChainParams{}, false
	}
	chain := common.GetChainFromChainID(c.BitcoinChainParams.ChainId)
	if chain == nil {
		panic(fmt.Sprintf("BTCChain is missing for chainID %d", c.BitcoinChainParams.ChainId))
	}
	return *chain, c.BitcoinChainParams, true
}

// UpdateChainParams updates core params for all chains
// this must be the ONLY function that writes to core params
func (c *Params) UpdateChainParams(
	keygen *observertypes.Keygen,
	newChains []common.Chain,
	evmChainParams map[int64]*observertypes.ChainParams,
	btcChainParams *observertypes.ChainParams,
	init bool,
	logger zerolog.Logger,
) {
	c.paramsLock.Lock()
	defer c.paramsLock.Unlock()

	// Ignore whatever order zetabridge organizes chain list in state
	sort.SliceStable(newChains, func(i, j int) bool {
		return newChains[i].ChainId < newChains[j].ChainId
	})
	if len(newChains) == 0 {
		logger.Warn().Msg("UpdateChainParams: No chains enabled in ZeroCore")
	}

	// Add some warnings if chain list changes at runtime
	if !init {
		if len(c.ChainsEnabled) != len(newChains) {
			logger.Warn().Msgf(
				"UpdateChainParams: ChainsEnabled changed at runtime!! current: %v, new: %v",
				c.ChainsEnabled,
				newChains,
			)
		} else {
			for i, chain := range newChains {
				if chain != c.ChainsEnabled[i] {
					logger.Warn().Msgf(
						"UpdateChainParams: ChainsEnabled changed at runtime!! current: %v, new: %v",
						c.ChainsEnabled,
						newChains,
					)
				}
			}
		}
	}
	c.Keygen = keygen
	c.ChainsEnabled = newChains
	// update chain params for bitcoin if it has config in file
	if c.BitcoinChainParams != nil && btcChainParams != nil {
		c.BitcoinChainParams = btcChainParams
	}
	// update core params for evm chains we have configs in file
	for _, params := range evmChainParams {
		_, found := c.EVMChainParams[params.ChainId]
		if !found {
			continue
		}
		c.EVMChainParams[params.ChainId] = params
	}
}
