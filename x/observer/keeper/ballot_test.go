package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/x/observer/keeper"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetBallot(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	identifier := sample.ZetaIndex(t)
	b := &types.Ballot{
		Index:                "123",
		BallotIdentifier:     identifier,
		VoterList:            nil,
		ObservationType:      0,
		BallotThreshold:      sdk.ZeroDec(),
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
		BallotThreshold:      sdk.ZeroDec(),
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
	t.Run("get existing ballot list", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
		b := &types.Ballot{
			Index:                "",
			BallotIdentifier:     identifier,
			VoterList:            nil,
			ObservationType:      0,
			BallotThreshold:      sdk.ZeroDec(),
			BallotStatus:         0,
			BallotCreationHeight: 1,
		}
		_, found := k.GetBallotListForHeight(ctx, 1)
		require.False(t, found)

		k.AddBallotToList(ctx, *b)
		list, found := k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)
		require.Equal(t, 1, len(list.BallotsIndexList))
		require.Equal(t, identifier, list.BallotsIndexList[0])
	})

	t.Run("get non-existing ballot list", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
		b := &types.Ballot{
			Index:                "",
			BallotIdentifier:     identifier,
			VoterList:            nil,
			ObservationType:      0,
			BallotThreshold:      sdk.ZeroDec(),
			BallotStatus:         0,
			BallotCreationHeight: 1,
		}
		_, found := k.GetBallotListForHeight(ctx, 1)
		require.False(t, found)

		k.AddBallotToList(ctx, *b)
		list, found := k.GetBallotListForHeight(ctx, -10)
		require.False(t, found)
		require.Nil(t, list.BallotsIndexList)
	})

}

func TestKeeper_GetMaturedBallots(t *testing.T) {
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
	ctx = ctx.WithBlockHeight(2)
	_, found := k.GetMaturedBallots(ctx, 1)
	require.False(t, found)

	k.AddBallotToList(ctx, *b)
	list, found := k.GetMaturedBallots(ctx, 1)
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
		BallotThreshold:      sdk.ZeroDec(),
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

func TestKeeper_DeleteBallot(t *testing.T) {
	t.Run("delete existing ballot", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
		b := &types.Ballot{
			BallotIdentifier: identifier,
		}
		k.SetBallot(ctx, b)
		_, found := k.GetBallot(ctx, identifier)
		require.True(t, found)

		//Act
		k.DeleteBallot(ctx, identifier)

		//Assert
		_, found = k.GetBallot(ctx, identifier)
		require.False(t, found)
	})

	t.Run("delete non-existing ballot,nothing happens", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		identifier := sample.ZetaIndex(t)
		numberOfBallots := 10
		for i := 0; i < numberOfBallots; i++ {
			k.SetBallot(ctx, &types.Ballot{
				BallotIdentifier: sample.ZetaIndex(t),
			})
		}

		require.Len(t, k.GetAllBallots(ctx), numberOfBallots)

		//Act
		k.DeleteBallot(ctx, identifier)

		//Assert
		_, found := k.GetBallot(ctx, identifier)
		require.False(t, found)
		require.Len(t, k.GetAllBallots(ctx), numberOfBallots)
	})
}

func TestKeeper_DeleteBallotList(t *testing.T) {
	t.Run("delete existing ballot list", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		numberOfBallotLists := 10
		for i := 0; i < numberOfBallotLists; i++ {
			k.AddBallotToList(ctx, types.Ballot{
				Index:                sample.ZetaIndex(t),
				BallotCreationHeight: 1,
			})
		}

		_, found := k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)

		//Act
		k.DeleteBallotListForHeight(ctx, 1)

		//Assert
		_, found = k.GetBallotListForHeight(ctx, 1)
		require.False(t, found)
	})

	t.Run("delete non-existing ballot list, nothing happens", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		numberOfBallotLists := 10
		for i := 0; i < numberOfBallotLists; i++ {
			k.AddBallotToList(ctx, types.Ballot{
				Index:                sample.ZetaIndex(t),
				BallotCreationHeight: 1,
			})
		}

		_, found := k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)

		//Act
		k.DeleteBallotListForHeight(ctx, 2)

		//Assert
		_, found = k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)
	})
}

