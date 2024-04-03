package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetBallot(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	identifier := sample.ZetaIndex(t)
	b := &types.Ballot{
		Index:                "123",
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

	// overwrite existing ballot
	b = &types.Ballot{
		Index:                "123",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      1,
		BallotThreshold:      sdk.Dec{},
		BallotStatus:         1,
		BallotCreationHeight: 2,
	}
	_, found = k.GetBallot(ctx, identifier)
	require.True(t, found)

	k.SetBallot(ctx, b)

	ballot, found = k.GetBallot(ctx, identifier)
	require.True(t, found)
	require.Equal(t, *b, ballot)
}

func TestKeeper_GetBallotList(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	identifier := sample.ZetaIndex(t)
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

func TestKeeper_GetAllBallots(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	identifier := sample.ZetaIndex(t)
	b := &types.Ballot{
		Index:                "",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      0,
		BallotThreshold:      sdk.Dec{},
		BallotStatus:         0,
		BallotCreationHeight: 1,
	}
	ballots := k.GetAllBallots(ctx)
	require.Empty(t, ballots)

	k.SetBallot(ctx, b)
	ballots = k.GetAllBallots(ctx)
	require.Equal(t, 1, len(ballots))
	require.Equal(t, b, ballots[0])
}

func TestKeeper_GetMaturedBallotList(t *testing.T) {
	t.Run("should return if maturity blocks less than height", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
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
		ctx = ctx.WithBlockHeight(101)
		k.AddBallotToList(ctx, *b)
		list = k.GetMaturedBallotList(ctx)
		require.Equal(t, 1, len(list))
		require.Equal(t, identifier, list[0])
	})

	t.Run("should return empty if maturity blocks greater than height", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
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
		require.Empty(t, list)
	})
}
