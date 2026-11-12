package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	io "io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// pbkdf2Iterations is the number of PBKDF2 iterations used for key derivation.
	// OWASP 2023 recommends 600000 for PBKDF2-HMAC-SHA256.
	pbkdf2Iterations = 600000
	// pbkdf2KeyLen is the AES-256 key length in bytes.
	pbkdf2KeyLen = 32
)

// getAESKey derives a 32-byte AES key from the given password using PBKDF2.
// This replaces the previous SHA-256-based derivation which was vulnerable to
// offline brute-force attacks.
func getAESKey(password string) []byte {
	return getAESKeyPBKDF2(password)
}

// getAESKeyPBKDF2 derives a key using PBKDF2-HMAC-SHA256 with a static salt.
func getAESKeyPBKDF2(password string) []byte {
	salt := []byte("zeta-chain-aes-256-gcm")
	return pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, pbkdf2KeyLen, sha256.New)
}

// getAESKeyLegacy derives a key using a single SHA-256 round (old insecure method).
// Retained for backward compatibility with encrypted data from before the PBKDF2 migration.
func getAESKeyLegacy(key string) []byte {
	h := sha256.New()
	h.Write([]byte(key))
	return h.Sum(nil)
}

// EncryptAES256GCMBase64 encrypts the given string plaintext using AES-256-GCM with the given password
// and returns the base64-encoded ciphertext.
func EncryptAES256GCMBase64(plaintext string, password string) (string, error) {
	if plaintext == "" {
		return "", errors.New("plaintext must not be empty")
	}
	if password == "" {
		return "", errors.New("password must not be empty")
	}

	ciphertext, err := EncryptAES256GCM([]byte(plaintext), password)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt string plaintext")
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES256GCMBase64 decrypts the given base64-encoded ciphertext using AES-256-GCM with the given password.
func DecryptAES256GCMBase64(ciphertextBase64 string, password string) (string, error) {
	if ciphertextBase64 == "" {
		return "", errors.New("ciphertext must not be empty")
	}
	if password == "" {
		return "", errors.New("password must not be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 ciphertext")
	}

	plaintext, err := DecryptAES256GCM(ciphertext, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt ciphertext")
	}
	return string(plaintext), nil
}

// EncryptAES256GCM encrypts the given plaintext using AES-256-GCM with the given password.
func EncryptAES256GCM(plaintext []byte, password string) ([]byte, error) {
	block, err := aes.NewCipher(getAESKey(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES256GCM decrypts the given ciphertext using AES-256-GCM with the given password.
// It first tries PBKDF2 key derivation, then falls back to legacy SHA-256 for backward compatibility.
func DecryptAES256GCM(ciphertext []byte, password string) ([]byte, error) {
	// try PBKDF2 key derivation first
	plaintext, err := decryptWithKey(ciphertext, getAESKeyPBKDF2(password))
	if err != nil {
		// fall back to legacy SHA-256 key derivation
		plaintext, err = decryptWithKey(ciphertext, getAESKeyLegacy(password))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt with both PBKDF2 and legacy key derivation")
		}
	}
	return plaintext, nil
}

// decryptWithKey decrypts ciphertext with the given AES key.
func decryptWithKey(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// #nosec G407 false positive https://github.com/securego/gosec/issues/1211
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
