package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_VoteOnOutboundBallot(t *testing.T) {
	t.Run("fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{})

		_, _, _, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			chains.ReceiveStatus_success,
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
			chains.ReceiveStatus_success,
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if receive status is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

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
			chains.ReceiveStatus(1000),
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidStatus)
	})

	t.Run("fail if sender is not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

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
			chains.ReceiveStatus_success,
			sample.AccAddress(),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrNotObserver)
	})

	t.Run("can add vote and create ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

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
			chains.ReceiveStatus_success,
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

	t.Run("fail if can not add vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

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
		ballot := types.Ballot{
			BallotIdentifier: "index",
			VoterList:        []string{observer},
			// already voted
			Votes:           []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:    types.BallotStatus_BallotInProgress,
			BallotThreshold: sdkmath.LegacyNewDec(2),
		}
		k.SetBallot(ctx, &ballot)
		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			chains.ReceiveStatus_success,
			observer,
		)
		require.Error(t, err)
		require.False(t, isFinalized)
		require.False(t, isNew)
	})

	t.Run("can add vote and create ballot without finalizing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// threshold high enough to not finalize the ballot
		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         getValidEthChainIDWithIndex(t, 0),
					IsSupported:     true,
					BallotThreshold: threshold,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{
				observer,
				sample.AccAddress(),
			},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			chains.ReceiveStatus_success,
			observer,
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.False(t, isFinalized)
		require.True(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, "index")
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})

	t.Run("can add vote to an existing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

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
		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)
		ballot := types.Ballot{
			BallotIdentifier: "index",
			VoterList: []string{
				sample.AccAddress(),
				sample.AccAddress(),
				observer,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(5),
			ObservationType: types.ObservationType_OutboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			chains.ReceiveStatus_success,
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

	t.Run("can add vote to an existing ballot and finalize ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

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
		threshold, err := sdkmath.LegacyNewDecFromStr("0.1")
		require.NoError(t, err)
		ballot := types.Ballot{
			BallotIdentifier: "index",
			VoterList: []string{
				observer,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(3),
			ObservationType: types.ObservationType_OutboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, ballot, _, err := k.VoteOnOutboundBallot(
			ctx,
			"index",
			getValidEthChainIDWithIndex(t, 0),
			chains.ReceiveStatus_success,
			observer,
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.True(t, isFinalized)
		require.False(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, "index")
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})
}
