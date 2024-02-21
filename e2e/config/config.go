package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"

	"gopkg.in/yaml.v2"
)

// Config contains the configuration for the e2e test
type Config struct {
	Accounts    Accounts  `yaml:"accounts"`
	RPCs        RPCs      `yaml:"rpcs"`
	Contracts   Contracts `yaml:"contracts"`
	ZetaChainID string    `yaml:"zeta_chain_id"`
	TestList    []string  `yaml:"test_list"`
}

// Accounts contains the configuration for the accounts
type Accounts struct {
	EVMAddress string `yaml:"evm_address"`
	EVMPrivKey string `yaml:"evm_priv_key"`
}

// RPCs contains the configuration for the RPC endpoints
type RPCs struct {
	Zevm         string     `yaml:"zevm"`
	EVM          string     `yaml:"evm"`
	Bitcoin      BitcoinRPC `yaml:"bitcoin"`
	ZetaCoreGRPC string     `yaml:"zetacore_grpc"`
	ZetaCoreRPC  string     `yaml:"zetacore_rpc"`
}

// BitcoinRPC contains the configuration for the Bitcoin RPC endpoint
type BitcoinRPC struct {
	User         string             `yaml:"user"`
	Pass         string             `yaml:"pass"`
	Host         string             `yaml:"host"`
	HTTPPostMode bool               `yaml:"http_post_mode"`
	DisableTLS   bool               `yaml:"disable_tls"`
	Params       BitcoinNetworkType `yaml:"params"`
}

// Contracts contains the addresses of predeployed contracts
type Contracts struct {
	EVM  EVM  `yaml:"evm"`
	ZEVM ZEVM `yaml:"zevm"`
}

// EVM contains the addresses of predeployed contracts on the EVM chain
type EVM struct {
	ZetaEthAddress   string `yaml:"zeta_eth"`
	ConnectorEthAddr string `yaml:"connector_eth"`
	CustodyAddr      string `yaml:"custody"`
	USDT             string `yaml:"usdt"`
}

// ZEVM contains the addresses of predeployed contracts on the zEVM chain
type ZEVM struct {
	SystemContractAddr string `yaml:"system_contract"`
	ETHZRC20Addr       string `yaml:"eth_zrc20"`
	USDTZRC20Addr      string `yaml:"usdt_zrc20"`
	BTCZRC20Addr       string `yaml:"btc_zrc20"`
	UniswapFactoryAddr string `yaml:"uniswap_factory"`
	UniswapRouterAddr  string `yaml:"uniswap_router"`
	ConnectorZEVMAddr  string `yaml:"connector_zevm"`
	WZetaAddr          string `yaml:"wzeta"`
	ZEVMSwapAppAddr    string `yaml:"zevm_swap_app"`
	ContextAppAddr     string `yaml:"context_app"`
	TestDappAddr       string `yaml:"test_dapp"`
}

// DefaultConfig returns the default config using values for localnet testing
func DefaultConfig() Config {
	return Config{
		RPCs: RPCs{
			Zevm: "http://zetacore0:8545",
			EVM:  "http://eth:8545",
			Bitcoin: BitcoinRPC{
				Host:         "bitcoin:18443",
				User:         "e2e",
				Pass:         "123",
				HTTPPostMode: true,
				DisableTLS:   true,
				Params:       Regnet,
			},
			ZetaCoreGRPC: "zetacore0:9090",
			ZetaCoreRPC:  "http://zetacore0:26657",
		},
		ZetaChainID: "athens_101-1",
		Contracts: Contracts{
			EVM: EVM{
				USDT: "0xff3135df4F2775f4091b81f4c7B6359CfA07862a",
			},
		},
	}
}

// ReadConfig reads the config file
func ReadConfig(file string) (config Config, err error) {
	if file == "" {
		return Config{}, errors.New("file name cannot be empty")
	}

	// #nosec G304 -- this is a config file
	b, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}

	return
}

// WriteConfig writes the config file
func WriteConfig(file string, config Config) error {
	if file == "" {
		return errors.New("file name cannot be empty")
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, b, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates the config
func (c Config) Validate() error {
	if c.RPCs.Bitcoin.Params != Mainnet &&
		c.RPCs.Bitcoin.Params != Testnet3 &&
		c.RPCs.Bitcoin.Params != Regnet {
		return errors.New("invalid bitcoin params")
	}
	return nil
}

// BitcoinNetworkType is a custom type to represent allowed network types
type BitcoinNetworkType string

// Enum values for BitcoinNetworkType
const (
	Mainnet  BitcoinNetworkType = "mainnet"
	Testnet3 BitcoinNetworkType = "testnet3"
	Regnet   BitcoinNetworkType = "regnet"
)

// GetParams returns the chaincfg.Params for the BitcoinNetworkType
func (bnt BitcoinNetworkType) GetParams() (chaincfg.Params, error) {
	switch bnt {
	case Mainnet:
		return chaincfg.MainNetParams, nil
	case Testnet3:
		return chaincfg.TestNet3Params, nil
	case Regnet:
		return chaincfg.RegressionNetParams, nil
	default:
		return chaincfg.Params{}, fmt.Errorf("invalid bitcoin params %s", bnt)
	}
}
