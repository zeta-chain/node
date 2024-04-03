package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetBlockHeader(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	blockHash := sample.Hash().Bytes()
	_, found := k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)

	bh := proofs.BlockHeader{
		Height:     1,
		Hash:       blockHash,
		ParentHash: sample.Hash().Bytes(),
		ChainId:    1,
		Header:     proofs.HeaderData{},
	}
	k.SetBlockHeader(ctx, bh)
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.True(t, found)

	k.RemoveBlockHeader(ctx, blockHash)
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)
}
