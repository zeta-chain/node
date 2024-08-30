package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	io "io"

	"github.com/pkg/errors"
)

// EncryptAES256GCMBase64 encrypts the given string plaintext using AES-256-GCM with the given password and returns the base64-encoded ciphertext.
func EncryptAES256GCMBase64(plaintext string, password string) (string, error) {
	// validate the input
	if plaintext == "" {
		return "", errors.New("plaintext must not be empty")
	}
	if password == "" {
		return "", errors.New("password must not be empty")
	}

	// encrypt the plaintext
	ciphertext, err := EncryptAES256GCM([]byte(plaintext), password)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt string plaintext")
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES256GCMBase64 decrypts the given base64-encoded ciphertext using AES-256-GCM with the given password.
func DecryptAES256GCMBase64(ciphertextBase64 string, password string) (string, error) {
	// validate the input
	if ciphertextBase64 == "" {
		return "", errors.New("ciphertext must not be empty")
	}
	if password == "" {
		return "", errors.New("password must not be empty")
	}

	// decode the base64-encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 ciphertext")
	}

	// decrypt the ciphertext
	plaintext, err := DecryptAES256GCM(ciphertext, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt ciphertext")
	}
	return string(plaintext), nil
}

// EncryptAES256GCM encrypts the given plaintext using AES-256-GCM with the given password.
func EncryptAES256GCM(plaintext []byte, password string) ([]byte, error) {
	// create AES cipher
	block, err := aes.NewCipher(getAESKey(password))
	if err != nil {
		return nil, err
	}

	// create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptAES256GCM decrypts the given ciphertext using AES-256-GCM with the given password.
func DecryptAES256GCM(ciphertext []byte, password string) ([]byte, error) {
	// create AES cipher
	block, err := aes.NewCipher(getAESKey(password))
	if err != nil {
		return nil, err
	}

	// create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// get the nonce size
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// extract the nonce from the ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// getAESKey uses SHA-256 to create a 32-byte key for AES encryption.
func getAESKey(key string) []byte {
	h := sha256.New()
	h.Write([]byte(key))

	return h.Sum(nil)
}
