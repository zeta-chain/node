package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

// TestKeeper_GetBlockHeader tests get, set, and remove block header
func TestKeeper_GetBlockHeader(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	blockHash := sample.Hash().Bytes()
	_, found := k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)

	k.SetBlockHeader(ctx, sample.BlockHeader(blockHash))
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.True(t, found)

	k.RemoveBlockHeader(ctx, blockHash)
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)
}

func TestKeeper_GetAllBlockHeaders(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	b1 := sample.BlockHeader(sample.Hash().Bytes())
	b2 := sample.BlockHeader(sample.Hash().Bytes())
	b3 := sample.BlockHeader(sample.Hash().Bytes())

	k.SetBlockHeader(ctx, b1)
	k.SetBlockHeader(ctx, b2)
	k.SetBlockHeader(ctx, b3)

	list := k.GetAllBlockHeaders(ctx)
	require.Len(t, list, 3)
	require.EqualValues(t, b1, list[0])
	require.EqualValues(t, b2, list[1])
	require.EqualValues(t, b3, list[2])
}
