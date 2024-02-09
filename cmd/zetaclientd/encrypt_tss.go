package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var encTssCmd = &cobra.Command{
	Use:   "tss-encrypt [file-path] [secret-key]",
	Short: "Utility command to encrypt existing tss key-share file",
	Args:  cobra.ExactArgs(2),
	RunE:  EncryptTSSFile,
}

func init() {
	RootCmd.AddCommand(encTssCmd)
}

func EncryptTSSFile(_ *cobra.Command, args []string) error {
	filePath := args[0]
	secretKey := args[1]

	filePath = filepath.Clean(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return errors.New("file does not contain valid json, may already be encrypted")
	}

	block, err := aes.NewCipher(getFragmentSeed(secretKey))
	if err != nil {
		return err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return os.WriteFile(filePath, cipherText, 0o600)
}

func getFragmentSeed(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	seed := h.Sum(nil)
	return seed
}
