package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	cfg.TssPath = GetPath(cfg.TssPath)
	cfg.PreParamsPath = GetPath(cfg.PreParamsPath)
	return cfg, nil
}

func GetPath(inputPath string) string {
	path := strings.Split(inputPath, "/")
	if len(path) > 0 {
		if path[0] == "~" {
			home, _ := os.UserHomeDir()
			path[0] = home
		}
	}
	return filepath.Join(path...)
}
