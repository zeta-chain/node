package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/crypto"
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

// EncryptTSSFile encrypts the given file with the given secret key
func EncryptTSSFile(_ *cobra.Command, args []string) error {
	filePath := args[0]
	password := args[1]

	filePath = filepath.Clean(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return errors.New("file does not contain valid json, may already be encrypted")
	}

	// encrypt the data
	cipherText, err := crypto.EncryptAES256GCM(data, password)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt data")
	}

	return os.WriteFile(filePath, cipherText, 0o600)
}
