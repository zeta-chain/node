package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_VoteOnOutboundBallot(t *testing.T) {
	t.Run("fail if chain is not supported", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{})

		_, _, _, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus_Success,
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)

		// set the chain but not supported
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: false,
				},
			},
		})

		_, _, _, _, err = k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus_Success,
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if receive status is invalid", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})

		_, _, _, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus(1000),
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidStatus)
	})

	t.Run("fail if sender is not authorized", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{})

		_, _, _, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus_Success,
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrNotObserver)
	})

	t.Run("can add vote and create ballot", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus_Success,
			observer,
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.True(t, isFinalized)
		require.True(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, "index")
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})

	t.Run("can add vote to an existing ballot", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		// set a ballot
		threshold, err := sdk.NewDecFromStr("0.7")
		require.NoError(t, err)
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList: []string{
				sample.AccAddress(),
				sample.AccAddress(),
				observer,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(5),
			ObservationType: types.ObservationType_OutBoundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			common.ReceiveStatus_Success,
			observer,
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.False(t, isFinalized)
		require.False(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, "index")
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})
}
