package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

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
}
