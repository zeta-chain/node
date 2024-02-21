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

func TestKeeper_VoteOnInboundBallot(t *testing.T) {

	t.Run("fail if inbound not enabled", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: false,
		})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)

		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInboundDisabled)
	})

	t.Run("fail if sender chain not supported", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
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

		_, _, err = k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if not authorized", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrNotObserver)
	})

	t.Run("fail if receiver chain not supported", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
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

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)

		// set the chain but not supported
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     common.ZetaPrivnetChain().ChainId,
					IsSupported: false,
				},
			},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		_, _, err = k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if inbound contain ZETA but receiver chain doesn't support ZETA", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:                  getValidEthChainIDWithIndex(t, 1),
					IsSupported:              true,
					ZetaTokenContractAddress: "",
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			common.CoinType_Zeta,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidZetaCoinTypes)
	})

	t.Run("can add vote and create ballot", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			common.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.True(t, isFinalized)
		require.True(t, isNew)
	})

	t.Run("can add vote to an existing ballot", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
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
			ObservationType: types.ObservationType_InBoundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			common.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should not be finalized as the threshold is not reached
		require.False(t, isFinalized)
		require.False(t, isNew)
	})
}
