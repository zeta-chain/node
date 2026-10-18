package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const auditPin = "0x3333333333333333333333333333333333333333"

// auditExport is the path of a ported audit export in the snapshot testdata.
func auditExport(name string) string {
	return filepath.Join("..", "snapshot", "testdata", name+".json")
}

// runSnapshotCmd runs `zetatool snapshot` on input with extra flags and
// returns its stdout.
func runSnapshotCmd(t *testing.T, input string, extra ...string) (string, error) {
	t.Helper()
	cmd := NewSnapshotCMD()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	args := []string{"--input", input, "--output", t.TempDir(), "--chain-id", "7000", "--pin", auditPin}
	cmd.SetArgs(append(args, extra...))
	err := cmd.Execute()
	return out.String(), err
}

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

func TestSnapshotCompleteExport(t *testing.T) {
	// ACT
	out, err := runSnapshotCmd(t, auditExport("complete"))

	// ASSERT
	require.NoError(t, err)
	require.Contains(t, out, "remainder by category (azeta):")
	require.Contains(t, out, "pins:          5000000000000000000")
	require.Contains(t, out, "conversion dust (18->9): 0 azeta")
}

// TestSnapshotRefusalNotSkippable runs an export that lost its unbonding record
// and checks it is refused even with --supply-check=false.
func TestSnapshotRefusalNotSkippable(t *testing.T) {
	for _, supplyCheck := range []string{"true", "false"} {
		t.Run("supply-check="+supplyCheck, func(t *testing.T) {
			// ACT
			_, err := runSnapshotCmd(t, auditExport("omitted_unbonding"), "--supply-check="+supplyCheck)

			// ASSERT
			require.ErrorContains(t, err, "not-bonded pool balance 100000000000000000000 != ")
		})
	}
}

func TestSnapshotMissingDistribution(t *testing.T) {
	// ARRANGE
	raw, err := os.ReadFile(auditExport("complete"))
	require.NoError(t, err)
	var doc map[string]map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &doc))
	delete(doc["app_state"], "distribution")
	stripped, err := json.Marshal(doc)
	require.NoError(t, err)
	input := filepath.Join(t.TempDir(), "export.json")
	require.NoError(t, os.WriteFile(input, stripped, 0o600))

	// ACT
	_, err = runSnapshotCmd(t, input)

	// ASSERT
	require.ErrorContains(t, err, `app_state is missing the "distribution" module`)
}

// TestSnapshotMainnetRequiresPins runs mainnet with no pins against an input that
// does not exist, so passing proves the refusal comes before the export is read.
func TestSnapshotMainnetRequiresPins(t *testing.T) {
	// ARRANGE
	cmd := NewSnapshotCMD()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	input := filepath.Join(t.TempDir(), "does-not-exist.json")
	cmd.SetArgs([]string{"--input", input, "--output", t.TempDir(), "--chain-id", "7000"})

	// ACT
	err := cmd.Execute()

	// ASSERT
	require.ErrorContains(t, err, "mainnet (--chain-id 7000) requires pins: pass --pin-file contrib/snapshot/pins_mainnet.txt")
}
