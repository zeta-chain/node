package config

import (
	"github.com/pelletier/go-toml"
	"github.com/zeta-chain/zetacore/cmd"
	"os"
	"path/filepath"
)

const filename string = "zeta-client.toml"

func Save(config *Config, path string) error {
	file := filepath.Join(path, filename)
	file = filepath.Clean(file)
	fp, err := os.Create(file)
	if err != nil {
		// failed to create/open the file
		return err
	}
	if err := toml.NewEncoder(fp).Encode(config); err != nil {
		// failed to encode
		return err
	}
	if err := fp.Close(); err != nil {
		// failed to close the file
		return err
	}
	return nil
}

func Load(path string) (*Config, error) {
	file := filepath.Join(path, filename)
	file = filepath.Clean(file)
	cfg := &Config{}
	fp, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	if err := toml.NewDecoder(fp).Decode(cfg); err != nil {
		return cfg, err
	}
	if err := fp.Close(); err != nil {
		// failed to close the file
		return cfg, err
	}

	// Initialize Global config variables
	BitcoinConfig = cfg.BitcoinConfig
	cmd.CHAINID = cfg.ChainID
	for _, chain := range cfg.ChainsEnabled {
		if chain.IsEVMChain() {
			cfg.EVMChainConfigs[chain.ChainName.String()].CoreParams = NewCoreParams()
		}
	}

	return cfg, nil
}
