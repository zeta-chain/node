package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMsgServer_RemoveChainInfo(t *testing.T) {
	t.Run("can't remove chain info if not authorized", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgRemoveChainInfo{
			Creator: sample.AccAddress(),
		}
		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(&msg),
				AuthorizedPolicy: types.PolicyType_groupAdmin,
			},
		}})

		_, err := msgServer.RemoveChainInfo(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)
	})

	t.Run("can remove chain from chain info", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainID := int64(42)
		chainInfo := sample.ChainInfo(chainID)
		k.SetChainInfo(ctx, chainInfo)

		// Act
		_, err := msgServer.RemoveChainInfo(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainInfo{
			Creator: admin,
			ChainId: chainID,
		})

		// Assert
		require.NoError(t, err)
		storedChains, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Len(t, storedChains.Chains, 2)
		require.NotContains(t, storedChains.Chains, chainInfo.Chains[0])
	})

	t.Run("can remove chain from chain info containing only 1 chain", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainID := int64(42)
		chainInfo := types.ChainInfo{Chains: []chains.Chain{{ChainId: chainID}}}
		k.SetChainInfo(ctx, chainInfo)

		// Act
		_, err := msgServer.RemoveChainInfo(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainInfo{
			Creator: admin,
			ChainId: chainID,
		})

		// Assert
		require.NoError(t, err)
		storedChains, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Len(t, storedChains.Chains, 0)
		require.NotContains(t, storedChains.Chains, chainInfo.Chains[0])
	})

	t.Run("can't remove chain from chain info if chain info not found", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})

		// Act
		_, err := msgServer.RemoveChainInfo(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainInfo{
			Creator: admin,
			ChainId: 42,
		})

		// Assert
		require.ErrorIs(t, err, types.ErrChainInfoNotFound)
	})

	t.Run("can't remove chain from chain info if chain not found", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainInfo := sample.ChainInfo(42)
		k.SetChainInfo(ctx, chainInfo)

		// Act
		_, err := msgServer.RemoveChainInfo(sdk.WrapSDKContext(ctx), &types.MsgRemoveChainInfo{
			Creator: admin,
			ChainId: 103,
		})

		// Assert
		require.NoError(t, err)
		storedChains, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Len(t, storedChains.Chains, 3)
	})
}

func TestMsgServer_RemoveChain(t *testing.T) {
	tt := []struct {
		name          string
		chainInfo     types.ChainInfo
		removeChainID int64
		expected      types.ChainInfo
	}{
		{
			name:          "can remove chain from chain info",
			chainInfo:     types.ChainInfo{Chains: []chains.Chain{{ChainId: 42}, {ChainId: 43}}},
			removeChainID: 42,
			expected:      types.ChainInfo{Chains: []chains.Chain{{ChainId: 43}}},
		},
		{
			name:          "can't remove chain from chain info if chain not found",
			chainInfo:     types.ChainInfo{Chains: []chains.Chain{{ChainId: 42}, {ChainId: 43}}},
			removeChainID: 103,
			expected:      types.ChainInfo{Chains: []chains.Chain{{ChainId: 42}, {ChainId: 43}}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, keeper.RemoveChain(tc.chainInfo, tc.removeChainID))
		})
	}

}
