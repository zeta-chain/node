package keys

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	// RelayerKeyFileSolana is the file name for the Solana relayer key
	RelayerKeyFileSolana = "solana.json"
)

// RelayerKey is the structure that holds the relayer private key
type RelayerKey struct {
	PrivateKey string `json:"private_key"`
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
	fileName = "/root/.zetacored/relayer-keys/solana.json"
	fmt.Println("Reading relayer key from file: ", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return RelayerKey{}, errors.Wrapf(err, "unable to open relayer key file: %s", fileName)
	}
	defer file.Close()

	// read the file contents
	fileData, err := io.ReadAll(file)
	if err != nil {
		return RelayerKey{}, errors.Wrapf(err, "unable to read relayer key data: %s", fileName)
	}

	// unmarshal the JSON data into the struct
	var key RelayerKey
	err = json.Unmarshal(fileData, &key)
	if err != nil {
		return RelayerKey{}, errors.Wrap(err, "unable to unmarshal relayer key")
	}

	return key, nil
}
