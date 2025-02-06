package sui

import (
	"bytes"
	"encoding/base64"
	"testing"

	sui_signer "github.com/block-vision/sui-go-sdk/signer"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"
)

// test vector:
// see https://github.com/MystenLabs/sui/blob/199f06d25ce85f0270a1a5a0396156bb2b83122c/sdk/typescript/test/unit/cryptography/secp256k1-keypair.test.ts

var VALID_SECP256K1_SECRET_KEY = []byte{
	59, 148, 11, 85, 134, 130, 61, 253, 2, 174, 59, 70, 27, 180, 51, 107, 94, 203,
	174, 253, 102, 39, 170, 146, 46, 252, 4, 143, 236, 12, 136, 28,
}

var VALID_SECP256K1_PUBLIC_KEY = []byte{
	2, 29, 21, 35, 7, 198, 183, 43, 14, 208, 65, 139, 14, 112, 205, 128, 231, 245,
	41, 91, 141, 134, 245, 114, 45, 63, 82, 19, 251, 210, 57, 79, 54,
}

func TestSignerSecp256k1FromSecretKey(t *testing.T) {
	// Create signer from secret key
	signer := NewSignerSecp256k1FromSecretKey(VALID_SECP256K1_SECRET_KEY)

	// Get public key bytes
	pubKey := signer.GetPublicKey()

	// Compare with expected public key
	if !bytes.Equal(pubKey, VALID_SECP256K1_PUBLIC_KEY) {
		t.Errorf("Public key mismatch\nexpected: %v\ngot: %v", VALID_SECP256K1_PUBLIC_KEY, pubKey)
	}
}

// Test keypairs generated and exported with
//
// sui client new-address secp256k1
// sui keytool export
//
// See https://github.com/sui-foundation/sips/blob/main/sips/sip-15.md for encoding info
func TestSuiSecp256k1Keypair(t *testing.T) {
	tests := []struct {
		name                 string
		privKey              string
		expectedAddress      string
		expectedPubkeyBase64 string
	}{
		{
			name:                 "example 1",
			privKey:              "suiprivkey1q8h7ejwfcdn6gc2x9kddtd9vld3kvuvtr9gdtj9qcf7nqnw94tvl79cwfq4",
			expectedAddress:      "0x68f6d05fd44366bd431405e8ea32e2d2f8e330d98e0089c9447ecfbbdf6e236f",
			expectedPubkeyBase64: "AQOhmtkY2bTZGXRZmXZLo495i5Dz+FgJvM7bbnUCWlL2hg==",
		},
		{
			name:                 "example 2",
			privKey:              "suiprivkey1qxghtp2vncr94s8h7ctvgf58gy27l9nch75ka2jh37qr90xuyxhlwk5khxc",
			expectedAddress:      "0x8ec6f13affbf48c73550567f2a1cb8f05474c0133ebf745f91a2f3b971c1ec9a",
			expectedPubkeyBase64: "AQIgu/14lUhVMEWjIB0RQ80ARJiH/xQIw4KJTEhDhHcjEQ==",
		},
		{
			name:                 "example 3",
			privKey:              "suiprivkey1q99wkv3fj352cn97yu5r9jwqhcvyyk6t9scwrftyjgqgflandfgc70r74hg",
			expectedAddress:      "0xa0f6b839f7945065ebdd030cec8e6e30d01a74c8cb31b1945909fd426c2cef80",
			expectedPubkeyBase64: "AQIRTpSIUP2GOAWTEHOTYbyUjlfxpHKwPvsgwZw3G6h/RQ==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, data, err := bech32.DecodeAndConvert(tt.privKey)
			require.NoError(t, err)

			// assert this is a secp256k1 private key
			require.Equal(t, byte(sui_signer.SigntureFlagSecp256k1), data[0])

			flagLessData := data[1:]

			signer := NewSignerSecp256k1FromSecretKey(flagLessData)

			address := signer.Address()
			require.Equal(t, tt.expectedAddress, address)

			pubkeyBytes := signer.GetFlaggedPublicKey()
			pubkeyEncoded := base64.StdEncoding.EncodeToString(pubkeyBytes)
			require.Equal(t, tt.expectedPubkeyBase64, pubkeyEncoded)

			// we don't have a good way outside e2e to verify the signature is correct, but let's excertise it anyway
			_, err = signer.SignTransactionBlock("ZXhhbXBsZQo=")
			require.NoError(t, err)
		})
	}
}
