// Package config provides functions to load and save ZetaClient config
package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	_ "github.com/zeta-chain/node/pkg/sdkconfig/default"

	"github.com/zeta-chain/node/pkg/chains"
)

// restrictedAddressBook is a map of restricted addresses
var restrictedAddressBook = map[string]bool{}
var restrictedAddressBookLock sync.RWMutex

const restrictedAddressesPath string = "zetaclient_restricted_addresses.json"

// filename is config file name for ZetaClient
const filename string = "zetaclient_config.json"

// folder is the folder name for ZetaClient config
const folder string = "config"

// Save saves ZetaClient config
func Save(config *Config, path string) error {
	// validate config
	if err := Validate(*config); err != nil {
		return errors.Wrapf(err, "config file validation failed")
	}

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
func Load(basePath string) (Config, error) {
	// retrieve file
	file := filepath.Join(basePath, folder, filename)
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

	// fields sanitization
	cfg.TssPath = GetPath(cfg.TssPath)
	cfg.PreParamsPath = GetPath(cfg.PreParamsPath)
	cfg.ZetaCoreHome = basePath

	// validate config
	if err := Validate(cfg); err != nil {
		return Config{}, errors.Wrapf(err, "config file validation failed")
	}

	return cfg, nil
}

// Validate performs basic validation on the config fields
func Validate(cfg Config) error {
	// go-tss requires a valid IPv4 address
	if cfg.PublicIP != "" && !govalidator.IsIPv4(cfg.PublicIP) {
		return errors.Errorf("reason: invalid public IP, got: %s", cfg.PublicIP)
	}

	if cfg.PublicDNS != "" && !govalidator.IsDNSName(cfg.PublicDNS) {
		return errors.Errorf("reason: invalid public DNS, got: %s", cfg.PublicDNS)
	}

	if _, err := chains.ZetaChainFromCosmosChainID(cfg.ChainID); err != nil {
		return errors.Errorf("reason: invalid chain id, got: %s", cfg.ChainID)
	}

	// ZetaCoreURL can be either an IP address or a hostname (e.g., Docker service name)
	if cfg.ZetaCoreURL != "" && !govalidator.IsIP(cfg.ZetaCoreURL) && !govalidator.IsDNSName(cfg.ZetaCoreURL) {
		return errors.Errorf("reason: invalid zetacore URL, got: %s", cfg.ZetaCoreURL)
	}

	// validate granter address - should be a valid bech32 address
	if _, err := sdktypes.AccAddressFromBech32(cfg.AuthzGranter); err != nil {
		return errors.Errorf("reason: invalid bech32 granter address, got: %s", cfg.AuthzGranter)
	}

	// validate grantee name - should not be empty
	if strings.TrimSpace(cfg.AuthzHotkey) == "" {
		return errors.Errorf("reason: grantee name is empty")
	}

	// acceptable log levels are: 0:debug, 1:info, 2:warn, 3:error, 4:fatal, 5:panic
	if cfg.LogLevel < 0 || cfg.LogLevel > 5 {
		return errors.Errorf("reason: log level must be between 0 and 5, got: %d", cfg.LogLevel)
	}

	if cfg.ConfigUpdateTicker == 0 {
		return errors.Errorf("reason: config update ticker is 0")
	}

	if cfg.KeyringBackend != KeyringBackendFile && cfg.KeyringBackend != KeyringBackendTest {
		return errors.Errorf("reason: invalid keyring backend, got: %s", cfg.KeyringBackend)
	}

	if cfg.MaxBaseFee < 0 {
		return errors.Errorf("reason: max base fee cannot be negative, got: %d", cfg.MaxBaseFee)
	}

	if cfg.MempoolCongestionThreshold < 0 {
		return errors.Errorf(
			"reason: mempool congestion threshold cannot be negative, got: %d",
			cfg.MempoolCongestionThreshold,
		)
	}

	return nil
}

// SetRestrictedAddressesFromConfig loads compliance data (restricted addresses) from config.
func SetRestrictedAddressesFromConfig(cfg Config) {
	restrictedAddressBook = cfg.GetRestrictedAddressBook()
}

// GetRestrictedAddresses returns a list of restricted addresses
func GetRestrictedAddresses() []string {
	restrictedAddressBookLock.RLock()
	defer restrictedAddressBookLock.RUnlock()

	addresses := []string{}
	for addr := range restrictedAddressBook {
		addresses = append(addresses, addr)
	}
	return addresses
}

func getRestrictedAddressAbsPath(basePath string) (string, error) {
	file := filepath.Join(basePath, folder, restrictedAddressesPath)
	file, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrapf(err, "absolute path conversion for %s", file)
	}
	return file, nil
}

func loadRestrictedAddressesConfig(cfg Config, file string) error {
	input, err := os.ReadFile(file) // #nosec G304
	if err != nil {
		return errors.Wrapf(err, "reading file %s", file)
	}
	addresses := []string{}
	err = json.Unmarshal(input, &addresses)
	if err != nil {
		return errors.Wrap(err, "invalid json")
	}

	restrictedAddressBookLock.Lock()
	defer restrictedAddressBookLock.Unlock()

	// Clear the existing map, load addresses from main config, then load addresses
	// from dedicated config file
	SetRestrictedAddressesFromConfig(cfg)
	for _, addr := range addresses {
		restrictedAddressBook[strings.ToLower(addr)] = true
	}
	return nil
}

// LoadRestrictedAddressesConfig loads the restricted addresses from the config file
func LoadRestrictedAddressesConfig(cfg Config, basePath string) error {
	file, err := getRestrictedAddressAbsPath(basePath)
	if err != nil {
		return errors.Wrap(err, "getting restricted address path")
	}
	return loadRestrictedAddressesConfig(cfg, file)
}

// WatchRestrictedAddressesConfig monitors the restricted addresses config file
// for changes and reloads it when necessary
func WatchRestrictedAddressesConfig(ctx context.Context, cfg Config, basePath string, logger zerolog.Logger) error {
	file, err := getRestrictedAddressAbsPath(basePath)
	if err != nil {
		return errors.Wrap(err, "getting restricted address path")
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "creating file watcher")
	}
	defer watcher.Close()

	// Watch the config directory
	// If you only watch the file, the watch will be disconnected if/when
	// the config is recreated.
	dir := filepath.Dir(file)
	err = watcher.Add(dir)
	if err != nil {
		return errors.Wrapf(err, "watching directory %s", dir)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Name != file {
				continue
			}

			// only reload on create or write
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			logger.Info().Msg("restricted addresses config updated")

			err := loadRestrictedAddressesConfig(cfg, file)
			if err != nil {
				logger.Err(err).Msg("load restricted addresses config")
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return errors.Wrap(err, "watcher error")
		}
	}
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
	restrictedAddressBookLock.RLock()
	defer restrictedAddressBookLock.RUnlock()
	for _, addr := range addrs {
		if addr != "" && restrictedAddressBook[strings.ToLower(addr)] {
			return true
		}
	}
	return false
}

// ResolveDBPath resolves the path to chain observer database
func ResolveDBPath() (string, error) {
	const dbpath = ".zetaclient/chainobserver"

	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to resolve user home directory")
	}

	return filepath.Join(userDir, dbpath), nil
}
