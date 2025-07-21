package utils

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"

	"github.com/zeta-chain/node/pkg/parsers"
)

// ScriptPKToAddress is a hex string for P2WPKH script
func ScriptPKToAddress(scriptPKHex string, params *chaincfg.Params) string {
	pkh, err := hex.DecodeString(scriptPKHex[4:])
	if err == nil {
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pkh, params)
		if err == nil {
			return addr.EncodeAddress()
		}
	}
	return ""
}

type infoLogger interface {
	Info(message string, args ...interface{})
}

type NoopLogger struct{}

func (nl NoopLogger) Info(_ string, _ ...interface{}) {}

type testingKey struct{}

// WithTesting allows to store a testing.T instance in the context
func WithTesting(ctx context.Context, t require.TestingT) context.Context {
	return context.WithValue(ctx, testingKey{}, t)
}

// TestingFromContext extracts require.TestingT from the context or panics.
func TestingFromContext(ctx context.Context) require.TestingT {
	t, ok := ctx.Value(testingKey{}).(require.TestingT)
	if !ok {
		panic("context missing require.TestingT key")
	}

	return t
}

func MinimumVersionCheck(testVersion, zetacoredVersion string) bool {
	// If major version is "v0", return true regardless of comparison
	if semver.Major(zetacoredVersion) == "v0" {
		return true
	}

	// Otherwise, return true if zetacoredVersion >= testVersion
	return semver.Compare(zetacoredVersion, testVersion) >= 0
}

// FetchNodePubkey retrieves the public key of the new validator node.
func FetchNodePubkey(host string) (string, error) {
	// #nosec G204 validation of host should be done by the caller
	cmd := exec.Command("ssh", "-q", fmt.Sprintf("root@%s", host), "zetacored tendermint show-validator")
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s, stderr: %s, error: %w", cmd.String(), stderr.String(), err)
	}
	output := out.String()
	output = strings.TrimSpace(output)
	return output, nil
}

// FetchHotkeyAddress retrieves the hotkey address of a new validator.
func FetchHotkeyAddress(host string) (parsers.ObserverInfoReader, error) {
	// #nosec G204 validation of host should be done by the caller
	cmd := exec.Command("ssh", "-q", fmt.Sprintf("root@%s", host), "cat ~/.zetacored/os.json")
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return parsers.ObserverInfoReader{}, fmt.Errorf(
			"failed to run command: %s, stderr: %s, error: %w",
			cmd.String(),
			stderr.String(),
			err,
		)
	}
	output := out.String()

	observerInfo := parsers.ObserverInfoReader{}

	err = json.Unmarshal([]byte(output), &observerInfo)
	if err != nil {
		return parsers.ObserverInfoReader{}, fmt.Errorf("failed to unmarshal observer info: %w", err)
	}

	return observerInfo, nil
}

// WorkDir gets the current working directory of E2E test
func WorkDir(t require.TestingT) string {
	dir, err := os.Getwd()
	require.NoError(t, err)

	return dir
}
