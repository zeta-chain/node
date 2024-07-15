// Package config provides functions to load and save ZetaClient config
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// restrictedAddressBook is a map of restricted addresses
var restrictedAddressBook = map[string]bool{}

// filename is config file name for ZetaClient
const filename string = "zetaclient_config.json"

// folder is the folder name for ZetaClient config
const folder string = "config"

// Save saves ZetaClient config
func Save(config *Config, path string) error {
	folderPath := filepath.Join(path, folder)
	err := os.MkdirAll(folderPath, 0o750)
	if err != nil {
		return err
	}

	file := filepath.Join(path, folder, filename)
	file = filepath.Clean(file)

	jsonFile, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, jsonFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Load loads ZetaClient config from a filepath
func Load(path string) (Config, error) {
	// retrieve file
	file := filepath.Join(path, folder, filename)
	file, err := filepath.Abs(file)
	if err != nil {
		return Config{}, err
	}
	file = filepath.Clean(file)

	// read config
	cfg := New(false)
	input, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(input, &cfg)
	if err != nil {
		return Config{}, err
	}

	// read keyring backend and use test by default
	if cfg.KeyringBackend == KeyringBackendUndefined {
		cfg.KeyringBackend = KeyringBackendTest
	}
	if cfg.KeyringBackend != KeyringBackendFile && cfg.KeyringBackend != KeyringBackendTest {
		return Config{}, fmt.Errorf("invalid keyring backend %s", cfg.KeyringBackend)
	}

	// fields sanitization
	cfg.TssPath = GetPath(cfg.TssPath)
	cfg.PreParamsPath = GetPath(cfg.PreParamsPath)
	cfg.ZetaCoreHome = path

	// load compliance config
	LoadComplianceConfig(cfg)

	return cfg, nil
}

// LoadComplianceConfig loads compliance data (restricted addresses) from config
func LoadComplianceConfig(cfg Config) {
	restrictedAddressBook = cfg.GetRestrictedAddressBook()
}

// GetPath returns the absolute path of the input path
func GetPath(inputPath string) string {
	path := strings.Split(inputPath, "/")
	if len(path) > 0 {
		if path[0] == "~" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			path[0] = home
			return filepath.Join(path...)
		}
	}

	return inputPath
}

// ContainRestrictedAddress returns true if any one of the addresses is restricted
// Note: the addrs can contains both ETH and BTC addresses
func ContainRestrictedAddress(addrs ...string) bool {
	for _, addr := range addrs {
		if addr != "" && restrictedAddressBook[strings.ToLower(addr)] {
			return true
		}
	}
	return false
}
