package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_VoteBlame(t *testing.T) {
	t.Run("should error if supported chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		res, err := srv.VoteBlame(ctx, &types.MsgVoteBlame{
			ChainId: 1,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if not tombstoned observer", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chainId := getValidEthChainIDWithIndex(t, 0)
		setSupportedChain(ctx, *k, chainId)

		res, err := srv.VoteBlame(ctx, &types.MsgVoteBlame{
			ChainId: chainId,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return response and set blame if finalizing vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chainId := getValidEthChainIDWithIndex(t, 0)
		setSupportedChain(ctx, *k, chainId)

		r := rand.New(rand.NewSource(9))
		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             string(consAddress),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		blameInfo := sample.BlameRecord(t, "index")
		res, err := srv.VoteBlame(ctx, &types.MsgVoteBlame{
			Creator:   accAddressOfValidator.String(),
			ChainId:   chainId,
			BlameInfo: blameInfo,
		})
		require.NoError(t, err)
		require.Equal(t, &types.MsgVoteBlameResponse{}, res)

		blame, found := k.GetBlame(ctx, blameInfo.Index)
		require.True(t, found)
		require.Equal(t, blameInfo, blame)
	})

	t.Run("should error if add vote fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chainId := getValidEthChainIDWithIndex(t, 1)
		setSupportedChain(ctx, *k, chainId)

		r := rand.New(rand.NewSource(9))
		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             string(consAddress),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String(), "Observer2"},
		})
		blameInfo := sample.BlameRecord(t, "index")
		vote := &types.MsgVoteBlame{
			Creator:   accAddressOfValidator.String(),
			ChainId:   chainId,
			BlameInfo: blameInfo,
		}
		ballot := types.Ballot{
			BallotIdentifier: vote.Digest(),
			VoterList:        []string{accAddressOfValidator.String()},
			Votes:            []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:     types.BallotStatus_BallotInProgress,
			BallotThreshold:  sdkmath.LegacyNewDec(2),
		}
		k.SetBallot(ctx, &ballot)

		_, err = srv.VoteBlame(ctx, vote)
		require.Error(t, err)
	})

	t.Run("should return response and not set blame if not finalizing vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chainId := getValidEthChainIDWithIndex(t, 1)
		setSupportedChain(ctx, *k, chainId)

		r := rand.New(rand.NewSource(9))
		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             string(consAddress),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String(), "Observer2"},
		})
		blameInfo := sample.BlameRecord(t, "index")
		vote := &types.MsgVoteBlame{
			Creator:   accAddressOfValidator.String(),
			ChainId:   chainId,
			BlameInfo: blameInfo,
		}
		ballot := types.Ballot{
			BallotIdentifier: vote.Digest(),
			VoterList:        []string{accAddressOfValidator.String()},
			Votes:            []types.VoteType{types.VoteType_NotYetVoted},
			BallotStatus:     types.BallotStatus_BallotInProgress,
			BallotThreshold:  sdkmath.LegacyNewDec(2),
		}
		k.SetBallot(ctx, &ballot)

		res, err := srv.VoteBlame(ctx, vote)
		require.NoError(t, err)
		require.Equal(t, &types.MsgVoteBlameResponse{}, res)

		_, found := k.GetBlame(ctx, blameInfo.Index)
		require.False(t, found)
	})
}
