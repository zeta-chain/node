package config

import (
	"encoding/json"

	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

const (
	FlagConfig         = "config"
	defaultCfgFileName = "zetatool_config.json"
	ZetaChainGRPC      = "127.0.0.1:9090"
	EthRPCURL          = "http://127.0.0.1:8545"
	BtcRPC             = "smoketest"
	BtcRPCPassword     = "123"
	BtcRPCHost         = "127.0.0.1:18443"
	BtcRPCParams       = "regtest"
	SolanaRPC          = "http://127.0.0.1:8899"
	ZetaChainID        = 101
)

func TestnetConfig() *Config {
	return &Config{
		ZetaGRPC:     "zetachain-testnet-grpc.itrocket.net:443",
		EthRPCURL:    "https://ethereum-sepolia-rpc.publicnode.com",
		ZetaChainID:  101,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcRPCParams: "",
		SolanaRPC:    "",
	}
}

func MainnetConfig() *Config {
	return &Config{
		ZetaGRPC:     "https://zetachain-grpc.f5nodes.com:9090",
		EthRPCURL:    "",
		ZetaChainID:  7001,
		BtcUser:      "",
		BtcPassword:  "",
		BtcHost:      "",
		BtcRPCParams: "",
		SolanaRPC:    "",
	}
}

func LocalNetConfig() *Config {
	return DefaultConfig()
}

// Config is a struct the defines the configuration fields used by zetatool
type Config struct {
	ZetaGRPC     string
	ZetaChainID  int64
	EthRPCURL    string
	BtcUser      string
	BtcPassword  string
	BtcHost      string
	BtcRPCParams string
	SolanaRPC    string
}

func DefaultConfig() *Config {
	return &Config{
		ZetaGRPC:     ZetaChainGRPC,
		EthRPCURL:    EthRPCURL,
		ZetaChainID:  ZetaChainID,
		BtcUser:      BtcRPC,
		BtcPassword:  BtcRPCPassword,
		BtcHost:      BtcRPCHost,
		BtcRPCParams: BtcRPCParams,
		SolanaRPC:    SolanaRPC,
	}
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

func GetConfig(filename string) (*Config, error) {
	//Check if cfgFile is empty, if so return default Config and save to file
	if filename == "" {
		cfg := TestnetConfig()
		err := cfg.Save()
		return cfg, err
	}

	//if file is specified, open file and return struct
	cfg := &Config{}
	err := cfg.Read(filename)
	return cfg, err
}
