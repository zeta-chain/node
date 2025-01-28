package sample

import (
	"bytes"
	_ "embed"
	"testing"

	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"
)

//go:embed genesis.json
var genesisJSON []byte

func AppGenesis(t *testing.T) *genutiltypes.AppGenesis {
	reader := bytes.NewReader(genesisJSON)
	genDoc, err := genutiltypes.AppGenesisFromReader(reader)
	require.NoError(t, err)
	return genDoc
}
