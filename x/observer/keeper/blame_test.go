package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_BlameByIdentifier(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	var chainId int64 = 97
	var nonce uint64 = 101
	digest := "85f5e10431f69bc2a14046a13aabaefc660103b6de7a84f75c4b96181d03f0b5"

	index := types.GetBlameIndex(chainId, nonce, digest, 123)

	keeper.SetBlame(ctx, &types.Blame{
		Index:         index,
		FailureReason: "failed to join party",
		Nodes:         nil,
	})

	blameRecords, found := keeper.GetBlame(ctx, index)
	require.True(t, found)
	require.Equal(t, index, blameRecords.Index)
}

func TestKeeper_BlameByChainAndNonce(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	var chainId int64 = 97
	var nonce uint64 = 101
	digest := "85f5e10431f69bc2a14046a13aabaefc660103b6de7a84f75c4b96181d03f0b5"

	index := types.GetBlameIndex(chainId, nonce, digest, 123)

	keeper.SetBlame(ctx, &types.Blame{
		Index:         index,
		FailureReason: "failed to join party",
		Nodes:         nil,
	})

	blameRecords, found := keeper.GetBlamesByChainAndNonce(ctx, chainId, int64(nonce))
	require.True(t, found)
	require.Equal(t, 1, len(blameRecords))
	require.Equal(t, index, blameRecords[0].Index)
}
