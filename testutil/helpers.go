package testutil

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const helpersFile = "testutil/helpers.go"

// Must a helper that terminates the program if the error is not nil.
func Must[T any](v T, err error) T {
	NoError(err)
	return v
}

// NoError terminates the program if the error is not nil.
func NoError(err error) {
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Printf("Unable to continue execution: %s.\nStacktrace:\n", err)

	for _, line := range assert.CallerInfo() {
		if strings.Contains(line, helpersFile) {
			continue
		}

		fmt.Println("  ", line)
	}

	os.Exit(1)
}

// HexToBytes convert hex string to bytes
func HexToBytes(t *testing.T, hexStr string) []byte {
	bytes, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	return bytes
}
