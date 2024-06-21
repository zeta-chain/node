package sample

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// create a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	tempPath, err := os.MkdirTemp("", "tempdir-")
	require.NoError(t, err)
	return tempPath
}
