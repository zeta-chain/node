package main

import (
	"encoding/json"
	"os"
)

// solanaTestKey is a local test private key for Solana
// TODO: use separate keys for each zetaclient in Solana E2E tests
// https://github.com/zeta-chain/node/issues/2614
var solanaTestKey = []uint8{
	199, 16, 63, 28, 125, 103, 131, 13, 6, 94, 68, 109, 13, 68, 132, 17,
	71, 33, 216, 51, 49, 103, 146, 241, 245, 162, 90, 228, 71, 177, 32, 199,
	31, 128, 124, 2, 23, 207, 48, 93, 141, 113, 91, 29, 196, 95, 24, 137,
	170, 194, 90, 4, 124, 113, 12, 222, 166, 209, 119, 19, 78, 20, 99, 5,
}

// createSolanaTestKeyFile creates a solana test key json file
func createSolanaTestKeyFile(keyFile string) error {
	// marshal the byte array to JSON
	keyBytes, err := json.Marshal(solanaTestKey)
	if err != nil {
		return err
	}

	// create file (or overwrite if it already exists)
	// #nosec G304 -- for E2E testing purposes only
	file, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the key bytes to the file
	_, err = file.Write(keyBytes)
	return err
}
