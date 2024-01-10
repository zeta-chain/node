package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains the configuration for the smoke test
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
	Zevm         string `yaml:"zevm"`
	EVM          string `yaml:"evm"`
	Bitcoin      string `yaml:"bitcoin"`
	ZetaCoreGRPC string `yaml:"zetacore_grpc"`
	ZetaCoreRPC  string `yaml:"zetacore_rpc"`
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
	ZEVMSwapAppAddr    string `yaml:"zevm_swap_app"`
	ContextAppAddr     string `yaml:"context_app"`
	TestDappAddr       string `yaml:"test_dapp"`
}

func DefaultConfig() Config {
	return Config{
		RPCs: RPCs{
			Zevm:         "http://zetacore0:8545",
			EVM:          "http://eth:8545",
			Bitcoin:      "bitcoin:18443",
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
