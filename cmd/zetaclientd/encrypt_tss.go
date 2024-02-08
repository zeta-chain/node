package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var encTssCmd = &cobra.Command{
	Use:   "tss-encrypt",
	Short: "Utility command to encrypt existing tss key-share file",
	RunE:  EncryptTSSFile,
}

type TSSArgs struct {
	secretKey string
	filePath  string
}

var tssArgs = TSSArgs{}

func init() {
	RootCmd.AddCommand(encTssCmd)

	encTssCmd.Flags().StringVar(&tssArgs.secretKey, "secret", "", "tss-encrpyt --secret p@$$w0rd")
	encTssCmd.Flags().StringVar(&tssArgs.filePath, "filepath", "", "tss-encrpyt --filepath ./file.json")
}

func EncryptTSSFile(_ *cobra.Command, _ []string) error {
	data, err := os.ReadFile(tssArgs.filePath)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return errors.New("file does not contain valid json, may already be encrypted")
	}

	block, err := aes.NewCipher(getFragmentSeed(tssArgs.secretKey))
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
	return os.WriteFile(tssArgs.filePath, cipherText, 0o600)
}

func getFragmentSeed(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	seed := h.Sum(nil)
	return seed
}
