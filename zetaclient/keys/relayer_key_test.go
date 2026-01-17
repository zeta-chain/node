package keys_test

import (
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/keys"
)

// createRelayerKeyFile creates a relayer key file for testing
func createRelayerKeyFile(t *testing.T, fileName, privKey, password string) {
	// encrypt the private key
	ciphertext, err := crypto.EncryptAES256GCMBase64(privKey, password)
	require.NoError(t, err)

	// create relayer key file
	err = keys.WriteRelayerKeyToFile(fileName, keys.RelayerKey{PrivateKey: ciphertext})
	require.NoError(t, err)
}

// createBadRelayerKeyFile creates a bad relayer key file for testing
func createBadRelayerKeyFile(t *testing.T, fileName string) {
	err := os.WriteFile(fileName, []byte("arbitrary data"), 0o600)
	require.NoError(t, err)
}

func Test_ResolveAddress(t *testing.T) {
	// sample test keys
	solanaPrivKey := sample.SolanaPrivateKey(t)

	tests := []struct {
		name                string
		network             chains.Network
		relayerKey          keys.RelayerKey
		expectedNetworkName string
		expectedAddress     string
		expectedError       string
	}{
		{
			name:    "should resolve solana address",
			network: chains.Network_solana,
			relayerKey: keys.RelayerKey{
				PrivateKey: solanaPrivKey.String(),
			},
			expectedNetworkName: "solana",
			expectedAddress:     solanaPrivKey.PublicKey().String(),
		},
		{
			name:    "should return error if private key is invalid",
			network: chains.Network_solana,
			relayerKey: keys.RelayerKey{
				PrivateKey: "invalid",
			},
			expectedError: "unable to construct solana private key",
		},
		{
			name:    "should return error if network is unsupported",
			network: chains.Network_eth,
			relayerKey: keys.RelayerKey{
				PrivateKey: solanaPrivKey.String(),
			},
			expectedError: "unsupported network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkName, address, err := tt.relayerKey.ResolveAddress(tt.network)
			if tt.expectedError != "" {
				require.Empty(t, networkName)
				require.Empty(t, address)
				require.ErrorContains(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNetworkName, networkName)
			require.Equal(t, tt.expectedAddress, address)
		})
	}
}

func Test_LoadRelayerKey(t *testing.T) {
	// sample test key and temp path
	solanaPrivKey := sample.SolanaPrivateKey(t)
	keyPath := sample.CreateTempDir(t)
	fileName := path.Join(keyPath, "solana.json")

	// create relayer key file
	createRelayerKeyFile(t, fileName, solanaPrivKey.String(), "password")

	// create a bad relayer key file
	keyPath2 := sample.CreateTempDir(t)
	badKeyFile := path.Join(keyPath2, "solana.json")
	createBadRelayerKeyFile(t, badKeyFile)

	// test cases
	tests := []struct {
		name        string
		keyPath     string
		network     chains.Network
		password    string
		expectedKey *keys.RelayerKey
		expectError string
	}{
		{
			name:        "should load relayer key successfully",
			keyPath:     keyPath,
			network:     chains.Network_solana,
			password:    "password",
			expectedKey: &keys.RelayerKey{PrivateKey: solanaPrivKey.String()},
		},
		{
			name:        "it's okay if the relayer key path is blank",
			keyPath:     "",
			network:     chains.Network_solana,
			password:    "",
			expectedKey: nil,
			expectError: "",
		},
		{
			name:        "it's okay if the relayer key path is invalid",
			keyPath:     sample.CreateTempDir(t), // create an empty directory
			network:     chains.Network_solana,
			password:    "",
			expectedKey: nil,
			expectError: "",
		},
		{
			name:        "should return error if network is unsupported",
			keyPath:     keyPath,
			network:     chains.Network_eth,
			password:    "",
			expectedKey: nil,
			expectError: "failed to resolve relayer key file name",
		},
		{
			name:        "should return error if unable to read relayer key file",
			keyPath:     keyPath2,
			network:     chains.Network_solana,
			password:    "",
			expectedKey: nil,
			expectError: "failed to read relayer key file",
		},
		{
			name:        "should return error if password is missing",
			keyPath:     keyPath,
			network:     chains.Network_solana,
			password:    "",
			expectedKey: nil,
			expectError: "password is required to decrypt the private key",
		},
		{
			name:        "should return error if password is incorrect",
			keyPath:     keyPath,
			network:     chains.Network_solana,
			password:    "incorrect",
			expectedKey: nil,
			expectError: "failed to decrypt private key",
		},
	}

	// Iterate over the test cases and run them
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relayerKey, err := keys.LoadRelayerKey(tt.keyPath, tt.network, tt.password)

			if tt.expectError != "" {
				require.ErrorContains(t, err, tt.expectError)
				require.Nil(t, relayerKey)
			} else {
				require.NoError(t, err)
				if tt.expectedKey != nil {
					require.Equal(t, tt.expectedKey.PrivateKey, relayerKey.PrivateKey)
				}
			}
		})
	}
}

