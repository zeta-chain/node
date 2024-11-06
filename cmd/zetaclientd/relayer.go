package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/crypto"
	zetaos "github.com/zeta-chain/node/pkg/os"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
)

// relayerOptions is the struct that holds arguments for the relayer commands
type relayerOptions struct {
	privateKey     string
	network        int32
	password       string
	relayerKeyPath string
}

var relayerOpts relayerOptions

func setupRelayerOptions() {
	f, cfg := RelayerCmd.PersistentFlags(), &relayerOpts

	// resolve default relayer key path
	defaultKeyPath := fmt.Sprintf("%s/%s", app.DefaultNodeHome, config.DefaultRelayerDir)

	f.Int32Var(&cfg.network, "network", 7, "network id, (7:solana)")
	f.StringVar(&cfg.password, "password", "", "the password to decrypt the relayer private key")
	f.StringVar(&cfg.relayerKeyPath, "key-path", defaultKeyPath, "path to relayer keys")

	// import command in addition has the private key option
	f = RelayerImportKeyCmd.Flags()
	f.StringVar(&cfg.privateKey, "private-key", "", "the relayer private key to import")
}

// RelayerShowAddress shows the relayer address
func RelayerShowAddress(_ *cobra.Command, _ []string) error {
	// try loading the relayer key if present
	network := chains.Network(relayerOpts.network)
	relayerKey, err := keys.LoadRelayerKey(relayerOpts.relayerKeyPath, network, relayerOpts.password)
	if err != nil {
		return errors.Wrap(err, "failed to load relayer key")
	}

	// relayer key does not exist, return error
	if relayerKey == nil {
		return fmt.Errorf(
			"relayer key not found for network %d in path: %s",
			relayerOpts.network,
			relayerOpts.relayerKeyPath,
		)
	}

	// resolve the relayer address
	networkName, address, err := relayerKey.ResolveAddress(network)
	if err != nil {
		return errors.Wrap(err, "failed to resolve relayer address")
	}

	fmt.Printf("relayer address (%s): %s\n", networkName, address)

	return nil
}

// RelayerImportKey imports a relayer private key
func RelayerImportKey(_ *cobra.Command, _ []string) error {
	// validate private key and password
	switch {
	case relayerOpts.privateKey == "":
		return errors.New("must provide a private key")
	case relayerOpts.password == "":
		return errors.New("must provide a password")
	case !keys.IsRelayerPrivateKeyValid(relayerOpts.privateKey, chains.Network(relayerOpts.network)):
		return errors.New("invalid private key")
	}

	// resolve the relayer key file path
	fileName, err := keys.ResolveRelayerKeyFile(relayerOpts.relayerKeyPath, chains.Network(relayerOpts.network))
	if err != nil {
		return errors.Wrap(err, "failed to resolve relayer key file path")
	}

	// create path (owner `rwx` permissions) if it does not exist
	keyPath := filepath.Dir(fileName)
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		if err := os.MkdirAll(keyPath, 0o700); err != nil {
			return errors.Wrapf(err, "failed to create relayer key path: %s", keyPath)
		}
	}

	// avoid overwriting existing key file
	if zetaos.FileExists(fileName) {
		return errors.Errorf(
			"relayer key %s already exists, please backup and remove it before importing a new key",
			fileName,
		)
	}

	// encrypt the private key
	ciphertext, err := crypto.EncryptAES256GCMBase64(relayerOpts.privateKey, relayerOpts.password)
	if err != nil {
		return errors.Wrap(err, "private key encryption failed")
	}

	// create the relayer key file
	err = keys.WriteRelayerKeyToFile(fileName, keys.RelayerKey{PrivateKey: ciphertext})
	if err != nil {
		return errors.Wrapf(err, "failed to create relayer key file: %s", fileName)
	}
	fmt.Printf("successfully imported relayer key: %s\n", fileName)

	return nil
}
