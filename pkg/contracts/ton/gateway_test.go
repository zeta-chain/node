package ton

import (
	"embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestFixtures(t *testing.T) {
	// ACT
	tx := getFixtureTX(t, "01")

	// ASSERT
	require.Equal(t, uint64(26023788000003), tx.Lt)
	require.Equal(t, "cbd6e2261334d08120e2fef428ecbb4e7773606ced878d0e6da204f2b4bf42bf", tx.Hash().Hex())
}

//go:embed testdata
var fixtures embed.FS

// testdata/$name.json tx
func getFixtureTX(t *testing.T, name string) ton.Transaction {
	t.Helper()

	filename := fmt.Sprintf("testdata/%s.json", name)

	b, err := fixtures.ReadFile(filename)
	require.NoError(t, err, filename)

	// bag of cells
	var raw struct {
		BOC string `json:"boc"`
	}

	require.NoError(t, json.Unmarshal(b, &raw))

	cells, err := boc.DeserializeBocHex(raw.BOC)
	require.NoError(t, err)
	require.Len(t, cells, 1)

	cell := cells[0]

	var tx ton.Transaction

	require.NoError(t, tx.UnmarshalTLB(cell, &tlb.Decoder{}))

	return tx
}
