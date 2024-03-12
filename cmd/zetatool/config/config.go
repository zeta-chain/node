package config

import (
	"encoding/json"
	"os"
)

const (
	Flag               = "config"
	defaultCfgFileName = "InboundTxFilter_config.json"
	ZetaURL            = "http://46.4.15.110:1317" //http://100.71.167.102:26657
	TssAddressBTC      = "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y"
	TssAddressEVM      = "0x70e967acfcc17c3941e87562161406d41676fd83"
	BtcExplorer        = "https://blockstream.info/api/address/bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y/txs"
	EthRPC             = "https://rpc.ankr.com/eth/2da24e4a1fd28f2bec1569eceb2c38a5694b7f5c83fd24c69ae714a89a514f9b"
	ConnectorAddress   = "0x000007Cf399229b2f5A4D043F20E90C9C98B7C6a"
	CustodyAddress     = "0x0000030Ec64DF25301d8414eE5a29588C4B0dE10"
	EvmStartBlock      = 19200110
	EvmMaxRange        = 1000
)

type Config struct {
	ZetaURL          string
	TssAddressBTC    string
	TssAddressEVM    string
	BtcExplorer      string
	EthRPC           string
	ConnectorAddress string
	CustodyAddress   string
	EvmStartBlock    uint64
	EvmMaxRange      uint64
}

func DefaultConfig() *Config {
	return &Config{
		ZetaURL:          ZetaURL,
		TssAddressBTC:    TssAddressBTC,
		TssAddressEVM:    TssAddressEVM,
		BtcExplorer:      BtcExplorer,
		EthRPC:           EthRPC,
		ConnectorAddress: ConnectorAddress,
		CustodyAddress:   CustodyAddress,
		EvmStartBlock:    EvmStartBlock,
		EvmMaxRange:      EvmMaxRange,
	}
}

func (c *Config) Save() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(defaultCfgFileName, file, 0600)
	return err
}

func (c *Config) Read(filename string) error {
	data, err := os.ReadFile(filename)
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
