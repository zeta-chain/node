package os_test

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	zetaos "github.com/zeta-chain/node/pkg/os"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestResolveHome(t *testing.T) {
	usr, err := user.Current()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		pathIn   string
		expected string
		fail     bool
	}{
		{
			name:     `should resolve home with leading "~/"`,
			pathIn:   "~/tmp/file.json",
			expected: filepath.Clean(filepath.Join(usr.HomeDir, "tmp/file.json")),
		},
		{
			name:     "should resolve '~'",
			pathIn:   `~`,
			expected: filepath.Clean(filepath.Join(usr.HomeDir, "")),
		},
		{
			name:     "should not resolve '~someuser/tmp'",
			pathIn:   `~someuser/tmp`,
			expected: `~someuser/tmp`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pathOut, err := zetaos.ExpandHomeDir(tc.pathIn)
			require.NoError(t, err)
			require.Equal(t, tc.expected, pathOut)
		})
	}
}

func TestFileExists(t *testing.T) {
	path := sample.CreateTempDir(t)

	// create a test file
	existingFile := filepath.Join(path, "test.txt")
	_, err := os.Create(existingFile)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		file     string
		expected bool
	}{
		{
			name:     "should return true for existing file",
			file:     existingFile,
			expected: true,
		},
		{
			name:     "should return false for non-existing file",
			file:     filepath.Join(path, "non-existing.txt"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists := zetaos.FileExists(tc.file)
			require.Equal(t, tc.expected, exists)
		})
	}
}
