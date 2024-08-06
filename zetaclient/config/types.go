package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"

	"github.com/zeta-chain/zetacore/pkg/chains"
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
)

// ClientConfiguration is a subset of zetaclient config that is used by zetacore client
type ClientConfiguration struct {
	ChainHost       string `json:"chain_host"        mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc"         mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name"       mapstructure:"signer_name"`
	SignerPasswd    string `json:"signer_passwd"`
	HsmMode         bool   `json:"hsm_mode"`
}

// EVMConfig is the config for EVM chain
type EVMConfig struct {
	Chain    chains.Chain
	Endpoint string
}

// BTCConfig is the config for Bitcoin chain
type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCHost     string
	RPCParams   string // "regtest", "mainnet", "testnet3"
}

// SolanaConfig is the config for Solana chain
type SolanaConfig struct {
	Endpoint string
}

// ComplianceConfig is the config for compliance
type ComplianceConfig struct {
	LogPath             string   `json:"LogPath"`
	RestrictedAddresses []string `json:"RestrictedAddresses"`
}

// Config is the config for ZetaClient
// TODO: use snake case for json fields
// https://github.com/zeta-chain/node/issues/1020
type Config struct {
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
	SolanaKeyFile       string         `json:"SolanaKeyFile"`

	// chain configs
	EVMChainConfigs map[int64]EVMConfig `json:"EVMChainConfigs"`
	BitcoinConfig   BTCConfig           `json:"BitcoinConfig"`
	SolanaConfig    SolanaConfig        `json:"SolanaConfig"`

	// compliance config
	ComplianceConfig ComplianceConfig `json:"ComplianceConfig"`

	mu *sync.RWMutex
}

// GetEVMConfig returns the EVM config for the given chain ID
func (c Config) GetEVMConfig(chainID int64) (EVMConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	evmCfg, found := c.EVMChainConfigs[chainID]
	return evmCfg, found
}

// GetAllEVMConfigs returns a map of all EVM configs
func (c Config) GetAllEVMConfigs() map[int64]EVMConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy evm configs
	copied := make(map[int64]EVMConfig, len(c.EVMChainConfigs))
	for chainID, evmConfig := range c.EVMChainConfigs {
		copied[chainID] = evmConfig
	}
	return copied
}

// GetBTCConfig returns the BTC config
func (c Config) GetBTCConfig() (BTCConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.BitcoinConfig, c.BitcoinConfig != (BTCConfig{})
}

// GetSolanaConfig returns the Solana config
func (c Config) GetSolanaConfig() (SolanaConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SolanaConfig, c.SolanaConfig != (SolanaConfig{})
}

// String returns the string representation of the config
func (c Config) String() string {
	s, err := json.MarshalIndent(c, "", "\t")
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

// LoadSolanaPrivateKey loads the Solana private key from the key file
func (c Config) LoadSolanaPrivateKey() (solana.PrivateKey, error) {
	// key file path
	fileName := path.Join(c.ZetaCoreHome, c.SolanaKeyFile)

	// load the gateway keypair from a JSON file
	// #nosec G304 -- user is allowed to specify the key file
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return solana.PrivateKey{}, errors.Wrapf(err, "unable to read Solana key file: %s", fileName)
	}

	// unmarshal the JSON content into a slice of bytes
	var keyBytes []byte
	err = json.Unmarshal(fileContent, &keyBytes)
	if err != nil {
		return solana.PrivateKey{}, errors.Wrap(err, "unable to unmarshal Solana key bytes")
	}

	// ensure the key length is 64 bytes
	if len(keyBytes) != 64 {
		return solana.PrivateKey{}, fmt.Errorf("invalid Solana key length: %d", len(keyBytes))
	}

	// create private key from the key bytes
	privKey := solana.PrivateKey(keyBytes)

	return privKey, nil
}

func (c EVMConfig) Empty() bool {
	return c.Endpoint == "" && c.Chain.IsEmpty()
}
