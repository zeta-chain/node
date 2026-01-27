package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"

	"github.com/zeta-chain/node/pkg/chains"
)

var AppFs = afero.NewOsFs()

const (
	FlagConfig         = "config"
	defaultCfgFileName = "zetatool_config.json"
	FlagDebug          = "debug"

	// Network name constants
	NetworkMainnet  = "mainnet"
	NetworkTestnet  = "testnet"
	NetworkSignet   = "signet"
	NetworkLocalnet = "localnet"
)

func TestnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "https://zetachain-athens.g.allthatnode.com/archive/tendermint",
		EthereumRPC:  "https://ethereum-sepolia-rpc.publicnode.com",
		ZetaChainID:  chains.ZetaChainTestnet.ChainId,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "https://api.testnet.solana.com",
		BscRPC:       "https://bsc-testnet-rpc.publicnode.com",
		PolygonRPC:   "https://rpc-amoy.polygon.technology/",
		BaseRPC:      "https://base-sepolia-rpc.publicnode.com",
		SuiRPC:       "https://fullnode.testnet.sui.io:443",
		TonRPC:       "",
		ArbitrumRPC:  "https://sepolia-rollup.arbitrum.io/rpc",
		OptimismRPC:  "https://sepolia.optimism.io",
		AvalancheRPC: "https://avalanche-fuji-c-chain-rpc.publicnode.com",
		WorldRPC:     "https://worldchain-sepolia.g.alchemy.com/public",
	}
}

func DevnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "",
		EthereumRPC:  "",
		ZetaChainID:  chains.ZetaChainDevnet.ChainId,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "",
		BscRPC:       "",
		PolygonRPC:   "",
		BaseRPC:      "",
		SuiRPC:       "",
		TonRPC:       "",
	}
}

func MainnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "https://zetachain-mainnet.g.allthatnode.com:443/archive/tendermint",
		EthereumRPC:  "https://eth-mainnet.public.blastapi.io",
		ZetaChainID:  chains.ZetaChainMainnet.ChainId,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "https://api.mainnet-beta.solana.com",
		BaseRPC:      "https://base-mainnet.public.blastapi.io",
		BscRPC:       "https://bsc-mainnet.public.blastapi.io",
		PolygonRPC:   "https://polygon-bor-rpc.publicnode.com",
		SuiRPC:       "https://fullnode.mainnet.sui.io:443",
		TonRPC:       "",
		ArbitrumRPC:  "https://arb1.arbitrum.io/rpc",
		OptimismRPC:  "https://mainnet.optimism.io",
		AvalancheRPC: "https://api.avax.network/ext/bc/C/rpc",
		WorldRPC:     "https://worldchain-mainnet.g.alchemy.com/public",
	}
}

// PrivateNetConfig returns a config for a private network, used for localnet testing
func PrivateNetConfig() *Config {
	return &Config{
		ZetaChainRPC: "http://127.0.0.1:26657",
		EthereumRPC:  "http://127.0.0.1:8545",
		ZetaChainID:  chains.ZetaChainPrivnet.ChainId,
		BtcUser:      "smoketest",
		BtcPassword:  "123",
		BtcHost:      "127.0.0.1:18443",
		BtcParams:    "regtest",
		SolanaRPC:    "http://127.0.0.1:8899",
		SuiRPC:       "http://127.0.0.1:9000",
		TonRPC:       "http://127.0.0.1:8081",
	}
}

// Config is a struct the defines the configuration fields used by zetatool
type Config struct {
	ZetaChainRPC string `json:"zeta_chain_rpc"`
	ZetaChainID  int64  `json:"zeta_chain_id"`
	EthereumRPC  string `json:"ethereum_rpc"`
	BtcUser      string `json:"btc_user"`
	BtcPassword  string `json:"btc_password"`
	BtcHost      string `json:"btc_host"`
	BtcParams    string `json:"btc_params"`
	SolanaRPC    string `json:"solana_rpc"`
	BscRPC       string `json:"bsc_rpc"`
	PolygonRPC   string `json:"polygon_rpc"`
	BaseRPC      string `json:"base_rpc"`
	SuiRPC       string `json:"sui_rpc"`
	TonRPC       string `json:"ton_rpc"`
	ArbitrumRPC  string `json:"arbitrum_rpc"`
	OptimismRPC  string `json:"optimism_rpc"`
	AvalancheRPC string `json:"avalanche_rpc"`
	WorldRPC     string `json:"world_rpc"`
}

func (c *Config) Save() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	err = afero.WriteFile(AppFs, defaultCfgFileName, file, 0600)
	return err
}
func (c *Config) Read(filename string) error {
	// #nosec G304 reading file is safe
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	return err
}

func GetConfig(chain chains.Chain, filename string) (*Config, error) {
	//Check if cfgFile is empty, if so return default Config based on network type
	if filename == "" {
		return map[chains.NetworkType]*Config{
			chains.NetworkType_mainnet: MainnetConfig(),
			chains.NetworkType_testnet: TestnetConfig(),
			chains.NetworkType_privnet: PrivateNetConfig(),
			chains.NetworkType_devnet:  DevnetConfig(),
		}[chain.NetworkType], nil
	}

	//if a file is specified, use the config in the file
	cfg := &Config{}
	err := cfg.Read(filename)
	return cfg, err
}

// GetConfigByNetwork returns a config based on network name string.
// Valid network names: "mainnet", "testnet", "localnet", "devnet"
func GetConfigByNetwork(network, filename string) (*Config, error) {
	// If a custom config file is specified, use it
	if filename != "" {
		cfg := &Config{}
		err := cfg.Read(filename)
		return cfg, err
	}

	// Return default config based on network name
	switch network {
	case "mainnet":
		return MainnetConfig(), nil
	case "testnet":
		return TestnetConfig(), nil
	case "localnet":
		return PrivateNetConfig(), nil
	case "devnet":
		return DevnetConfig(), nil
	default:
		return nil, fmt.Errorf("unknown network: %s (valid: mainnet, testnet, localnet, devnet)", network)
	}
}