func TestKeeper_ClearMaturedBallots(t *testing.T) {
	t.Run("clear matured ballots successfully", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		numberOfBallots := 10
		ballots := make([]types.Ballot, numberOfBallots)
		for i := 0; i < numberOfBallots; i++ {
			b := types.Ballot{
				BallotIdentifier:     sample.ZetaIndex(t),
				BallotCreationHeight: 1,
			}
			k.AddBallotToList(ctx, b)
			k.SetBallot(ctx, &b)
			ballots[i] = b
		}
		_, found := k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)
		require.Equal(t, numberOfBallots, len(k.GetAllBallots(ctx)))

		//Act
		k.ClearMaturedBallotsAndBallotList(ctx, 0)

		//Assert
		for _, b := range ballots {
			_, found = k.GetBallot(ctx, b.BallotIdentifier)
			require.False(t, found)
		}
		_, found = k.GetBallotListForHeight(ctx, 0)
		require.False(t, found)
		eventCount := 0
		for _, event := range ctx.EventManager().Events() {
			if event.Type == "zetachain.zetacore.observer.EventBallotDeleted" {
				eventCount++
			}
		}
		require.Equal(t, numberOfBallots, eventCount)
	})

	t.Run("clear only ballotList if no ballots are found", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		numberOfBallots := 10
		ballots := make([]types.Ballot, numberOfBallots)
		for i := 0; i < numberOfBallots; i++ {
			b := types.Ballot{
				BallotIdentifier:     sample.ZetaIndex(t),
				BallotCreationHeight: 1,
			}
			k.AddBallotToList(ctx, b)
			ballots[i] = b
		}
		_, found := k.GetBallotListForHeight(ctx, 1)
		require.True(t, found)
		require.Equal(t, 0, len(k.GetAllBallots(ctx)))

		//Act
		k.ClearMaturedBallotsAndBallotList(ctx, 0)

		//Assert
		_, found = k.GetBallotListForHeight(ctx, 1)
		require.False(t, found)
		require.Equal(t, 0, len(k.GetAllBallots(ctx)))
	})

	t.Run("do nothing if ballot list for height is not found", func(t *testing.T) {
		// Note this condition should never happen in production as the ballot list for height should always be set when saving a ballot to state
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		numberOfBallots := 10
		ballots := make([]types.Ballot, numberOfBallots)
		for i := 0; i < numberOfBallots; i++ {
			b := types.Ballot{
				BallotIdentifier:     sample.ZetaIndex(t),
				BallotCreationHeight: 1,
			}
			k.SetBallot(ctx, &b)
			ballots[i] = b
		}
		_, found := k.GetBallotListForHeight(ctx, 1)
		require.False(t, found)
		require.Equal(t, numberOfBallots, len(k.GetAllBallots(ctx)))

		//Act
		k.ClearMaturedBallotsAndBallotList(ctx, 0)

		//Assert
		for _, b := range ballots {
			_, found = k.GetBallot(ctx, b.BallotIdentifier)
			require.True(t, found)
		}
	})

}

func TestGetMaturedBallotHeight(t *testing.T) {
	tt := []struct {
		name           string
		currentHeight  int64
		maturityBlocks int64
		expectedHeight int64
	}{
		{
			name:           "maturity blocks is 0",
			currentHeight:  10,
			maturityBlocks: 0,
			expectedHeight: 10,
		},
		{
			name:           "maturity blocks is same as current height",
			currentHeight:  10,
			maturityBlocks: 10,
			expectedHeight: 0,
		},
		{
			name:           "maturity blocks is less than current height",
			currentHeight:  10,
			maturityBlocks: 5,
			expectedHeight: 5,
		},
		{
			name:           "maturity blocks is greater than current height",
			currentHeight:  5,
			maturityBlocks: 10,
			expectedHeight: -5,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, ctx, _, _ := keepertest.ObserverKeeper(t)
			ctx = ctx.WithBlockHeight(tc.currentHeight)
			require.Equal(t, tc.expectedHeight, keeper.GetMaturedBallotHeightFunc(ctx, tc.maturityBlocks))
		})
	}
}
