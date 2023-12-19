package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// TODO: support pre-deployed addresses for zEVM contracts
// https://github.com/zeta-chain/node-private/issues/41

// Config contains the configuration for the smoke test
type Config struct {
	RPCs        RPCs      `yaml:"rpcs"`
	Contracts   Contracts `yaml:"contracts"`
	ZetaChainID string    `yaml:"zeta_chain_id"`
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
	EVM EVM `yaml:"evm"`
}

// EVM contains the addresses of predeployed contracts on the EVM chain
type EVM struct {
	ZetaEthAddress   string `yaml:"zeta_eth"`
	ConnectorEthAddr string `yaml:"connector_eth"`
	USDT             string `yaml:"usdt"`
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
