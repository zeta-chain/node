package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetBallot(t *testing.T) {
	k, ctx := SetupKeeper(t)
	identifier := "0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d"
	b := &types.Ballot{
		Index:                "",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      0,
		BallotThreshold:      sdk.Dec{},
		BallotStatus:         0,
		BallotCreationHeight: 1,
	}
	_, found := k.GetBallot(ctx, identifier)
	require.False(t, found)

	k.SetBallot(ctx, b)

	ballot, found := k.GetBallot(ctx, identifier)
	require.True(t, found)
	require.Equal(t, *b, ballot)
}

func TestKeeper_GetBallotList(t *testing.T) {
	k, ctx := SetupKeeper(t)
	identifier := "0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d"
	b := &types.Ballot{
		Index:                "",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      0,
		BallotThreshold:      sdk.Dec{},
		BallotStatus:         0,
		BallotCreationHeight: 1,
	}
	_, found := k.GetBallotList(ctx, 1)
	require.False(t, found)

	k.AddBallotToList(ctx, *b)
	list, found := k.GetBallotList(ctx, 1)
	require.True(t, found)
	require.Equal(t, 1, len(list.BallotsIndexList))
	require.Equal(t, identifier, list.BallotsIndexList[0])
}

func TestKeeper_GetMaturedBallotList(t *testing.T) {
	k, ctx := SetupKeeper(t)
	identifier := "0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d"
	b := &types.Ballot{
		Index:                "",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      0,
		BallotThreshold:      sdk.Dec{},
		BallotStatus:         0,
		BallotCreationHeight: 1,
	}
	list := k.GetMaturedBallotList(ctx)
	require.Empty(t, list)
	ctx = ctx.WithBlockHeight(1)
	k.AddBallotToList(ctx, *b)
	list = k.GetMaturedBallotList(ctx)
	require.Equal(t, 1, len(list))
	require.Equal(t, identifier, list[0])
}
