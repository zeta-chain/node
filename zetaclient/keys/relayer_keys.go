package keys

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/pkg/chains"
	zetaos "github.com/zeta-chain/zetacore/pkg/os"
)

const (
	// RelayerKeyFileSolana is the file name for the Solana relayer key
	RelayerKeyFileSolana = "solana.json"
)

// RelayerKey is the structure that holds the relayer private key
type RelayerKey struct {
	PrivateKey string `json:"private_key"`
}

// ResolveAddress returns the network name and address of the relayer key
func (rk RelayerKey) ResolveAddress(network int32) (string, string, error) {
	// get network name
	networkName, found := chains.GetNetworkName(network)
	if !found {
		return "", "", errors.Errorf("network name not found for network %d", network)
	}

	switch chains.Network(network) {
	case chains.Network_solana:
		privKey, err := solana.PrivateKeyFromBase58(rk.PrivateKey)
		if err != nil {
			return "", "", errors.Wrap(err, "unable to construct solana private key")
		}
		return networkName, privKey.PublicKey().String(), nil
	default:
		return "", "", errors.Errorf("cannot derive relayer address for unsupported network %d", network)
	}
}

// LoadRelayerKey loads a relayer key from given path and chain
func LoadRelayerKey(keyPath string, chain chains.Chain) (RelayerKey, error) {
	// determine relayer key file name based on chain
	var fileName string
	switch chain.Network {
	case chains.Network_solana:
		fileName = path.Join(keyPath, RelayerKeyFileSolana)
	default:
		return RelayerKey{}, errors.Errorf("relayer key not supported for network %s", chain.Network)
	}

	// read the relayer key file
	relayerKey, err := ReadRelayerKeyFromFile(fileName)
	if err != nil {
		return RelayerKey{}, errors.Wrap(err, "ReadRelayerKeyFile failed")
	}

	return relayerKey, nil
}

// ReadRelayerKeyFromFile reads the relayer key file and returns the key
func ReadRelayerKeyFromFile(fileName string) (RelayerKey, error) {
	// expand home directory in the file path if it exists
	fileNameFull, err := zetaos.ExpandHomeDir(fileName)
	if err != nil {
		return RelayerKey{}, errors.Wrapf(err, "ExpandHome failed for file: %s", fileName)
	}

	// open the file
	file, err := os.Open(fileNameFull)
	if err != nil {
		return RelayerKey{}, errors.Wrapf(err, "unable to open relayer key file: %s", fileNameFull)
	}
	defer file.Close()

	// read the file contents
	fileData, err := io.ReadAll(file)
	if err != nil {
		return RelayerKey{}, errors.Wrapf(err, "unable to read relayer key data: %s", fileNameFull)
	}

	// unmarshal the JSON data into the struct
	var key RelayerKey
	err = json.Unmarshal(fileData, &key)
	if err != nil {
		return RelayerKey{}, errors.Wrap(err, "unable to unmarshal relayer key")
	}

	return key, nil
}

// GetRelayerKeyFileByNetwork returns the relayer key file name based on network
func GetRelayerKeyFileByNetwork(network int32) (string, error) {
	// get network name
	networkName, found := chains.GetNetworkName(network)
	if !found {
		return "", errors.Errorf("network name not found for network %d", network)
	}

	// return file name for supported networks only
	switch chains.Network(network) {
	case chains.Network_solana:
		return networkName + ".json", nil
	default:
		return "", errors.Errorf("network %d does not support relayer key", network)
	}
}
