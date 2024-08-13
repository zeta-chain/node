package keys

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/crypto"
	zetaos "github.com/zeta-chain/zetacore/pkg/os"
)

// RelayerKey is the structure that holds the relayer private key
type RelayerKey struct {
	PrivateKey string `json:"private_key"`
}

// ResolveAddress returns the network name and address of the relayer key
func (rk RelayerKey) ResolveAddress(network chains.Network) (string, string, error) {
	// get network name
	networkName, found := chains.GetNetworkName(int32(network))
	if !found {
		return "", "", errors.Errorf("network name not found for network %d", network)
	}

	switch network {
	case chains.Network_solana:
		privKey, err := crypto.SolanaPrivateKeyFromString(rk.PrivateKey)
		if err != nil {
			return "", "", errors.Wrap(err, "unable to construct solana private key")
		}
		return networkName, privKey.PublicKey().String(), nil
	default:
		return "", "", errors.Errorf("cannot derive relayer address for unsupported network %d", network)
	}
}

// LoadRelayerKey loads the relayer key for given network and password
func LoadRelayerKey(relayerKeyPath string, network chains.Network, password string) (*RelayerKey, error) {
	// resolve the relayer key file name
	fileName, err := ResolveRelayerKeyFile(relayerKeyPath, network)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve relayer key file name")
	}

	// load the relayer key if it is present
	if zetaos.FileExists(fileName) {
		// read the relayer key file
		relayerKey, err := ReadRelayerKeyFromFile(fileName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read relayer key file: %s", fileName)
		}

		// decrypt the private key
		privateKey, err := crypto.DecryptAES256GCMBase64(relayerKey.PrivateKey, password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt private key")
		}

		relayerKey.PrivateKey = privateKey
		return relayerKey, nil
	}

	// relayer key is optional, so it's okay if the relayer key is not provided
	return nil, nil
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

	// open the file
	file, err := os.Open(fileNameFull)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open relayer key file: %s", fileNameFull)
	}
	defer file.Close()

	// read the file contents
	fileData, err := io.ReadAll(file)
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

// relayerKeyFileByNetwork returns the relayer key JSON file name based on network
func relayerKeyFileByNetwork(network chains.Network) (string, error) {
	// get network name
	networkName, found := chains.GetNetworkName(int32(network))
	if !found {
		return "", errors.Errorf("network name not found for network %d", network)
	}

	// JSONFileSuffix is the suffix for the relayer key file
	const JSONFileSuffix = ".json"

	// return file name for supported networks only
	switch network {
	case chains.Network_solana:
		return networkName + JSONFileSuffix, nil
	default:
		return "", errors.Errorf("network %d does not support relayer key", network)
	}
}
