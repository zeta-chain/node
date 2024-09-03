package crypto_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/crypto"
)

func Test_SolanaPrivateKeyFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output *solana.PrivateKey
		errMsg string
	}{
		{
			name:  "valid private key",
			input: "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
			output: func() *solana.PrivateKey {
				privKey, _ := solana.PrivateKeyFromBase58(
					"3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
				)
				return &privKey
			}(),
		},
		{
			name:   "invalid private key - too short",
			input:  "oR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ",
			output: nil,
			errMsg: "invalid private key length: 38",
		},
		{
			name:   "invalid private key - too long",
			input:  "3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQdJ",
			output: nil,
			errMsg: "invalid private key length: 66",
		},
		{
			name:   "invalid private key - bad base58 encoding",
			input:  "!!!InvalidBase58!!!",
			output: nil,
			errMsg: "invalid base58 private key",
		},
		{
			name:   "invalid private key - empty string",
			input:  "",
			output: nil,
			errMsg: "invalid base58 private key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := crypto.SolanaPrivateKeyFromString(tt.input)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.output.String(), result.String())
		})
	}
}
