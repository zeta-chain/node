package sample

import (
	_ "embed"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

//go:embed genesis.json
var genesisJSON []byte

func GenDoc(t *testing.T) *types.GenesisDoc {
	genDoc, err := types.GenesisDocFromJSON(genesisJSON)
	require.NoError(t, err)
	return genDoc
}
