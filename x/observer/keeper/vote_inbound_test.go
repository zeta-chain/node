package keeper_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"testing"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_VoteOnInboundBallot(t *testing.T) {

	t.Run("fail if inbound not enabled", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: false,
		})

		_, err := k.VoteOnInboundBallot(
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
		k, ctx := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{})

		_, err := k.VoteOnInboundBallot(
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

		_, err = k.VoteOnInboundBallot(
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
		k, ctx := keepertest.ObserverKeeper(t)

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

		_, err := k.VoteOnInboundBallot(
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
		k, ctx := keepertest.ObserverKeeper(t)

		observer := sample.AccAddress()

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

		_, err := k.VoteOnInboundBallot(
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

		_, err = k.VoteOnInboundBallot(
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
		k, ctx := keepertest.ObserverKeeper(t)

		observer := sample.AccAddress()

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
					ChainId:                  common.ZetaPrivnetChain().ChainId,
					IsSupported:              true,
					ZetaTokenContractAddress: "",
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})

		_, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			common.ZetaPrivnetChain().ChainId,
			common.CoinType_Zeta,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidZetaCoinTypes)
	})

	t.Run("can add vote to inbound ballot", func(t *testing.T) {

	})
}
