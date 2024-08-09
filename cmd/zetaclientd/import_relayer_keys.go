package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/pkg/crypto"
	zetaos "github.com/zeta-chain/zetacore/pkg/os"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
)

var CmdImportRelayerKey = &cobra.Command{
	Use:     "import-relayer-key [network] [private-key] [password] [relayer-key-path]",
	Short:   "Import a relayer private key",
	Example: `zetaclientd import-relayer-key --network=7 --private-key=3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ --password=my_password`,
	RunE:    ImportRelayerKey,
}

var CmdRelayerAddress = &cobra.Command{
	Use:     "relayer-address [network] [password] [relayer-key-path]",
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
	defaultRelayerKeyPath := "~/.zetacored/relayer-keys"
	defaultRelayerKeyPath, err := zetaos.ExpandHomeDir(defaultRelayerKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to resolve default relayer key path")
	}

	CmdImportRelayerKey.Flags().Int32Var(&importArgs.network, "network", 7, "network id, (7: solana)")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.privateKey, "private-key", "", "the relayer private key to import")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.password, "password", "", "the password to encrypt the private key")
	CmdImportRelayerKey.Flags().
		StringVar(&importArgs.relayerKeyPath, "relayer-key-path", defaultRelayerKeyPath, "path to relayer keys")

	CmdRelayerAddress.Flags().Int32Var(&addressArgs.network, "network", 7, "network id, (7:solana)")
	CmdRelayerAddress.Flags().
		StringVar(&addressArgs.password, "password", "", "the password to decrypt the private key")
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

	// resolve the relayer key file path
	keyPath, fileName, err := resolveRelayerKeyPath(importArgs.network, importArgs.relayerKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to resolve relayer key file path")
	}

	// create path (owner `rwx` permissions) if it does not exist
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

	// construct the relayer key struct and write to file as json
	keyData, err := json.Marshal(keys.RelayerKey{PrivateKey: ciphertext})
	if err != nil {
		return errors.Wrap(err, "failed to marshal relayer key")
	}

	// create relay key file (owner `rw` permissions)
	err = os.WriteFile(fileName, keyData, 0o600)
	if err != nil {
		return errors.Wrapf(err, "failed to create relayer key file: %s", fileName)
	}
	fmt.Printf("successfully imported relayer key: %s\n", fileName)

	return nil
}

// ShowRelayerAddress shows the relayer address
func ShowRelayerAddress(_ *cobra.Command, _ []string) error {
	// resolve the relayer key file path
	_, fileName, err := resolveRelayerKeyPath(addressArgs.network, addressArgs.relayerKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to resolve relayer key file path")
	}

	// read the relayer key file
	relayerKey, err := keys.ReadRelayerKeyFromFile(fileName)
	if err != nil {
		return err
	}

	// decrypt the private key
	privateKey, err := crypto.DecryptAES256GCMBase64(relayerKey.PrivateKey, addressArgs.password)
	if err != nil {
		return errors.Wrap(err, "private key decryption failed")
	}
	relayerKey.PrivateKey = privateKey

	// resolve the address
	networkName, address, err := relayerKey.ResolveAddress(addressArgs.network)
	if err != nil {
		return errors.Wrap(err, "failed to resolve relayer address")
	}
	fmt.Printf("relayer address (%s): %s\n", networkName, address)

	return nil
}

// resolveRelayerKeyPath is a helper function to resolve the relayer key file path and name
func resolveRelayerKeyPath(network int32, relayerKeyPath string) (string, string, error) {
	// get relayer key file name by network
	name, err := keys.GetRelayerKeyFileByNetwork(network)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get relayer key file name")
	}

	// resolve relayer key path if it contains a tilde
	keyPath, err := zetaos.ExpandHomeDir(relayerKeyPath)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve relayer key path")
	}

	// build file name
	fileName := filepath.Join(keyPath, name)

	return keyPath, fileName, err
}