func Test_ResolveRelayerKeyPath(t *testing.T) {
	usr, err := user.Current()
	require.NoError(t, err)

	tests := []struct {
		name           string
		relayerKeyPath string
		network        chains.Network
		expectedName   string
		errMessage     string
	}{
		{
			name:           "should resolve relayer key path",
			relayerKeyPath: "~/.zetacored/relayer-keys",
			network:        chains.Network_solana,
			expectedName:   path.Join(usr.HomeDir, ".zetacored/relayer-keys/solana.json"),
		},
		{
			name:           "should return error if network is invalid",
			relayerKeyPath: "~/.zetacored/relayer-keys",
			network:        chains.Network(999),
			errMessage:     "failed to get relayer key file name",
		},
		{
			name:           "should return error if network does not support relayer key",
			relayerKeyPath: "~/.zetacored/relayer-keys",
			network:        chains.Network_eth,
			errMessage:     "does not support relayer key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := keys.ResolveRelayerKeyFile(tt.relayerKeyPath, tt.network)
			if tt.errMessage != "" {
				require.Empty(t, name)
				require.ErrorContains(t, err, tt.errMessage)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedName, name)
		})
	}
}

func Test_ReadWriteRelayerKeyFile(t *testing.T) {
	// sample test key and temp path
	solanaPrivKey := sample.SolanaPrivateKey(t)
	keyPath := sample.CreateTempDir(t)
	fileName := path.Join(keyPath, "solana.json")

	t.Run("should write and read relayer key file", func(t *testing.T) {
		// create relayer key file
		err := keys.WriteRelayerKeyToFile(fileName, keys.RelayerKey{PrivateKey: solanaPrivKey.String()})
		require.NoError(t, err)

		// read relayer key file
		relayerKey, err := keys.ReadRelayerKeyFromFile(fileName)
		require.NoError(t, err)
		require.Equal(t, solanaPrivKey.String(), relayerKey.PrivateKey)
	})

	t.Run("should return error if relayer key file does not exist", func(t *testing.T) {
		noFileName := path.Join(keyPath, "non-existing.json")
		_, err := keys.ReadRelayerKeyFromFile(noFileName)
		require.ErrorContains(t, err, "unable to read relayer key data")
	})

	t.Run("should return error if unmarshalling fails", func(t *testing.T) {
		// create a bad key file
		badKeyFile := path.Join(keyPath, "bad.json")
		createBadRelayerKeyFile(t, badKeyFile)

		// try reading bad key file
		key, err := keys.ReadRelayerKeyFromFile(badKeyFile)
		require.ErrorContains(t, err, "unable to unmarshal relayer key")
		require.Nil(t, key)
	})
}

func Test_IsRelayerPrivateKeyValid(t *testing.T) {
	tests := []struct {
		name    string
		privKey string
		network chains.Network
		result  bool
	}{
		{
			name:    "valid private key - solana",
			privKey: "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
			network: chains.Network(7), // solana
			result:  true,
		},
		{
			name:    "invalid private key - unsupported network",
			privKey: "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
			network: chains.Network(0), // eth
			result:  false,
		},
		{
			name:    "invalid private key - invalid solana private key",
			privKey: "3EMjCcCJg53fMEGVj13UPQpo6p", // too short
			network: chains.Network(7),            // solana
			result:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keys.IsRelayerPrivateKeyValid(tt.privKey, chains.Network(tt.network))
			require.Equal(t, tt.result, result)
		})
	}
}
