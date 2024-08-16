package crypto

import (
	fmt "fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
)

// SolanaPrivateKeyFromString converts a base58 encoded private key to a solana.PrivateKey
func SolanaPrivateKeyFromString(privKeyBase58 string) (*solana.PrivateKey, error) {
	privateKey, err := solana.PrivateKeyFromBase58(privKeyBase58)
	if err != nil {
		return nil, errors.Wrap(err, "invalid base58 private key")
	}

	// Solana private keys are 64 bytes long
	if len(privateKey) != 64 {
		return nil, fmt.Errorf("invalid private key length: %d", len(privateKey))
	}

	return &privateKey, nil
}
