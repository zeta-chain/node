package crypto_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/crypto"
)

func Test_EncryptDecryptAES256GCM(t *testing.T) {
	tests := []struct {
		name        string
		plaintext   string
		encryptPass string
		decryptPass string
		modifyFunc  func([]byte) []byte
		fail        bool
		errMsg      string
	}{
		{
			name:        "Successful encryption and decryption",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password",
			fail:        false,
		},
		{
			name:        "Decryption with incorrect key should fail",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password2",
			fail:        true,
		},
		{
			name:        "Decryption with ciphertext too short should fail",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password",
			modifyFunc: func(ciphertext []byte) []byte {
				// truncate the ciphertext, nonce size is 12 bytes
				return ciphertext[:10]
			},
			fail:   true,
			errMsg: "ciphertext too short",
		},
		{
			name:        "Decryption with corrupted ciphertext should fail",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password",
			modifyFunc: func(ciphertext []byte) []byte {
				// flip the last bit of the ciphertext
				ciphertext[len(ciphertext)-1] ^= 0x01
				return ciphertext
			},
			fail: true,
		},
		{
			name:        "Decryption with incorrect nonce should fail",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password",
			modifyFunc: func(ciphertext []byte) []byte {
				// flip the first bit of the nonce
				ciphertext[0] ^= 0x01
				return ciphertext
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := crypto.EncryptAES256GCM([]byte(tt.plaintext), tt.encryptPass)
			require.NoError(t, err)

			// modify the encrypted data if needed
			if tt.modifyFunc != nil {
				encrypted = tt.modifyFunc(encrypted)
			}

			// decrypt the data
			decrypted, err := crypto.DecryptAES256GCM(encrypted, tt.decryptPass)
			if tt.fail {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			require.True(t, bytes.Equal(decrypted, []byte(tt.plaintext)), "decrypted plaintext does not match")
		})
	}
}

func Test_EncryptAES256GCMBase64(t *testing.T) {
	tests := []struct {
		name         string
		plaintext    string
		encryptPass  string
		decryptPass  string
		errorMessage string
	}{
		{
			name:        "Successful encryption and decryption",
			plaintext:   "Hello, World!",
			encryptPass: "my_password",
			decryptPass: "my_password",
		},
		{
			name:         "Encryption with empty plaintext should fail",
			plaintext:    "",
			errorMessage: "plaintext must not be empty",
		},
		{
			name:         "Encryption with empty password should fail",
			plaintext:    "Hello, World!",
			encryptPass:  "",
			errorMessage: "password must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// encrypt the data
			ciphertextBase64, err := crypto.EncryptAES256GCMBase64(tt.plaintext, tt.encryptPass)
			if tt.errorMessage != "" {
				require.ErrorContains(t, err, tt.errorMessage)
				return
			}

			// decrypt the data
			decrypted, err := crypto.DecryptAES256GCMBase64(ciphertextBase64, tt.decryptPass)
			require.NoError(t, err)

			require.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func Test_DecryptAES256GCMBase64(t *testing.T) {
	tests := []struct {
		name             string
		ciphertextBase64 string
		plaintext        string
		decryptKey       string
		modifyFunc       func(string) string
		errorMessage     string
	}{
		{
			name:             "Successful decryption",
			ciphertextBase64: "CXLWgHdVeZQwVOZZyHeZ5n5VB+eVSLaWFF0v0QOm9DyB7XSiHDwhNwQ=",
			plaintext:        "Hello, World!",
			decryptKey:       "my_password",
		},
		{
			name:             "Decryption with empty ciphertext should fail",
			ciphertextBase64: "",
			decryptKey:       "my_password",
			errorMessage:     "ciphertext must not be empty",
		},
		{
			name:             "Decryption with empty password should fail",
			ciphertextBase64: "CXLWgHdVeZQwVOZZyHeZ5n5VB+eVSLaWFF0v0QOm9DyB7XSiHDwhNwQ=",
			decryptKey:       "",
			errorMessage:     "password must not be empty",
		},
		{
			name:             "Decryption with invalid base64 ciphertext should fail",
			ciphertextBase64: "CXLWgHdVeZQwVOZZyHeZ5n5VB*eVSLaWFF0v0QOm9DyB7XSiHDwhNwQ=", // use '*' instead of '+'
			decryptKey:       "my_password",
			errorMessage:     "failed to decode base64 ciphertext",
		},
		{
			name:             "Decryption with incorrect password should fail",
			ciphertextBase64: "CXLWgHdVeZQwVOZZyHeZ5n5VB+eVSLaWFF0v0QOm9DyB7XSiHDwhNwQ=",
			decryptKey:       "my_password2",
			errorMessage:     "failed to decrypt ciphertext",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertextBase64 := tt.ciphertextBase64

			// modify the encrypted data if needed
			if tt.modifyFunc != nil {
				ciphertextBase64 = tt.modifyFunc(ciphertextBase64)
			}

			// decrypt the data
			decrypted, err := crypto.DecryptAES256GCMBase64(ciphertextBase64, tt.decryptKey)
			if tt.errorMessage != "" {
				require.ErrorContains(t, err, tt.errorMessage)
				return
			}

			require.Equal(t, tt.plaintext, decrypted)
		})
	}
}
