package config

import (
	"encoding/json"

	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

const (
	FlagConfig         = "config"
	defaultCfgFileName = "zetatool_config.json"
	ZetaURL            = "127.0.0.1:1317"
	BtcExplorerURL     = "https://blockstream.info/api/"
	EthRPCURL          = "https://ethereum-rpc.publicnode.com"
	ConnectorAddress   = "0x000007Cf399229b2f5A4D043F20E90C9C98B7C6a"
	CustodyAddress     = "0x0000030Ec64DF25301d8414eE5a29588C4B0dE10"
)

// Config is a struct the defines the configuration fields used by zetatool
type Config struct {
	ZetaURL          string
	BtcExplorerURL   string
	EthRPCURL        string
	EtherscanAPIkey  string
	ConnectorAddress string
	CustodyAddress   string
}

func DefaultConfig() *Config {
	return &Config{
		ZetaURL:          ZetaURL,
		BtcExplorerURL:   BtcExplorerURL,
		EthRPCURL:        EthRPCURL,
		ConnectorAddress: ConnectorAddress,
		CustodyAddress:   CustodyAddress,
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
		cfg := DefaultConfig()
		err := cfg.Save()
		return cfg, err
	}

	//if file is specified, open file and return struct
	cfg := &Config{}
	err := cfg.Read(filename)
	return cfg, err
}
