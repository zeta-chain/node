package config

import (
	"encoding/json"

	"github.com/spf13/afero"
	"github.com/zeta-chain/node/pkg/chains"
)

var AppFs = afero.NewOsFs()

const (
	FlagConfig         = "config"
	defaultCfgFileName = "zetatool_config.json"
)

func TestnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "https://zetachain-testnet-grpc.itrocket.net:443",
		EthereumRPC:  "https://ethereum-sepolia-rpc.publicnode.com",
		ZetaChainID:  101,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "",
		BscRPC:       "https://bsc-testnet-rpc.publicnode.com",
		PolygonRPC:   "https://polygon-amoy.gateway.tenderly.co",
		BaseRPC:      "https://base-sepolia-rpc.publicnode.com",
	}
}

func DevnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "",
		EthereumRPC:  "",
		ZetaChainID:  101,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "",
		BscRPC:       "",
		PolygonRPC:   "",
		BaseRPC:      "",
	}
}

func MainnetConfig() *Config {
	return &Config{
		ZetaChainRPC: "https://zetachain-mainnet.g.allthatnode.com:443/archive/tendermint",
		EthereumRPC:  "https://eth-mainnet.public.blastapi.io",
		ZetaChainID:  7000,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcParams:    "",
		SolanaRPC:    "",
		BaseRPC:      "https://base-mainnet.public.blastapi.io",
		BscRPC:       "https://bsc-mainnet.public.blastapi.io",
		PolygonRPC:   "https://polygon-bor-rpc.publicnode.com",
	}
}

func PrivateNetConfig() *Config {
	return &Config{
		ZetaChainRPC: "http://127.0.0.1:26657",
		EthereumRPC:  "http://127.0.0.1:8545",
		ZetaChainID:  101,
		BtcUser:      "smoketest",
		BtcPassword:  "123",
		BtcHost:      "127.0.0.1:18443",
		BtcParams:    "regtest",
		SolanaRPC:    "http://127.0.0.1:8899",
	}
}

// Config is a struct the defines the configuration fields used by zetatool
type Config struct {
	ZetaChainRPC string
	ZetaChainID  int64
	EthereumRPC  string
	BtcUser      string
	BtcPassword  string
	BtcHost      string
	BtcParams    string
	SolanaRPC    string
	BscRPC       string
	PolygonRPC   string
	BaseRPC      string
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
	data, err := afero.ReadFile(AppFs, filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	return err
}

func GetConfig(chain chains.Chain, filename string) (*Config, error) {
	//Check if cfgFile is empty, if so return default Config and save to file
	if filename == "" {
		return map[chains.NetworkType]*Config{
			chains.NetworkType_mainnet: MainnetConfig(),
			chains.NetworkType_testnet: TestnetConfig(),
			chains.NetworkType_privnet: PrivateNetConfig(),
			chains.NetworkType_devnet:  DevnetConfig(),
		}[chain.NetworkType], nil
	}

	//if file is specified, open file and return struct
	cfg := &Config{}
	err := cfg.Read(filename)
	return cfg, err
}
