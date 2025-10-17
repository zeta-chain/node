package config

import (
	"encoding/json"
	"maps"
	"strings"
	"sync"

	"github.com/showa-93/go-mask"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/mode"
)

// KeyringBackend is the type of keyring backend to use for the hotkey
type KeyringBackend string

const (
	// KeyringBackendUndefined is undefined keyring backend
	KeyringBackendUndefined KeyringBackend = ""

	// KeyringBackendTest is the test Cosmos keyring backend
	KeyringBackendTest KeyringBackend = "test"

	// KeyringBackendFile is the file Cosmos keyring backend
	KeyringBackendFile KeyringBackend = "file"

	DefaultRelayerDir = "relayer-keys"

	// DefaultRelayerKeyPath is the default path that relayer keys are stored
	DefaultRelayerKeyPath = "~/.zetacored/" + DefaultRelayerDir

	// DefaultMempoolCongestionThreshold is the default threshold of unconfirmed txs in zetacore
	// mempool to consider it congested.
	// Leave 30% of mempool space to allow txs get processed, otherwise the congestion may get
	// even worse.
	DefaultMempoolCongestionThreshold = constant.DefaultAppMempoolSize * 7 / 10
)

// ClientConfiguration is a subset of zetaclient config that is used by zetacore client
type ClientConfiguration struct {
	ChainHost       string `json:"chain_host"        mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc"         mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name"       mapstructure:"signer_name"`
	SignerPasswd    string `json:"signer_passwd"`
}

// EVMConfig is the config for EVM chain
type EVMConfig struct {
	Endpoint string `mask:"filled"`
}

// BTCConfig is the config for Bitcoin chain
type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string `mask:"filled"`
	RPCPassword string `mask:"filled"`
	RPCHost     string `mask:"filled"`
	RPCParams   string // "regtest", "mainnet", "testnet3" , "signet", "testnet4"
}

// SolanaConfig is the config for Solana chain
type SolanaConfig struct {
	Endpoint string `mask:"filled"`
}

// SuiConfig is the config for Sui chain
type SuiConfig struct {
	Endpoint string `mask:"filled"`
}

// TONConfig is the config for TON chain
type TONConfig struct {
	// Endpoint url (toncenter V2 api) e.g. https://toncenter.com/api/v2/
	Endpoint string `mask:"filled"`
}

// ComplianceConfig is the config for compliance
type ComplianceConfig struct {
	LogPath string `json:"LogPath"`
	// Deprecated: use the separate restricted addresses config
	RestrictedAddresses []string `json:"RestrictedAddresses" mask:"zero"`
}

// FeatureFlags contains feature flags for controlling new and experimental features
type FeatureFlags struct {
	// EnableMultipleCalls enables multiple calls from the same transaction
	EnableMultipleCalls bool `json:"EnableMultipleCalls"`

	// EnableSolanaAddressLookupTable enables using Solana Address Lookup Table for withdraw and call
	EnableSolanaAddressLookupTable bool `json:"EnableSolanaAddressLookupTable"`
}

