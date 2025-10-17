package keys

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/crypto"
	zetaos "github.com/zeta-chain/node/pkg/os"
)

// RelayerKey is the structure that holds the relayer private key
type RelayerKey struct {
	PrivateKey string `json:"private_key"`
}

// ResolveAddress returns the network name and address of the relayer key
func (rk RelayerKey) ResolveAddress(network chains.Network) (string, string, error) {
	var address string

	switch network {
	case chains.Network_solana:
		privKey, err := solana.PrivateKeyFromBase58(rk.PrivateKey)
		if err != nil {
			return "", "", errors.Wrap(err, "unable to construct solana private key")
		}
		address = privKey.PublicKey().String()
	default:
		return "", "", errors.Errorf("unsupported network %d: unable to derive relayer address", network)
	}

	// return network name and address
	return network.String(), address, nil
}

// LoadRelayerKey loads the relayer key for given network and password.
// Note: returns (nil,nil) if the relayer key is not present.
func LoadRelayerKey(relayerKeyPath string, network chains.Network, password string) (*RelayerKey, error) {
	// resolve the relayer key file name
	fileName, err := ResolveRelayerKeyFile(relayerKeyPath, network)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve relayer key file name")
	}

	// relayer key is optional, so it's okay if the relayer key is not provided
	if fileName == "" {
		log.Logger.Warn().Msg("blank relayer key file")
		return nil, nil
	}

	// still returns no error if the relayer key that was provided is invalid
	if !zetaos.FileExists(fileName) {
		log.Logger.Warn().Str("file_name", fileName).Msg("invalid relayer key file")
		return nil, nil
	}

	// read the relayer key file
	relayerKey, err := ReadRelayerKeyFromFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read relayer key file: %s", fileName)
	}

	// password must be set by operator
	if password == "" {
		return nil, errors.New("password is required to decrypt the private key")
	}

	// decrypt the private key
	privateKey, err := crypto.DecryptAES256GCMBase64(relayerKey.PrivateKey, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt private key")
	}

	relayerKey.PrivateKey = privateKey

	return relayerKey, nil
}

// ResolveRelayerKeyFile is a helper function to resolve the relayer key file with full path
func ResolveRelayerKeyFile(relayerKeyPath string, network chains.Network) (string, error) {
	// resolve relayer key path if it contains a tilde
	keyPath, err := zetaos.ExpandHomeDir(relayerKeyPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve relayer key path")
	}

	// get relayer key file name by network
	name, err := relayerKeyFileByNetwork(network)
	if err != nil {
		return "", errors.Wrap(err, "failed to get relayer key file name")
	}

	return filepath.Join(keyPath, name), nil
}

// WriteRelayerKeyToFile writes the relayer key to a file
func WriteRelayerKeyToFile(fileName string, relayerKey RelayerKey) error {
	keyData, err := json.Marshal(relayerKey)
	if err != nil {
		return errors.Wrap(err, "failed to marshal relayer key")
	}

	// create relay key file (owner `rw` permissions)
	return os.WriteFile(fileName, keyData, 0o600)
}

// ReadRelayerKeyFromFile reads the relayer key file and returns the key
func ReadRelayerKeyFromFile(fileName string) (*RelayerKey, error) {
	// expand home directory in the file path if it exists
	fileNameFull, err := zetaos.ExpandHomeDir(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "ExpandHome failed for file: %s", fileName)
	}

	// read the file contents
	// #nosec G304 -- relayer key file is controlled by the operator
	fileData, err := os.ReadFile(fileNameFull)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read relayer key data: %s", fileNameFull)
	}

	// unmarshal the JSON data into the struct
	var key RelayerKey
	err = json.Unmarshal(fileData, &key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal relayer key")
	}

	return &key, nil
}

// IsRelayerPrivateKeyValid checks if the given private key is valid for the given network
func IsRelayerPrivateKeyValid(privateKey string, network chains.Network) bool {
	switch network {
	case chains.Network_solana:
		_, err := solana.PrivateKeyFromBase58(privateKey)
		if err != nil {
			return false
		}
	default:
		// unsupported network
		return false
	}
	return true
}

// relayerKeyFileByNetwork returns the relayer key JSON file name based on network
func relayerKeyFileByNetwork(network chains.Network) (string, error) {
	// JSONFileSuffix is the suffix for the relayer key file
	const JSONFileSuffix = ".json"

	// return file name for supported networks only
	switch network {
	case chains.Network_solana:
		// return network name + '.json'
		return network.String() + JSONFileSuffix, nil
	default:
		return "", errors.Errorf("network %d does not support relayer key", network)
	}
}
