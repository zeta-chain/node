package config

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

type EVMConfig struct {
	observertypes.CoreParams
	Chain    common.Chain
	Endpoint string
}

type BTCConfig struct {
	observertypes.CoreParams

	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCHost     string
	RPCParams   string // "regtest", "mainnet", "testnet3"
}

// TODO: use snake case for json fields
type Config struct {
	Peer                string        `json:"Peer"`
	PublicIP            string        `json:"PublicIP"`
	LogFormat           string        `json:"LogFormat"`
	LogLevel            zerolog.Level `json:"LogLevel"`
	LogSampler          bool          `json:"LogSampler"`
	PreParamsPath       string        `json:"PreParamsPath"`
	ZetaCoreHome        string        `json:"ZetaCoreHome"`
	ChainID             string        `json:"ChainID"`
	ZetaCoreURL         string        `json:"ZetaCoreURL"`
	AuthzGranter        string        `json:"AuthzGranter"`
	AuthzHotkey         string        `json:"AuthzHotkey"`
	P2PDiagnostic       bool          `json:"P2PDiagnostic"`
	ConfigUpdateTicker  uint64        `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker uint64        `json:"P2PDiagnosticTicker"`
	TssPath             string        `json:"TssPath"`
	TestTssKeysign      bool          `json:"TestTssKeysign"`
	CurrentTssPubkey    string        `json:"CurrentTssPubkey"`
	SignerPass          string        `json:"SignerPass"`
	ZetaCoreHome        string        `json:"ZetaCoreHome"`

	// chain specific fields are updatable at runtime and shared across threads
	cfgLock         *sync.RWMutex        `json:"-"`
	Keygen          observertypes.Keygen `json:"Keygen"`
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

func (c *Config) GetKeygen() observertypes.Keygen {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	copiedPubkeys := make([]string, len(c.Keygen.GranteePubkeys))
	copy(copiedPubkeys, c.Keygen.GranteePubkeys)

	return observertypes.Keygen{
		Status:         c.Keygen.Status,
		GranteePubkeys: copiedPubkeys,
		BlockNumber:    c.Keygen.BlockNumber,
	}
}

func (c *Config) GetEnabledChains() []common.Chain {
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
func (c *Config) UpdateCoreParams(keygen *observertypes.Keygen, newChains []common.Chain, evmCoreParams map[int64]*observertypes.CoreParams, btcCoreParams *observertypes.CoreParams, init bool, logger zerolog.Logger) {
	c.cfgLock.Lock()
	defer c.cfgLock.Unlock()

	// Ignore whatever order zetacore organizes chain list in state
	sort.SliceStable(newChains, func(i, j int) bool {
		return newChains[i].ChainId < newChains[j].ChainId
	})
	if len(newChains) == 0 {
		logger.Warn().Msg("UpdateCoreParams: No chains enabled in ZeroCore")
	}

	// Add some warnings if chain list changes at runtime
	if !init {
		if len(c.ChainsEnabled) != len(newChains) {
			logger.Warn().Msgf("UpdateCoreParams: ChainsEnabled changed at runtime!! current: %v, new: %v", c.ChainsEnabled, newChains)
		}
		for i, chain := range newChains {
			if chain != c.ChainsEnabled[i] {
				logger.Warn().Msgf("UpdateCoreParams: ChainsEnabled changed at runtime!! current: %v, new: %v", c.ChainsEnabled, newChains)
			}
		}
	}
	c.Keygen = *keygen
	c.ChainsEnabled = newChains
	if c.BitcoinConfig != nil && btcCoreParams != nil { // update core params for bitcoin if it has config in file
		c.BitcoinConfig.CoreParams = *btcCoreParams
	}
	for _, params := range evmCoreParams { // update core params for evm chains we have configs in file
		curCfg, found := c.EVMChainConfigs[params.ChainId]
		if found {
			curCfg.CoreParams = *params
		}
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
		ChainsEnabled:   c.GetEnabledChains(),
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

// ValidateCoreParams performs some basic checks on core params
func ValidateCoreParams(coreParams *observertypes.CoreParams) error {
	if coreParams == nil {
		return fmt.Errorf("invalid core params: nil")
	}
	chain := common.GetChainFromChainID(coreParams.ChainId)
	if chain == nil {
		return fmt.Errorf("invalid core params: chain %d not supported", coreParams.ChainId)
	}
	if coreParams.ConfirmationCount < 1 {
		return fmt.Errorf("invalid core params: ConfirmationCount %d", coreParams.ConfirmationCount)
	}
	// zeta chain skips the rest of the checks for now
	if chain.IsZetaChain() {
		return nil
	}

	// check tickers
	if coreParams.GasPriceTicker < 1 {
		return fmt.Errorf("invalid core params: GasPriceTicker %d", coreParams.GasPriceTicker)
	}
	if coreParams.InTxTicker < 1 {
		return fmt.Errorf("invalid core params: InTxTicker %d", coreParams.InTxTicker)
	}
	if coreParams.OutTxTicker < 1 {
		return fmt.Errorf("invalid core params: OutTxTicker %d", coreParams.OutTxTicker)
	}
	if coreParams.OutboundTxScheduleInterval < 1 {
		return fmt.Errorf("invalid core params: OutboundTxScheduleInterval %d", coreParams.OutboundTxScheduleInterval)
	}
	if coreParams.OutboundTxScheduleLookahead < 1 {
		return fmt.Errorf("invalid core params: OutboundTxScheduleLookahead %d", coreParams.OutboundTxScheduleLookahead)
	}

	// chain type specific checks
	if common.IsBitcoinChain(coreParams.ChainId) && coreParams.WatchUtxoTicker < 1 {
		return fmt.Errorf("invalid core params: watchUtxo ticker %d", coreParams.WatchUtxoTicker)
	}
	if common.IsEVMChain(coreParams.ChainId) {
		if !validCoreContractAddress(coreParams.ZetaTokenContractAddress) {
			return fmt.Errorf("invalid core params: zeta token contract address %s", coreParams.ZetaTokenContractAddress)
		}
		if !validCoreContractAddress(coreParams.ConnectorContractAddress) {
			return fmt.Errorf("invalid core params: connector contract address %s", coreParams.ConnectorContractAddress)
		}
		if !validCoreContractAddress(coreParams.Erc20CustodyContractAddress) {
			return fmt.Errorf("invalid core params: erc20 custody contract address %s", coreParams.Erc20CustodyContractAddress)
		}
	}
	return nil
}

func validCoreContractAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return ethcommon.IsHexAddress(address)
}
