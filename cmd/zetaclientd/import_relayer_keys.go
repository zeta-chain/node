package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/crypto"
	zetaos "github.com/zeta-chain/node/pkg/os"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
)

var CmdImportRelayerKey = &cobra.Command{
	Use:     "import-relayer-key --network=<network> --private-key=<private-key> --password=<password> --relayer-key-path=<relayer-key-path>",
	Short:   "Import a relayer private key",
	Example: `zetaclientd import-relayer-key --network=7 --private-key=<your_private_key> --password=<your_password>`,
	RunE:    ImportRelayerKey,
}

var CmdRelayerAddress = &cobra.Command{
	Use:     "relayer-address --network=<network> --password=<password> --relayer-key-path=<relayer-key-path>",
	Short:   "Show the relayer address",
	Example: `zetaclientd relayer-address --network=7 --password=my_password`,
	RunE:    ShowRelayerAddress,
}

var importArgs = importRelayerKeyArguments{}
var addressArgs = relayerAddressArguments{}

// importRelayerKeyArguments is the struct that holds the arguments for the import command
type importRelayerKeyArguments struct {
	network        int32
	privateKey     string
	password       string
	relayerKeyPath string
}

// relayerAddressArguments is the struct that holds the arguments for the show command
type relayerAddressArguments struct {
	network        int32
	password       string
	relayerKeyPath string
}

func init() {
	RootCmd.AddCommand(CmdImportRelayerKey)
	RootCmd.AddCommand(CmdRelayerAddress)

	// resolve default relayer key path
	defaultRelayerKeyPath, err := zetaos.ExpandHomeDir(config.DefaultRelayerKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to resolve default relayer key path")
	}

	CmdImportRelayerKey.Flags().Int32Var(&importArgs.network, "network", 7, "network id, (7: solana)")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.privateKey, "private-key", "", "the relayer private key to import")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.password, "password", "", "the password to encrypt the relayer private key")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.relayerKeyPath, "relayer-key-path", defaultRelayerKeyPath, "path to relayer keys")

	CmdRelayerAddress.Flags().Int32Var(&addressArgs.network, "network", 7, "network id, (7:solana)")
	CmdRelayerAddress.Flags().
		StringVar(&addressArgs.password, "password", "", "the password to decrypt the relayer private key")
	CmdRelayerAddress.Flags().
		StringVar(&addressArgs.relayerKeyPath, "relayer-key-path", defaultRelayerKeyPath, "path to relayer keys")
}

// ImportRelayerKey imports a relayer private key
func ImportRelayerKey(_ *cobra.Command, _ []string) error {
	// validate private key and password
	if importArgs.privateKey == "" {
		return errors.New("must provide a private key")
	}
	if importArgs.password == "" {
		return errors.New("must provide a password")
	}
	if !keys.IsRelayerPrivateKeyValid(importArgs.privateKey, chains.Network(importArgs.network)) {
		return errors.New("invalid private key")
	}

	// resolve the relayer key file path
	fileName, err := keys.ResolveRelayerKeyFile(importArgs.relayerKeyPath, chains.Network(importArgs.network))
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
	ciphertext, err := crypto.EncryptAES256GCMBase64(importArgs.privateKey, importArgs.password)
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

// ShowRelayerAddress shows the relayer address
func ShowRelayerAddress(_ *cobra.Command, _ []string) error {
	// try loading the relayer key if present
	network := chains.Network(addressArgs.network)
	relayerKey, err := keys.LoadRelayerKey(addressArgs.relayerKeyPath, network, addressArgs.password)
	if err != nil {
		return errors.Wrap(err, "failed to load relayer key")
	}

	// relayer key does not exist, return error
	if relayerKey == nil {
		return fmt.Errorf(
			"relayer key not found for network %d in path: %s",
			addressArgs.network,
			addressArgs.relayerKeyPath,
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
