package config

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	ostypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

type EVMConfig struct {
	ostypes.CoreParams
	Chain    common.Chain
	Endpoint string
}

type BTCConfig struct {
	ostypes.CoreParams

	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCHost     string
	RPCParams   string // "regtest", "mainnet", "testnet3"
}

type Config struct {
	Peer                string        `json:"Peer"`
	PublicIP            string        `json:"PublicIP"`
	LogFormat           string        `json:"LogFormat"`
	LogLevel            zerolog.Level `json:"LogLevel"`
	LogSampler          bool          `json:"LogSampler"`
	PreParamsPath       string        `json:"PreParamsPath"`
	ChainID             string        `json:"ChainID"`
	ZetaCoreURL         string        `json:"ZetaCoreURL"`
	AuthzGranter        string        `json:"AuthzGranter"`
	AuthzHotkey         string        `json:"AuthzHotkey"`
	P2PDiagnostic       bool          `json:"P2PDiagnostic"`
	ConfigUpdateTicker  uint64        `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker uint64        `json:"P2PDiagnosticTicker"`
	TssPath             string        `json:"TssPath"`
	TestTssKeysign      bool          `json:"TestTssKeysign"`

	// chain specific fields are updatable at runtime and shared across threads
	cfgLock         *sync.RWMutex        `json:"-"`
	Keygen          ostypes.Keygen       `json:"Keygen"`
	ChainsEnabled   []common.Chain       `json:"ChainsEnabled"`
	EVMChainConfigs map[int64]*EVMConfig `json:"EVMChainConfigs"`
	BitcoinConfig   *BTCConfig           `json:"BitcoinConfig"`
}

func NewConfig() *Config {
	return &Config{
		cfgLock: &sync.RWMutex{},
	}
}

func (c *Config) String() string {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	s, _ := json.MarshalIndent(c, "", "\t")
	return string(s)
}

func (c *Config) GetKeygen() ostypes.Keygen {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	copiedPubkeys := make([]string, len(c.Keygen.GranteePubkeys))
	copy(copiedPubkeys, c.Keygen.GranteePubkeys)

	return ostypes.Keygen{
		Status:         c.Keygen.Status,
		GranteePubkeys: copiedPubkeys,
		BlockNumber:    c.Keygen.BlockNumber,
	}
}

func (c *Config) GetChainsEnabled() []common.Chain {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	copiedChains := make([]common.Chain, len(c.ChainsEnabled))
	copy(copiedChains, c.ChainsEnabled)
	return copiedChains
}

func (c *Config) GetEVMConfig(chainID int64) (EVMConfig, bool) {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	evmCfg, found := c.EVMChainConfigs[chainID]
	return *evmCfg, found
}

func (c *Config) GetAllEVMConfigs() map[int64]*EVMConfig {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()

	// deep copy evm configs
	copied := make(map[int64]*EVMConfig, len(c.EVMChainConfigs))
	for chainID, evmConfig := range c.EVMChainConfigs {
		copied[chainID] = &EVMConfig{}
		*copied[chainID] = *evmConfig
	}
	return copied
}

func (c *Config) GetBTCConfig() (common.Chain, BTCConfig, bool) {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()

	if c.BitcoinConfig == nil { // bitcoin is not enabled
		return common.Chain{}, BTCConfig{}, false
	}
	chain := common.GetChainFromChainID(c.BitcoinConfig.ChainId)
	if chain == nil {
		panic(fmt.Sprintf("BTCChain is missing for chainID %d", c.BitcoinConfig.ChainId))
	}
	return *chain, *c.BitcoinConfig, true
}

// This is the ONLY function that writes to core params
func (c *Config) UpdateCoreParams(keygen *ostypes.Keygen, newChains []common.Chain, evmCoreParams map[int64]*ostypes.CoreParams, btcCoreParams *ostypes.CoreParams, init bool) {
	c.cfgLock.Lock()
	defer c.cfgLock.Unlock()

	// Ignore whatever order zetacore organizes chains list in state
	sort.SliceStable(newChains, func(i, j int) bool {
		return newChains[i].ChainId < newChains[j].ChainId
	})
	if len(newChains) == 0 {
		panic("No chains enabled in ZeroCore")
	}

	// Put some limitations on core params updater for now
	if !init {
		if len(c.ChainsEnabled) != len(newChains) {
			panic(fmt.Sprintf("ChainsEnabled changed at runtime!! current: %v, new: %v", c.ChainsEnabled, newChains))
		}
		for i, chain := range newChains {
			if chain != c.ChainsEnabled[i] {
				panic(fmt.Sprintf("ChainsEnabled changed at runtime!! current: %v, new: %v", c.ChainsEnabled, newChains))
			}
		}
		for _, params := range evmCoreParams {
			curCfg, found := c.EVMChainConfigs[params.ChainId]
			if !found {
				panic(fmt.Sprintf("Unreachable code: EVMConfig not found for chainID %d", params.ChainId))
			}
			if curCfg.ZetaTokenContractAddress != params.ZetaTokenContractAddress ||
				curCfg.ConnectorContractAddress != params.ConnectorContractAddress ||
				curCfg.Erc20CustodyContractAddress != params.Erc20CustodyContractAddress {
				panic(fmt.Sprintf("Zetacore contract changed at runtime!! current cfg: %v, new cfg: %v", curCfg, params))
			}
		}
	}
	c.Keygen = *keygen
	c.ChainsEnabled = newChains
	if c.BitcoinConfig != nil && btcCoreParams != nil { // update core params for bitcoin if it's enabled
		c.BitcoinConfig.CoreParams = *btcCoreParams
	}
	for _, params := range evmCoreParams { // update core params for evm chains
		c.EVMChainConfigs[params.ChainId].CoreParams = *params
	}
}

// Make a separate (deep) copy of the config
func (c *Config) Clone() *Config {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	copied := &Config{
		Peer:                c.Peer,
		PublicIP:            c.PublicIP,
		LogFormat:           c.LogFormat,
		LogLevel:            c.LogLevel,
		LogSampler:          c.LogSampler,
		PreParamsPath:       c.PreParamsPath,
		ChainID:             c.ChainID,
		ZetaCoreURL:         c.ZetaCoreURL,
		AuthzGranter:        c.AuthzGranter,
		AuthzHotkey:         c.AuthzHotkey,
		P2PDiagnostic:       c.P2PDiagnostic,
		ConfigUpdateTicker:  c.ConfigUpdateTicker,
		P2PDiagnosticTicker: c.P2PDiagnosticTicker,
		TssPath:             c.TssPath,
		TestTssKeysign:      c.TestTssKeysign,

		cfgLock:         &sync.RWMutex{},
		Keygen:          c.GetKeygen(),
		ChainsEnabled:   c.GetChainsEnabled(),
		EVMChainConfigs: make(map[int64]*EVMConfig, len(c.EVMChainConfigs)),
		BitcoinConfig:   nil,
	}
	// deep copy evm & btc configs
	for chainID, evmConfig := range c.EVMChainConfigs {
		copied.EVMChainConfigs[chainID] = &EVMConfig{}
		*copied.EVMChainConfigs[chainID] = *evmConfig
	}
	if c.BitcoinConfig != nil {
		copied.BitcoinConfig = &BTCConfig{}
		*copied.BitcoinConfig = *c.BitcoinConfig
	}

	return copied
}
