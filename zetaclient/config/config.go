package config

import (
	"encoding/json"
	"github.com/zeta-chain/zetacore/cmd"
	"os"
	"path/filepath"
)

const filename string = "zeta-client.json"

func Save(config *Config, path string) error {
	file := filepath.Join(path, filename)
	file = filepath.Clean(file)
	//fp, err := os.Create(file)
	//if err != nil {
	//	// failed to create/open the file
	//	return err
	//}
	//if err := toml.NewEncoder(fp).Encode(config); err != nil {
	//	// failed to encode
	//	return err
	//}

	jsonFile, _ := json.MarshalIndent(config, "", "    ")
	err := os.WriteFile(file, jsonFile, 0600)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) (*Config, error) {
	file := filepath.Join(path, filename)
	file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	file = filepath.Clean(file)
	cfg := &Config{}
	input, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &cfg)
	if err != nil {
		return nil, err
	}
	// Initialize Global config variables
	cmd.CHAINID = cfg.ChainID
	return cfg, nil
}
