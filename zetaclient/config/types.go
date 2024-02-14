package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// KeyringBackend is the type of keyring backend to use for the hotkey
type KeyringBackend string

const (
	KeyringBackendUndefined KeyringBackend = ""
	KeyringBackendTest      KeyringBackend = "test"
	KeyringBackendFile      KeyringBackend = "file"
)

type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string `json:"signer_passwd"`
	HsmMode         bool   `json:"hsm_mode"`
}

type EVMConfig struct {
	Chain    common.Chain
	Endpoint string
}

type BTCConfig struct {
	ChainID int64
	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCHost     string
	RPCParams   string // "regtest", "mainnet", "testnet3"
}

// Config is the config for ZetaClient
// TODO: use snake case for json fields
// https://github.com/zeta-chain/node/issues/1020
type Config struct {
	cfgLock *sync.RWMutex `json:"-"`

	Peer                string         `json:"Peer"`
	PublicIP            string         `json:"PublicIP"`
	LogFormat           string         `json:"LogFormat"`
	LogLevel            int8           `json:"LogLevel"`
	LogSampler          bool           `json:"LogSampler"`
	PreParamsPath       string         `json:"PreParamsPath"`
	ZetaCoreHome        string         `json:"ZetaCoreHome"`
	ChainID             string         `json:"ChainID"`
	ZetaCoreURL         string         `json:"ZetaCoreURL"`
	AuthzGranter        string         `json:"AuthzGranter"`
	AuthzHotkey         string         `json:"AuthzHotkey"`
	P2PDiagnostic       bool           `json:"P2PDiagnostic"`
	ConfigUpdateTicker  uint64         `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker uint64         `json:"P2PDiagnosticTicker"`
	TssPath             string         `json:"TssPath"`
	TestTssKeysign      bool           `json:"TestTssKeysign"`
	KeyringBackend      KeyringBackend `json:"KeyringBackend"`
	HsmMode             bool           `json:"HsmMode"`
	HsmHotKey           string         `json:"HsmHotKey"`

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
	s, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
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

// TODO: get chain from params chainId, and only return config here?
func (c *Config) GetBTCConfig() (common.Chain, BTCConfig, bool) {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()

	if c.BitcoinConfig == nil { // bitcoin is not enabled
		return common.Chain{}, BTCConfig{}, false
	}
	chain := common.GetChainFromChainID(c.BitcoinConfig.ChainID)
	if chain == nil {
		panic(fmt.Sprintf("BTCChain is missing for chainID %d", c.BitcoinConfig.ChainID))
	}
	return *chain, *c.BitcoinConfig, true
}

func (c *Config) GetKeyringBackend() KeyringBackend {
	c.cfgLock.RLock()
	defer c.cfgLock.RUnlock()
	return c.KeyringBackend
}

// ValidateChainParams performs some basic checks on core params
func ValidateChainParams(chainParams *observertypes.ChainParams) error {
	if chainParams == nil {
		return fmt.Errorf("invalid chain params: nil")
	}
	chain := common.GetChainFromChainID(chainParams.ChainId)
	if chain == nil {
		return fmt.Errorf("invalid chain params: chain %d not supported", chainParams.ChainId)
	}
	if chainParams.ConfirmationCount < 1 {
		return fmt.Errorf("invalid chain params: ConfirmationCount %d", chainParams.ConfirmationCount)
	}
	// zeta chain skips the rest of the checks for now
	if chain.IsZetaChain() {
		return nil
	}

	// check tickers
	if chainParams.GasPriceTicker < 1 {
		return fmt.Errorf("invalid chain params: GasPriceTicker %d", chainParams.GasPriceTicker)
	}
	if chainParams.InTxTicker < 1 {
		return fmt.Errorf("invalid chain params: InTxTicker %d", chainParams.InTxTicker)
	}
	if chainParams.OutTxTicker < 1 {
		return fmt.Errorf("invalid chain params: OutTxTicker %d", chainParams.OutTxTicker)
	}
	if chainParams.OutboundTxScheduleInterval < 1 {
		return fmt.Errorf("invalid chain params: OutboundTxScheduleInterval %d", chainParams.OutboundTxScheduleInterval)
	}
	if chainParams.OutboundTxScheduleLookahead < 1 {
		return fmt.Errorf("invalid chain params: OutboundTxScheduleLookahead %d", chainParams.OutboundTxScheduleLookahead)
	}

	// chain type specific checks
	if common.IsBitcoinChain(chainParams.ChainId) && chainParams.WatchUtxoTicker < 1 {
		return fmt.Errorf("invalid chain params: watchUtxo ticker %d", chainParams.WatchUtxoTicker)
	}
	if common.IsEVMChain(chainParams.ChainId) {
		if !validCoreContractAddress(chainParams.ZetaTokenContractAddress) {
			return fmt.Errorf("invalid chain params: zeta token contract address %s", chainParams.ZetaTokenContractAddress)
		}
		if !validCoreContractAddress(chainParams.ConnectorContractAddress) {
			return fmt.Errorf("invalid chain params: connector contract address %s", chainParams.ConnectorContractAddress)
		}
		if !validCoreContractAddress(chainParams.Erc20CustodyContractAddress) {
			return fmt.Errorf("invalid chain params: erc20 custody contract address %s", chainParams.Erc20CustodyContractAddress)
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
