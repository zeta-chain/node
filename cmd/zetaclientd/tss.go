package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/crypto"
)

// TSSEncryptFile encrypts the given file with the given secret key
func TSSEncryptFile(_ *cobra.Command, args []string) error {
	var (
		filePath = filepath.Clean(args[0])
		password = args[1]
	)

	// #nosec G304 -- this is a config file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return fmt.Errorf("file %s is not a valid json, may already be encrypted", filePath)
	}

	// encrypt the data
	cipherText, err := crypto.EncryptAES256GCM(data, password)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt data")
	}

	if err := os.WriteFile(filePath, cipherText, 0o600); err != nil {
		return errors.Wrap(err, "failed to write encrypted data to file")
	}

	fmt.Printf("File %s successfully encrypted\n", filePath)

	return nil
}

func TSSGeneratePreParams(_ *cobra.Command, args []string) error {
	startTime := time.Now()
	preParams, err := keygen.GeneratePreParams(time.Second * 300)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(preParams); err != nil {
		return err
	}

	fmt.Printf("Generated new pre-parameters in %s\n", time.Since(startTime).String())

	return nil
}