// Config is the config for ZetaClient
// TODO: use snake case for json fields
// https://github.com/zeta-chain/node/issues/1020
type Config struct {
	ClientMode mode.ClientMode `json:"ClientMode"`

	Peer                    string         `json:"Peer"`
	PublicIP                string         `json:"PublicIP"`
	LogFormat               string         `json:"LogFormat"`
	LogLevel                int8           `json:"LogLevel"`
	LogSampler              bool           `json:"LogSampler"`
	PreParamsPath           string         `json:"PreParamsPath"`
	ZetaCoreHome            string         `json:"ZetaCoreHome"`
	ChainID                 string         `json:"ChainID"`
	ZetaCoreURL             string         `json:"ZetaCoreURL"`
	AuthzGranter            string         `json:"AuthzGranter"`
	AuthzHotkey             string         `json:"AuthzHotkey"`
	P2PDiagnostic           bool           `json:"P2PDiagnostic"`
	ConfigUpdateTicker      uint64         `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker     uint64         `json:"P2PDiagnosticTicker"`
	TssPath                 string         `json:"TssPath"`
	TSSMaxPendingSignatures uint64         `json:"TSSMaxPendingSignatures"`
	TestTssKeysign          bool           `json:"TestTssKeysign"`
	KeyringBackend          KeyringBackend `json:"KeyringBackend"`
	RelayerKeyPath          string         `json:"RelayerKeyPath"`

	// MaxBaseFee is the maximum base fee allowed for zetaclient to send ZetaChain transactions
	MaxBaseFee int64 `json:"MaxBaseFee"`

	// MempoolCongestionThreshold is the threshold number of unconfirmed txs in the zetacore
	// mempool to consider it congested
	//
	// Observation will stop if the number of unconfirmed txs in mempool is greater than to this
	// threshold.
	MempoolCongestionThreshold int64 `json:"MempoolCongestionThreshold"`

	// chain configs
	EVMChainConfigs map[int64]EVMConfig `json:"EVMChainConfigs"`
	BTCChainConfigs map[int64]BTCConfig `json:"BTCChainConfigs"`
	SolanaConfig    SolanaConfig        `json:"SolanaConfig"`
	SuiConfig       SuiConfig           `json:"SuiConfig"`
	TONConfig       TONConfig           `json:"TONConfig"`

	// compliance config
	ComplianceConfig ComplianceConfig `json:"ComplianceConfig"`

	// feature flags for controlling new and experimental features
	FeatureFlags FeatureFlags `json:"FeatureFlags"`

	mu *sync.RWMutex
}

// GetEVMConfig returns the EVM config for the given chain ID
func (c Config) GetEVMConfig(chainID int64) (EVMConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	evmCfg := c.EVMChainConfigs[chainID]
	return evmCfg, !evmCfg.Empty()
}

// GetAllEVMConfigs returns a map of all EVM configs
func (c Config) GetAllEVMConfigs() map[int64]EVMConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// shallow copy evm configs (sufficient for current struct with immutable fields)
	return maps.Clone(c.EVMChainConfigs)
}

// GetBTCConfig returns the BTC config for the given chain ID
func (c Config) GetBTCConfig(chainID int64) (BTCConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	btcCfg := c.BTCChainConfigs[chainID]

	return btcCfg, !btcCfg.Empty()
}

// GetSolanaConfig returns the Solana config
func (c Config) GetSolanaConfig() (SolanaConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SolanaConfig, c.SolanaConfig != (SolanaConfig{})
}

// GetSuiConfig returns the Sui config
func (c Config) GetSuiConfig() (SuiConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SuiConfig, c.SuiConfig != (SuiConfig{})
}

// GetTONConfig returns the TONConfig and a bool indicating if it's present.
func (c Config) GetTONConfig() (TONConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.TONConfig, c.TONConfig != TONConfig{}
}

// StringMasked returns the string representation of the config with sensitive fields masked.
// Currently only the endpoints and bitcoin credentials are masked.
func (c Config) StringMasked() string {
	// create a masker
	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskAnyFunc(mask.MaskTypeFilled, masker.MaskZero)

	// mask the config
	masked, err := masker.Mask(c)
	if err != nil {
		return ""
	}

	s, err := json.MarshalIndent(masked, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
}

// GetRestrictedAddressBook returns a map of restricted addresses
// Note: the restricted address book contains both ETH and BTC addresses
func (c Config) GetRestrictedAddressBook() map[string]bool {
	restrictedAddresses := make(map[string]bool)
	for _, address := range c.ComplianceConfig.RestrictedAddresses {
		if address != "" {
			restrictedAddresses[strings.ToLower(address)] = true
		}
	}
	return restrictedAddresses
}

// GetKeyringBackend returns the keyring backend
func (c Config) GetKeyringBackend() KeyringBackend {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.KeyringBackend
}

// GetRelayerKeyPath returns the relayer key path
func (c Config) GetRelayerKeyPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// use default path if not configured
	if c.RelayerKeyPath == "" {
		return DefaultRelayerKeyPath
	}
	return c.RelayerKeyPath
}

// GetMaxBaseFee returns the max base fee
func (c Config) GetMaxBaseFee() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MaxBaseFee
}

// GetMempoolCongestionThreshold returns the threshold of unconfirmed txs in zetacore mempool to
// consider it congested
func (c Config) GetMempoolCongestionThreshold() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MempoolCongestionThreshold
}

func (c EVMConfig) Empty() bool {
	return c.Endpoint == ""
}

func (c BTCConfig) Empty() bool {
	return c.RPCHost == ""
}

// GetFeatureFlags returns the feature flags
func (c Config) GetFeatureFlags() FeatureFlags {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FeatureFlags
}

// IsEnableMultipleCallsEnabled returns true if multiple calls from same transaction are enabled
func (c Config) IsEnableMultipleCallsEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FeatureFlags.EnableMultipleCalls
}

// IsEnableSolanaAddressLookupTable returns true if Solana Address Lookup Table is enabled for withdraw and call
func (c Config) IsEnableSolanaAddressLookupTable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FeatureFlags.EnableSolanaAddressLookupTable
}
