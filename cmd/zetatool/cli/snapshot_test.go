package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writePinFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "pins.txt")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

// TestReadPinFile covers the shapes an operator-supplied pin file takes: a plain
// list, comment and blank lines, trailing comments, and a TSV whose first column
// is the address.
func TestReadPinFile(t *testing.T) {
	// ARRANGE
	path := writePinFile(t, `# non-claimable contracts (mainnet)
# source: eth_getCode over the top holders

0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf  # WZETA, canonical wrapped ZETA
   0x1840aabe58042e4080e248c7384381cb37178215

0x0fc8f47c5bd1c89b0e2f08bfe72df4ea9dff61aa	680661	TTUV2Native
zeta1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq
`)

	// ACT
	pins, err := readPinFile(path)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, []string{
		"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf",
		"0x1840aabe58042e4080e248c7384381cb37178215",
		"0x0fc8f47c5bd1c89b0e2f08bfe72df4ea9dff61aa",
		"zeta1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
	}, pins, "only the first field of each non-comment line, in file order")
}

func TestReadPinFileCommentsOnly(t *testing.T) {
	// ARRANGE
	path := writePinFile(t, "# nothing here\n\n   \n#0xdeadbeef disabled\n")

	// ACT
	pins, err := readPinFile(path)

	// ASSERT
	require.NoError(t, err)
	require.Empty(t, pins)
}

func TestReadPinFileMissing(t *testing.T) {
	// ARRANGE
	path := filepath.Join(t.TempDir(), "does-not-exist.txt")

	// ACT
	_, err := readPinFile(path)

	// ASSERT
	require.ErrorContains(t, err, "read pin file")
}

// TestReadPinFileMainnetList asserts the committed mainnet pin list parses to the
// 16 contract addresses it documents.
func TestReadPinFileMainnetList(t *testing.T) {
	// ACT
	pins, err := readPinFile("../../../contrib/snapshot/pins_mainnet.txt")

	// ASSERT
	require.NoError(t, err)
	require.Len(t, pins, 16)
	require.Equal(t, "0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf", pins[0], "WZETA is the largest holder")
}
