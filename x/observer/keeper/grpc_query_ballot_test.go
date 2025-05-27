package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_BallotListForHeight(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// Act
		res, err := k.BallotListForHeight(ctx, nil)

		// Assert
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return empty list if no ballots", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetBallotList(ctx, &types.BallotListForHeight{
			Height:           1,
			BallotsIndexList: nil,
		})

		// Act
		res, err := k.BallotListForHeight(ctx, &types.QueryBallotListForHeightRequest{
			Height: 1,
		})

		// Assert
		require.NoError(t, err)
		require.Nil(t, res.BallotList.BallotsIndexList)
	})

	t.Run("should return error if list not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// Act
		res, err := k.BallotListForHeight(ctx, &types.QueryBallotListForHeightRequest{
			Height: 1,
		})

		// Assert
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return list if exists", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		ballotList := &types.BallotListForHeight{
			Height:           1,
			BallotsIndexList: []string{"index-1", "index-2"},
		}

		k.SetBallotList(ctx, ballotList)

		// Act
		res, err := k.BallotListForHeight(ctx, &types.QueryBallotListForHeightRequest{
			Height: 1,
		})

		// Assert
		require.NoError(t, err)
		require.Equal(t, ballotList.BallotsIndexList, res.BallotList.BallotsIndexList)
	})

}

func TestKeeper_HasVoted(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.HasVoted(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return false if ballot not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.HasVoted(wctx, &types.QueryHasVotedRequest{
			BallotIdentifier: "test",
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryHasVotedResponse{
			HasVoted: false,
		}, res)
	})

	t.Run("should return true if ballot found and voted", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		voter := sample.AccAddress()
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList:        []string{voter},
			Votes:            []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:     types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		res, err := k.HasVoted(wctx, &types.QueryHasVotedRequest{
			BallotIdentifier: "index",
			VoterAddress:     voter,
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryHasVotedResponse{
			HasVoted: true,
		}, res)
	})

	t.Run("should return false if ballot found and not voted", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		voter := sample.AccAddress()
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList:        []string{voter},
			Votes:            []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:     types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		res, err := k.HasVoted(wctx, &types.QueryHasVotedRequest{
			BallotIdentifier: "index",
			VoterAddress:     sample.AccAddress(),
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryHasVotedResponse{
			HasVoted: false,
		}, res)
	})
}

func TestKeeper_BallotByIdentifier(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BallotByIdentifier(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return nil if ballot not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BallotByIdentifier(wctx, &types.QueryBallotByIdentifierRequest{
			BallotIdentifier: "test",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return ballot if exists", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		voter := sample.AccAddress()
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList:        []string{voter},
			Votes:            []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:     types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		res, err := k.BallotByIdentifier(wctx, &types.QueryBallotByIdentifierRequest{
			BallotIdentifier: "index",
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryBallotByIdentifierResponse{
			BallotIdentifier: ballot.BallotIdentifier,
			Voters: []*types.VoterList{
				{
					VoterAddress: voter,
					VoteType:     types.VoteType_SuccessObservation,
				},
			},
			ObservationType: ballot.ObservationType,
			BallotStatus:    ballot.BallotStatus,
		}, res)
	})

	t.Run("should return 100 ballots if more exist and limit is not provided", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		numOfBallots := 1000

		ballots := make([]types.Ballot, numOfBallots)
		for i := 0; i < numOfBallots; i++ {
			ballot := types.Ballot{
				Index:                "",
				BallotIdentifier:     fmt.Sprintf("index-%d", i),
				VoterList:            []string{sample.AccAddress()},
				Votes:                []types.VoteType{types.VoteType_SuccessObservation},
				BallotStatus:         types.BallotStatus_BallotInProgress,
				BallotCreationHeight: 1,
				BallotThreshold:      sdkmath.LegacyMustNewDecFromStr("0.5"),
			}
			k.SetBallot(ctx, &ballot)
			ballots[i] = ballot
		}

		res, err := k.Ballots(wctx, &types.QueryBallotsRequest{})
		require.NoError(t, err)
		require.Len(t, res.Ballots, 100)
	})

	t.Run("should return limit number of ballots if limit is provided", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		numOfBallots := 1000

		ballots := make([]types.Ballot, numOfBallots)
		for i := 0; i < numOfBallots; i++ {
			ballot := types.Ballot{
				Index:                "",
				BallotIdentifier:     fmt.Sprintf("index-%d", i),
				VoterList:            []string{sample.AccAddress()},
				Votes:                []types.VoteType{types.VoteType_SuccessObservation},
				BallotStatus:         types.BallotStatus_BallotInProgress,
				BallotCreationHeight: 1,
				BallotThreshold:      sdkmath.LegacyMustNewDecFromStr("0.5"),
			}
			k.SetBallot(ctx, &ballot)
			ballots[i] = ballot
		}

		res, err := k.Ballots(wctx, &types.QueryBallotsRequest{
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Ballots, 10)
	})
}

func TestKeeper_Ballots(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.Ballots(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return empty list if no ballots", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.Ballots(wctx, &types.QueryBallotsRequest{})
		require.NoError(t, err)
		require.Empty(t, res.Ballots)
	})

	t.Run("should return all ballots", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		ballots := make([]types.Ballot, 10)
		for i := 0; i < 10; i++ {
			ballot := types.Ballot{
				Index:                "",
				BallotIdentifier:     fmt.Sprintf("index-%d", i),
				VoterList:            []string{sample.AccAddress()},
				Votes:                []types.VoteType{types.VoteType_SuccessObservation},
				BallotStatus:         types.BallotStatus_BallotInProgress,
				BallotCreationHeight: 1 + int64(i),
				BallotThreshold:      sdkmath.LegacyMustNewDecFromStr("0.5"),
			}
			k.SetBallot(ctx, &ballot)
			ballots[i] = ballot
		}

		res, err := k.Ballots(wctx, &types.QueryBallotsRequest{})
		require.NoError(t, err)
		require.ElementsMatch(t, ballots, res.Ballots)

		firstBallotCreationHeight := res.Ballots[0].BallotCreationHeight
		for _, ballot := range res.Ballots {
			require.GreaterOrEqual(t, ballot.BallotCreationHeight, firstBallotCreationHeight)
		}
	})
}
